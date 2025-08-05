package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetDisplayRank(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	dispResp := model.SetDisplayRankResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(dispResp)
}

func SetDeck(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	userId, err := strconv.Atoi(ctx.GetString("userid"))
	if ss.CheckErr(err) {
		return
	}

	deckReq := model.UnitDeckReq{}
	err = json.Unmarshal([]byte(ctx.PostForm("request_data")), &deckReq)
	if ss.CheckErr(err) {
		return
	}

	// 原有队伍信息
	var userDeckId []int
	err = ss.UserEng.Table("user_deck_m").Cols("id").Where("user_id = ?", userId).Find(&userDeckId)
	if ss.CheckErr(err) {
		return
	}

	// 删除全部原有队伍成员
	_, err = ss.UserEng.Table("deck_unit_m").In("user_deck_id", userDeckId).Delete()
	if ss.CheckErr(err) {
		return
	}

	// 删除全部原有队伍
	_, err = ss.UserEng.Table("user_deck_m").In("id", userDeckId).Delete()
	if ss.CheckErr(err) {
		return
	}

	// 遍历新队伍
	for _, deck := range deckReq.UnitDeckList {
		// 新队伍信息
		userDeck := model.UserDeckData{
			DeckID:     deck.UnitDeckID,
			MainFlag:   deck.MainFlag,
			DeckName:   deck.DeckName,
			UserID:     userId,
			InsertDate: time.Now().Unix(),
		}
		_, err = ss.UserEng.Table("user_deck_m").Insert(&userDeck)
		if ss.CheckErr(err) {
			return
		}
		userDeckId := userDeck.ID
		// fmt.Println("新队伍 ID:", userDeckId)

		// 队伍成员信息
		for _, unit := range deck.UnitDeckDetail {
			// 成员信息
			newUnitData := model.UnitData{}
			exists, err := ss.UserEng.Table("user_unit_m").Where("unit_owning_user_id = ?", unit.UnitOwningUserID).Exist()
			if ss.CheckErr(err) {
				return
			}
			if exists {
				// fmt.Println("新成员为用户增加成员")
				_, err = ss.UserEng.Table("user_unit_m").Where("unit_owning_user_id = ?", unit.UnitOwningUserID).Get(&newUnitData)
				if ss.CheckErr(err) {
					return
				}
			} else {
				exists, err := ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", unit.UnitOwningUserID).Exist()
				if ss.CheckErr(err) {
					return
				}
				if exists {
					// fmt.Println("新成员为公共成员")
					_, err = ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", unit.UnitOwningUserID).Get(&newUnitData)
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
			newUnitDeckData := model.UnitDeckData{}
			b, err := json.Marshal(newUnitData)
			if ss.CheckErr(err) {
				return
			}
			err = json.Unmarshal(b, &newUnitDeckData)
			if ss.CheckErr(err) {
				return
			}
			newUnitDeckData.BeforeLove = newUnitDeckData.MaxLove
			newUnitDeckData.Position = unit.Position
			newUnitDeckData.UserDeckID = userDeckId
			newUnitDeckData.InsertData = time.Now().Unix()

			_, err = ss.UserEng.Table("deck_unit_m").Insert(&newUnitDeckData)
			if ss.CheckErr(err) {
				return
			}
		}
	}

	dispResp := model.SetDeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(dispResp)
}

func SetDeckName(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	userId, err := strconv.Atoi(ctx.GetString("userid"))
	if ss.CheckErr(err) {
		return
	}

	deckReq := model.DeckNameReq{}
	err = json.Unmarshal([]byte(ctx.PostForm("request_data")), &deckReq)
	if ss.CheckErr(err) {
		return
	}

	exists, err := ss.UserEng.Table("user_deck_m").Where("user_id = ? AND deck_id = ?", userId, deckReq.UnitDeckID).Exist()
	if ss.CheckErr(err) {
		return
	}
	if !exists {
		ctx.String(http.StatusForbidden, ErrorMsg)
		return
	}
	userDeck := model.UserDeckData{
		DeckName: deckReq.DeckName,
	}
	_, err = ss.UserEng.Table("user_deck_m").Update(&userDeck, &model.UserDeckData{
		UserID: userId,
		DeckID: deckReq.UnitDeckID,
	})
	if ss.CheckErr(err) {
		return
	}

	dispResp := model.SetDeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(dispResp)
}

func WearAccessory(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	req := model.WearAccessoryReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &req)
	if ss.CheckErr(err) {
		return
	}

	// 取下饰品
	for _, v := range req.Remove {
		fmt.Println("Remove:", v.AccessoryOwningUserID, v.UnitOwningUserID)
		_, err := ss.UserEng.Table("accessory_wear_m").
			Where("accessory_owning_user_id = ? AND unit_owning_user_id = ? AND user_id = ?", v.AccessoryOwningUserID, v.UnitOwningUserID, ctx.GetString("userid")).
			Delete()
		if ss.CheckErr(err) {
			return
		}
	}

	// 佩戴饰品
	for _, v := range req.Wear {
		fmt.Println("Wear:", v.AccessoryOwningUserID, v.UnitOwningUserID)
		data := model.AccessoryWearData{
			AccessoryOwningUserID: v.AccessoryOwningUserID,
			UnitOwningUserID:      v.UnitOwningUserID,
			UserId:                ctx.GetString("userid"),
		}
		_, err := ss.UserEng.Table("accessory_wear_m").Insert(&data)
		if ss.CheckErr(err) {
			return
		}
	}

	wearResp := model.AwardSetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(wearResp)
}

func RemoveSkillEquip(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	req := model.SkillEquipReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &req)
	if ss.CheckErr(err) {
		return
	}

	// 取下宝石
	for _, v := range req.Remove {
		fmt.Println("Remove:", v.UnitOwningUserID, v.UnitRemovableSkillID)
		_, err := ss.UserEng.Table("skill_equip_m").
			Where("unit_removable_skill_id = ? AND unit_owning_user_id = ? AND user_id = ?", v.UnitRemovableSkillID, v.UnitOwningUserID, ctx.GetString("userid")).
			Delete()
		if ss.CheckErr(err) {
			return
		}
	}

	// 佩戴宝石
	for _, v := range req.Equip {
		fmt.Println("Equip:", v.UnitOwningUserID, v.UnitRemovableSkillID)
		data := model.SkillEquipData{
			UnitRemovableSkillId: v.UnitRemovableSkillID,
			UnitOwningUserID:     v.UnitOwningUserID,
			UserId:               ctx.GetString("userid"),
		}
		_, err := ss.UserEng.Table("skill_equip_m").Insert(&data)
		if ss.CheckErr(err) {
			return
		}
	}

	wearResp := model.AwardSetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(wearResp)
}
