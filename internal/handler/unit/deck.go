package unit

import (
	"encoding/json"
	"errors"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	profileapischema "honoka-chan/internal/schema/api/profile"
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func deck(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	deckReq := unit.DeckReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &deckReq)
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
			newUnitData := profileapischema.UnitData{}
			exists, err := ss.UserEng.Table("user_unit").Where("unit_owning_user_id = ?", u.UnitOwningUserID).Exist()
			if ss.CheckErr(err) {
				return
			}
			if exists {
				// fmt.Println("新成员为用户增加成员")
				_, err = ss.UserEng.Table("user_unit").Where("unit_owning_user_id = ?", u.UnitOwningUserID).Get(&newUnitData)
				if ss.CheckErr(err) {
					return
				}
			} else {
				exists, err := ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", u.UnitOwningUserID).Exist()
				if ss.CheckErr(err) {
					return
				}
				if exists {
					// fmt.Println("新成员为公共成员")
					_, err = ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", u.UnitOwningUserID).Get(&newUnitData)
					if ss.CheckErr(err) {
						return
					}
				} else {
					// fmt.Println("新成员不存在")
					err = errors.New("新成员不存在")
				}
				if ss.CheckErr(err) {
					return
				}
			}
			// fmt.Println("新的成员信息:", newUnitData)

			// 插入新成员信息
			newUnitDeckData := unitapischema.UnitDeckData{}
			b, err := json.Marshal(newUnitData)
			if ss.CheckErr(err) {
				return
			}
			err = json.Unmarshal(b, &newUnitDeckData)
			if ss.CheckErr(err) {
				return
			}
			newUnitDeckData.BeforeLove = newUnitDeckData.MaxLove
			newUnitDeckData.Position = u.Position
			newUnitDeckData.UserDeckID = userDeckId
			newUnitDeckData.InsertData = time.Now().Unix()

			_, err = ss.UserEng.Table("user_deck_unit").Insert(&newUnitDeckData)
			if ss.CheckErr(err) {
				return
			}
		}
	}

	ss.Respond(unit.DeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/deck", middleware.Common, deck)
}
