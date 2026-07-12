package unit

import (
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitapischema "honoka-chan/internal/schema/api/unit"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func deck(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	deckReq := unitschema.DeckReq{}
	err := honokautils.ParseRequestData(ctx, &deckReq)
	if ss.CheckErr(err) {
		return
	}

	// 原有队伍信息
	var userDeckID []int
	err = ss.UserEng.Table("user_deck").Cols("id").Where("user_id = ?", ss.UserID).Find(&userDeckID)
	if ss.CheckErr(err) {
		return
	}

	// 删除全部原有队伍成员
	_, err = ss.UserEng.Table("user_deck_unit").In("user_deck_id", userDeckID).Delete()
	if ss.CheckErr(err) {
		return
	}

	// 删除全部原有队伍
	_, err = ss.UserEng.Table("user_deck").In("id", userDeckID).Delete()
	if ss.CheckErr(err) {
		return
	}

	// 遍历新队伍
	for _, deck := range deckReq.UnitDeckList {
		// 新队伍信息
		userDeck := unitapischema.UserDeckData{
			DeckID:     deck.UnitDeckID,
			MainFlag:   deck.MainFlag,
			DeckName:   deck.DeckName,
			UserID:     ss.UserID,
			InsertDate: time.Now().Unix(),
		}
		_, err = ss.UserEng.Table("user_deck").Insert(&userDeck)
		if ss.CheckErr(err) {
			return
		}
		userDeckId := userDeck.ID
		// fmt.Println("新队伍 ID:", userDeckId)

		// 队伍成员信息
		for _, u := range deck.UnitDeckDetail {
			// 成员信息
			unitData := unitmodel.UnitDataMap{}
			_, err = ss.GetBasicUnitInfo().
				Where("a.unit_owning_user_id = ?", u.UnitOwningUserID).Get(&unitData)
			if ss.CheckErr(err) {
				return
			}
			// fmt.Println("新的成员信息:", unitData)

			// 插入新成员信息
			newUnitDeckData := usermodel.UserDeckUnit{
				UserDeckID:       userDeckId,
				UnitOwningUserID: unitData.UnitOwningUserID,
				UnitID:           unitData.UnitID,
				Position:         u.Position,
				Level:            unitData.Level,
				LevelLimitID:     unitData.LevelLimitID,
				DisplayRank:      unitData.DisplayRank,
				Love:             unitData.Love,
				UnitSkillLevel:   unitData.UnitSkillLevel,
				IsRankMax:        unitData.IsRankMax,
				IsLoveMax:        unitData.IsLoveMax,
				IsLevelMax:       unitData.IsLevelMax,
				IsSigned:         unitData.IsSigned,
				BeforeLove:       unitData.MaxLove,
				MaxLove:          unitData.MaxLove,
				UserID:           ss.UserID,
				InsertDate:       time.Now().Unix(),
			}
			_, err = ss.UserEng.Table("user_deck_unit").Insert(&newUnitDeckData)
			if ss.CheckErr(err) {
				return
			}
		}
	}

	ss.Respond(unitschema.DeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/deck", middleware.Common, deck)
}
