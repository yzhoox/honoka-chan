package live

import (
	"encoding/json"
	"fmt"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	honokautils "honoka-chan/pkg/utils"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func play(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playReq := live.PlayReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &playReq)
	if ss.CheckErr(err) {
		return
	}

	tDifficultyId := playReq.LiveDifficultyID
	difficultyId, err := strconv.Atoi(tDifficultyId)
	if ss.CheckErr(err) {
		return
	}
	deckId := playReq.UnitDeckID

	// Save Deck Id for /live/reward
	key := "live_deck_" + strconv.Itoa(ss.UserID)
	err = db.Ldb.Set([]byte(key), []byte(strconv.Itoa(deckId)))
	if ss.CheckErr(err) {
		return
	}

	// Song type: normal / special
	// sqlite3 doesn't support FULL OUTER JOIN so use UNION ALL here.
	sql := `SELECT notes_setting_asset,c_rank_score,b_rank_score,a_rank_score,s_rank_score,ac_flag,swing_flag FROM live_setting_m WHERE live_setting_id IN (SELECT live_setting_id FROM normal_live_m WHERE live_difficulty_id = ? UNION ALL SELECT live_setting_id FROM special_live_m WHERE live_difficulty_id = ?)`
	var notes_setting_asset string
	var c_rank_score, b_rank_score, a_rank_score, s_rank_score, ac_flag, swing_flag int
	err = ss.MainEng.DB().QueryRow(sql, difficultyId, difficultyId).Scan(&notes_setting_asset, &c_rank_score, &b_rank_score, &a_rank_score, &s_rank_score, &ac_flag, &swing_flag)
	if ss.CheckErr(err) {
		return
	}

	// fmt.Println(notes_setting_asset)
	// fmt.Println(c_rank_score, b_rank_score, a_rank_score, s_rank_score)

	notes := []live.NotesList{}
	// fmt.Println("./assets/serverdata/beatmaps/" + notes_setting_asset)
	notes_list := honokautils.ReadAllText("./assets/serverdata/beatmaps/" + notes_setting_asset)
	err = json.Unmarshal([]byte(notes_list), &notes)
	if ss.CheckErr(err) {
		return
	}

	ranks := []live.RankInfo{}
	ranks = append(ranks, live.RankInfo{
		Rank:    5,
		RankMin: 0,
		RankMax: c_rank_score,
	}, live.RankInfo{
		Rank:    4,
		RankMin: c_rank_score + 1,
		RankMax: b_rank_score,
	}, live.RankInfo{
		Rank:    3,
		RankMin: b_rank_score + 1,
		RankMax: a_rank_score,
	}, live.RankInfo{
		Rank:    2,
		RankMin: a_rank_score + 1,
		RankMax: s_rank_score,
	}, live.RankInfo{
		Rank:    1,
		RankMin: s_rank_score + 1,
		RankMax: 0,
	})

	//  ss.UserEng.ShowSQL(true)
	// ss.MainEng.ShowSQL(true)
	owningIdList := []int{}
	err = ss.UserEng.Table("user_deck_unit").Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_id = ? AND deck_id = ?", ss.UserID, deckId).Cols("unit_owning_user_id").
		OrderBy("user_deck_unit.position ASC").Find(&owningIdList)
	if ss.CheckErr(err) {
		return
	}

	unitList := []live.UnitList{}
	var totalSmile, totalPure, totalCool, maxLove float64
	var totalHp int
	for _, owningId := range owningIdList {
		var uId int
		exists, err := ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", owningId).Cols("unit_id").Get(&uId)
		if ss.CheckErr(err) {
			return
		}

		var maxHp, attrId, unitTypeId int
		var baseSmile, basePure, baseCool, smileMax, pureMax, coolMax float64
		if exists {
			// 公共卡片仅为100级属性
			_, err = ss.MainEng.Table("unit_m").Where("unit_id = ?", uId).
				Join("LEFT", "unit_rarity_m", "unit_m.rarity = unit_rarity_m.rarity").
				Cols("attribute_id,hp_max,smile_max,pure_max,cool_max,after_love_max,unit_type_id").
				Get(&attrId, &maxHp, &baseSmile, &basePure, &baseCool, &maxLove, &unitTypeId)
			if ss.CheckErr(err) {
				return
			}
		} else {
			// 用户卡片暂时固定为满级350级
			exists, err := ss.UserEng.Table("user_unit").Where("unit_owning_user_id = ?", owningId).Cols("unit_id").Get(&uId)
			if ss.CheckErr(err) {
				return
			}

			if exists {
				// 卡片100级基础属性
				_, err = ss.MainEng.Table("unit_m").Where("unit_id = ?", uId).
					Join("LEFT", "unit_rarity_m", "unit_m.rarity = unit_rarity_m.rarity").
					Cols("attribute_id,hp_max,smile_max,pure_max,cool_max,after_love_max,unit_type_id").
					Get(&attrId, &maxHp, &baseSmile, &basePure, &baseCool, &maxLove, &unitTypeId)
				if ss.CheckErr(err) {
					return
				}

				// 增量属性
				var diffSmile, diffPure, diffCool float64
				_, err = ss.MainEng.Table("unit_level_limit_pattern_m").Where("unit_level_limit_id = 1 AND unit_level = 350").
					Cols("smile_diff,pure_diff,cool_diff").Get(&diffSmile, &diffPure, &diffCool)
				if ss.CheckErr(err) {
					return
				}

				// 更新卡片属性（注意这里是负数，要用减号）
				baseSmile -= diffSmile
				basePure -= diffPure
				baseCool -= diffCool
			} else {
				err = fmt.Errorf("no such unit owning id: %d", owningId)
			}
			if ss.CheckErr(err) {
				return
			}
		}

		// 饰品属性加成（满级）
		var accessoryOwningId int
		_, err = ss.UserEng.Table("user_accessory_wear").Where("unit_owning_user_id = ?", owningId).
			Cols("accessory_owning_user_id").Get(&accessoryOwningId)
		if ss.CheckErr(err) {
			return
		}
		var smileAccessory, pureAccessory, coolAccessory float64
		_, err = ss.MainEng.Table("common_accessory_m").Join("LEFT", "accessory_m", "common_accessory_m.accessory_id = accessory_m.accessory_id").
			Where("accessory_owning_user_id = ?", accessoryOwningId).Cols("smile_max,pure_max,cool_max").
			Get(&smileAccessory, &pureAccessory, &coolAccessory)
		if ss.CheckErr(err) {
			return
		}
		// fmt.Println("基础属性:", baseSmile, basePure, baseCool)

		// 饰品属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		baseSmile += smileAccessory
		basePure += pureAccessory
		baseCool += coolAccessory
		// fmt.Println("饰品属性加成:", smileAccessory, pureAccessory, coolAccessory)
		// fmt.Println("饰品属性加成后的基础属性:", baseSmile, basePure, baseCool)

		// 回忆画廊属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		var smileBuff, pureBuff, coolBuff float64
		_, err = ss.MainEng.Table("museum_contents_m").Select("SUM(smile_buff),SUM(pure_buff),SUM(cool_buff)").Get(&smileBuff, &pureBuff, &coolBuff)
		if ss.CheckErr(err) {
			return
		}
		baseSmile += smileBuff
		basePure += pureBuff
		baseCool += coolBuff
		// fmt.Println("回忆画廊属性加成:", smileBuff, pureBuff, coolBuff)
		// fmt.Println("回忆画廊属性加成后的基础属性:", baseSmile, basePure, baseCool)

		// 绊属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		switch attrId {
		case 1:
			baseSmile += maxLove
		case 2:
			basePure += maxLove
		case 3:
			baseCool += maxLove
		}
		// fmt.Println("绊属性加成:", maxLove)
		// fmt.Println("绊属性加成后的基础属性:", baseSmile, basePure, baseCool)

		// 宝石属性加成
		var kissSmile, kissPure, kissCool float64
		var skillSmile, skillPure, skillCool float64
		// var mainCenterSmile, mainCenterPure, mainCenterCool float64
		// var secCenterSmile, secCenterPure, secCenterCool float64

		// 宝石加成（满级）
		removableSkillIds := []int{}
		err = ss.UserEng.Table("user_unit_skill_equip").Where("unit_owning_user_id = ? AND user_id = ?", owningId, ss.UserID).
			Cols("unit_removable_skill_id").Find(&removableSkillIds)
		if ss.CheckErr(err) {
			return
		}

		for _, sk := range removableSkillIds {
			// 判断宝石效果类型（效果范围、效果类型、效果值、是否固定数值）
			var effectRange, effectType, fixedValueFlag, refType int
			var effectValue float64
			_, err = ss.MainEng.Table("unit_removable_skill_m").Where("unit_removable_skill_id = ?", sk).
				Cols("effect_range,effect_type,effect_value,fixed_value_flag,target_reference_type").
				Get(&effectRange, &effectType, &effectValue, &fixedValueFlag, &refType)
			if ss.CheckErr(err) {
				return
			}

			if fixedValueFlag == 1 {
				// 吻、眼神属性加成（固定数值）
				switch effectType {
				case 1:
					kissSmile += effectValue
				case 2:
					kissPure += effectValue
				case 3:
					kissCool += effectValue
				}
				// fmt.Println("吻、眼神属性加成:", kissSmile, kissPure, kissCool)
			} else {
				// 仅效果类型为1、2、3的有属性加成
				if effectType == 1 || effectType == 2 || effectType == 3 {
					// 加成范围：2:全员 1:非全员
					if effectRange == 2 {
						switch effectType {
						case 1:
							skillSmile += math.Ceil(baseSmile * (effectValue / 100))
						case 2:
							skillPure += math.Ceil(basePure * (effectValue / 100))
						case 3:
							skillCool += math.Ceil(baseCool * (effectValue / 100))
						}
						// fmt.Println("全员类宝石属性加成:", skillSmile, skillPure, skillCool)
					} else {
						// refType: 1 -> 年级类加成, target_type -> 指定年级（这里不需要使用，因为能装上宝石肯定是符合的）
						// refType: 2 -> 个宝
						// refType: 3 -> 爆分、奶、判宝石, 0 -> 竞技场宝石
						if refType == 1 || refType == 2 { // 年级类和个宝都是百分比加成
							if effectType == 1 {
								skillSmile += math.Ceil(baseSmile * (effectValue / 100))
							} else if effectType == 2 {
								skillPure += math.Ceil(basePure * (effectValue / 100))
							} else if effectValue == 3 {
								skillCool += math.Ceil(baseCool * (effectValue / 100))
							}
							// fmt.Println("年级类宝石、个宝属性加成:", skillSmile, skillPure, skillCool)
						}
					}
				}
			}
		}

		// 单卡属性
		smileMax = baseSmile + kissSmile + skillSmile
		pureMax = basePure + kissPure + skillPure
		coolMax = baseCool + kissCool + skillCool

		// 主唱技能加成
		var myCenterUnitId int
		_, err = ss.UserEng.Table("user_deck_unit").Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
			Where("user_deck.deck_id = ? AND user_deck.user_id = ? AND user_deck_unit.position = 5", playReq.UnitDeckID, ss.UserID).
			Cols("user_deck_unit.unit_id").Get(&myCenterUnitId)
		if ss.CheckErr(err) {
			return
		}

		// 主唱技能加成：主C技能（这里不使用新C技能，即以某属性的百分比提升另一属性）
		var myAttrId int
		var myEffectValue float64
		_, err = ss.MainEng.Table("unit_m").
			Join("LEFT", "unit_leader_skill_m", "unit_m.default_leader_skill_id = unit_leader_skill_m.unit_leader_skill_id").
			Where("unit_m.unit_id = ?", myCenterUnitId).
			Cols("unit_m.attribute_id,unit_leader_skill_m.effect_value").Get(&myAttrId, &myEffectValue)
		if ss.CheckErr(err) {
			return
		}
		var myCenterSmile, myCenterPure, myCenterCool float64
		switch myAttrId {
		case 1:
			myCenterSmile = math.Ceil(smileMax * (myEffectValue / 100))
		case 2:
			myCenterPure = math.Ceil(pureMax * (myEffectValue / 100))
		case 3:
			myCenterCool = math.Ceil(coolMax * (myEffectValue / 100))
		}
		// fmt.Println("主C技能属性加成:", myCenterSmile, myCenterPure, myCenterCool)

		// 主唱技能加成：副C技能
		var mySubEffectValue float64
		var myMemberTagId int
		_, err = ss.MainEng.Table("unit_m").
			Join("LEFT", "unit_leader_skill_extra_m", "unit_m.default_leader_skill_id = unit_leader_skill_extra_m.unit_leader_skill_id").
			Where("unit_m.unit_id = ?", myCenterUnitId).
			Cols("unit_leader_skill_extra_m.effect_value,unit_leader_skill_extra_m.member_tag_id").Get(&mySubEffectValue, &myMemberTagId)
		if ss.CheckErr(err) {
			return
		}

		exists, err = ss.MainEng.Table("unit_type_member_tag_m").
			Where("unit_type_id = ? AND member_tag_id = ?", unitTypeId, myMemberTagId).Exist()
		if ss.CheckErr(err) {
			return
		}
		var mySubSmile, mySubPure, mySubCool float64
		if exists {
			switch myAttrId {
			case 1:
				mySubSmile = math.Ceil(smileMax * (mySubEffectValue / 100))
			case 2:
				mySubPure = math.Ceil(pureMax * (mySubEffectValue / 100))
			case 3:
				mySubCool = math.Ceil(coolMax * (mySubEffectValue / 100))
			}
			// fmt.Println("副C技能属性加成:", mySubSmile, mySubPure, mySubCool)
		}

		// 好友主唱技能加成
		// TODO 好友支援存入数据库
		var tomoUnitId int64
		partyList := gjson.Parse(honokautils.ReadAllText("assets/serverdata/partylist.json")).Get("response_data.party_list")
		partyList.ForEach(func(key, value gjson.Result) bool {
			if value.Get("user_info.user_id").Int() == playReq.PartyUserID {
				tomoUnitId = value.Get("center_unit_info.unit_id").Int()
				return false
			}
			return true
		})
		// fmt.Println("好友UnitID:", tomoUnitId)

		// 好友主唱技能加成：主C技能（这里不使用新C技能，即以某属性的百分比提升另一属性）
		var tomoAttrId int
		var tomoEffectValue float64
		_, err = ss.MainEng.Table("unit_m").
			Join("LEFT", "unit_leader_skill_m", "unit_m.default_leader_skill_id = unit_leader_skill_m.unit_leader_skill_id").
			Where("unit_m.unit_id = ?", tomoUnitId).
			Cols("unit_m.attribute_id,unit_leader_skill_m.effect_value").Get(&tomoAttrId, &tomoEffectValue)
		if ss.CheckErr(err) {
			return
		}
		var tomoCenterSmile, tomoCenterPure, tomoCenterCool float64
		switch myAttrId {
		case 1:
			tomoCenterSmile = math.Ceil(smileMax * (tomoEffectValue / 100))
		case 2:
			tomoCenterPure = math.Ceil(pureMax * (tomoEffectValue / 100))
		case 3:
			tomoCenterCool = math.Ceil(coolMax * (tomoEffectValue / 100))
		}
		// fmt.Println("好友主C技能属性加成:", tomoCenterSmile, tomoCenterPure, tomoCenterCool)

		// 好友主唱技能加成：副C技能
		var tomoSubEffectValue float64
		var tomoMemberTagId int
		_, err = ss.MainEng.Table("unit_m").
			Join("LEFT", "unit_leader_skill_extra_m", "unit_m.default_leader_skill_id = unit_leader_skill_extra_m.unit_leader_skill_id").
			Where("unit_m.unit_id = ?", tomoUnitId).
			Cols("unit_leader_skill_extra_m.effect_value,unit_leader_skill_extra_m.member_tag_id").Get(&tomoSubEffectValue, &tomoMemberTagId)
		if ss.CheckErr(err) {
			return
		}

		exists, err = ss.MainEng.Table("unit_type_member_tag_m").
			Where("unit_type_id = ? AND member_tag_id = ?", unitTypeId, tomoMemberTagId).Exist()
		if ss.CheckErr(err) {
			return
		}
		var tomoSubSmile, tomoSubPure, tomoSubCool float64
		if exists {
			switch myAttrId {
			case 1:
				tomoSubSmile = math.Ceil(smileMax * (tomoSubEffectValue / 100))
			case 2:
				tomoSubPure = math.Ceil(pureMax * (tomoSubEffectValue / 100))
			case 3:
				tomoSubCool = math.Ceil(coolMax * (tomoSubEffectValue / 100))
			}
			// fmt.Println("好友副C技能属性加成:", tomoSubSmile, tomoSubPure, tomoSubCool)
		}

		// 全部卡属性
		totalSmile += smileMax + myCenterSmile + mySubSmile + tomoCenterSmile + tomoSubSmile
		totalPure += pureMax + myCenterPure + mySubPure + tomoCenterPure + tomoSubPure
		totalCool += coolMax + myCenterCool + mySubCool + tomoCenterCool + tomoSubCool
		totalHp += maxHp

		// 单卡属性计算结果取上取整
		fixedSmileMax := int(smileMax)
		fixedPureMax := int(pureMax)
		fixedCoolMax := int(coolMax)
		// fmt.Println("单卡属性:", fixedSmileMax, fixedPureMax, fixedCoolMax)

		unitList = append(unitList, live.UnitList{
			Smile: fixedSmileMax,
			Cute:  fixedPureMax,
			Cool:  fixedCoolMax,
		})
	}

	// 全部卡属性计算结果取上取整
	fixedTotalSmile := int(math.Ceil(totalSmile))
	fixedTotalPure := int(math.Ceil(totalPure))
	fixedTotalCool := int(math.Ceil(totalCool))
	// fmt.Println("全卡组属性:", fixedTotalSmile, fixedTotalPure, fixedTotalCool)

	lives := []live.PlayLiveList{}
	lives = append(lives, live.PlayLiveList{
		LiveInfo: live.LiveInfo{
			LiveDifficultyID: difficultyId,
			IsRandom:         false,
			AcFlag:           ac_flag,
			SwingFlag:        swing_flag,
			NotesList:        notes,
		},
		DeckInfo: live.DeckInfo{
			UnitDeckID:       deckId,
			TotalSmile:       fixedTotalSmile,
			TotalCute:        fixedTotalPure,
			TotalCool:        fixedTotalCool,
			TotalHp:          totalHp,
			PreparedHpDamage: 0,
			UnitList:         unitList,
		},
	})

	playResp := live.PlayResp{
		ResponseData: live.PlayData{
			RankInfo:            ranks,
			EnergyFullTime:      "2023-03-20 01:28:55",
			OverMaxEnergy:       0,
			AvailableLiveResume: false,
			LiveList:            lives,
			IsMarathonEvent:     false,
			MarathonEventID:     nil,
			NoSkill:             false,
			CanActivateEffect:   true,
			ServerTimestamp:     time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/play", middleware.Common, play)
}
