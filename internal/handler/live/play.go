package live

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func play(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playReq := liveschema.PlayReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &playReq)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println(ctx.MustGet("request_data").(string))

	ss.RegisterLiveInProgress(playReq.UnitDeckID)

	difficultyID, _ := strconv.Atoi(playReq.LiveDifficultyID)

	// 歌曲类型: normal / special
	// sqlite3 不支持 FULL OUTER JOIN 所以这里使用 UNION ALL
	var liveSetting struct {
		NotesSettingAsset string `xorm:"notes_setting_asset"`
		ARankScore        int    `xorm:"a_rank_score"`
		BRankScore        int    `xorm:"b_rank_score"`
		CRankScore        int    `xorm:"c_rank_score"`
		SRankScore        int    `xorm:"s_rank_score"`
		AcFlag            int    `xorm:"ac_flag"`
		SwingFlag         int    `xorm:"swing_flag"`
	}
	sql := `
		SELECT notes_setting_asset,
			a_rank_score,
			b_rank_score,
			c_rank_score,
			s_rank_score,
			ac_flag,
			swing_flag
		FROM live_setting_m
		WHERE live_setting_id IN (
			SELECT live_setting_id
			FROM normal_live_m
			WHERE live_difficulty_id = ?
			UNION ALL
			SELECT live_setting_id
			FROM special_live_m
			WHERE live_difficulty_id = ?
		)
		`
	_, err = ss.MainEng.SQL(sql, difficultyID, difficultyID).Get(&liveSetting)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println("liveSetting", liveSetting)

	notes := []liveschema.NotesList{}
	notes_list := honokautils.ReadAllText("./assets/serverdata/beatmaps/" + liveSetting.NotesSettingAsset)
	err = json.Unmarshal([]byte(notes_list), &notes)
	if ss.CheckErr(err) {
		return
	}

	ranks := []liveschema.RankInfo{}
	ranks = append(ranks, liveschema.RankInfo{
		Rank:    5,
		RankMin: 0,
		RankMax: liveSetting.CRankScore,
	}, liveschema.RankInfo{
		Rank:    4,
		RankMin: liveSetting.CRankScore + 1,
		RankMax: liveSetting.BRankScore,
	}, liveschema.RankInfo{
		Rank:    3,
		RankMin: liveSetting.BRankScore + 1,
		RankMax: liveSetting.ARankScore,
	}, liveschema.RankInfo{
		Rank:    2,
		RankMin: liveSetting.ARankScore + 1,
		RankMax: liveSetting.SRankScore,
	}, liveschema.RankInfo{
		Rank:    1,
		RankMin: liveSetting.SRankScore + 1,
		RankMax: 0,
	})

	owningIdList := []int{}
	err = ss.UserEng.Table("user_deck_unit").
		Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_deck.user_id = ? AND user_deck.deck_id = ?", ss.UserID, playReq.UnitDeckID).
		Cols("unit_owning_user_id").OrderBy("user_deck_unit.position ASC").Find(&owningIdList)
	if ss.CheckErr(err) {
		return
	}

	unitList := []liveschema.UnitList{}
	var totalSmile, totalPure, totalCool float64
	var totalHp int
	for _, owningId := range owningIdList {
		// 卡片基础属性
		var baseSmile, basePure, baseCool, smileMax, pureMax, coolMax float64
		var unitData unitmodel.UnitDataMap
		_, err = ss.GetBasicUnitInfo().Where("unit_owning_user_id = ?", owningId).Get(&unitData)
		if ss.CheckErr(err) {
			return
		}
		baseSmile = float64(unitData.Smile)
		basePure = float64(unitData.Cute)
		baseCool = float64(unitData.Cool)
		// fmt.Println("================================")
		// fmt.Println("基础属性:", baseSmile, basePure, baseCool)

		// 饰品属性加成（满级）
		var accessoryOwningId int
		exists, err := ss.UserEng.Table("user_accessory_wear").
			Where("unit_owning_user_id = ?", owningId).
			Cols("accessory_owning_user_id").Get(&accessoryOwningId)
		if ss.CheckErr(err) {
			return
		}
		if exists {
			var smileAccessory, pureAccessory, coolAccessory float64
			_, err = ss.MainEng.Table("common_accessory_m").
				Join("LEFT", "accessory_m", "common_accessory_m.accessory_id = accessory_m.accessory_id").
				Where("accessory_owning_user_id = ?", accessoryOwningId).Cols("smile_max,pure_max,cool_max").
				Get(&smileAccessory, &pureAccessory, &coolAccessory)
			if ss.CheckErr(err) {
				return
			}

			// 饰品属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
			baseSmile += smileAccessory
			basePure += pureAccessory
			baseCool += coolAccessory
			// fmt.Println("饰品属性加成:", smileAccessory, pureAccessory, coolAccessory)
			// fmt.Println("饰品属性加成后的基础属性:", baseSmile, basePure, baseCool)
		}

		// 回忆画廊属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		var smileBuff, pureBuff, coolBuff float64
		_, err = ss.MainEng.Table("museum_contents_m").
			Select("SUM(smile_buff),SUM(pure_buff),SUM(cool_buff)").
			Get(&smileBuff, &pureBuff, &coolBuff)
		if ss.CheckErr(err) {
			return
		}
		baseSmile += smileBuff
		basePure += pureBuff
		baseCool += coolBuff
		// fmt.Println("回忆画廊属性加成:", smileBuff, pureBuff, coolBuff)
		// fmt.Println("回忆画廊属性加成后的基础属性:", baseSmile, basePure, baseCool)

		// 绊属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		switch unitData.Attribute {
		case 1:
			baseSmile += float64(unitData.MaxLove)
		case 2:
			basePure += float64(unitData.MaxLove)
		case 3:
			baseCool += float64(unitData.MaxLove)
		}
		// fmt.Println("绊属性加成:", unitData.MaxLove)
		// fmt.Println("绊属性加成后的基础属性:", baseSmile, basePure, baseCool)

		// 宝石属性加成
		var kissSmile, kissPure, kissCool float64
		var skillSmile, skillPure, skillCool float64

		// 宝石加成（满级）
		removableSkillIds := []int{}
		err = ss.UserEng.Table("user_unit_skill_equip").
			Where("unit_owning_user_id = ?", owningId).
			Cols("unit_removable_skill_id").Find(&removableSkillIds)
		if ss.CheckErr(err) {
			return
		}

		for _, sk := range removableSkillIds {
			// 判断宝石效果类型（效果范围、效果类型、效果值、是否固定数值）
			var effectRange, effectType, fixedValueFlag, refType int
			var effectValue float64
			_, err = ss.MainEng.Table("unit_removable_skill_m").
				Where("unit_removable_skill_id = ?", sk).
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
		_, err = ss.UserEng.Table("user_deck_unit").
			Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
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
			Cols("unit_leader_skill_extra_m.effect_value,unit_leader_skill_extra_m.member_tag_id").
			Get(&mySubEffectValue, &myMemberTagId)
		if ss.CheckErr(err) {
			return
		}

		exists, err = ss.MainEng.Table("unit_type_member_tag_m").
			Where("unit_type_id = ? AND member_tag_id = ?", unitData.UnitTypeID, myMemberTagId).Exist()
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
		// var tomoUnitId int64
		// partyList := gjson.Parse(honokautils.ReadAllText("assets/serverdata/partylist.json")).Get("response_data.party_list")
		// partyList.ForEach(func(key, value gjson.Result) bool {
		// 	if value.Get("user_info.user_id").Int() == playReq.PartyUserID {
		// 		tomoUnitId = value.Get("center_unit_info.unit_id").Int()
		// 		return false
		// 	}
		// 	return true
		// })
		// TODO: 好友功能实装前先使用自己的卡组助战
		tomoUnitId := myCenterUnitId
		// fmt.Println("好友UnitID:", tomoUnitId)

		// 好友主唱技能加成：主C技能（这里不使用新C技能，即以某属性的百分比提升另一属性）
		var tomoEffectValue float64
		_, err = ss.MainEng.Table("unit_m").
			Join("LEFT", "unit_leader_skill_m", "unit_m.default_leader_skill_id = unit_leader_skill_m.unit_leader_skill_id").
			Where("unit_m.unit_id = ?", tomoUnitId).
			Cols("unit_leader_skill_m.effect_value").
			Get(&tomoEffectValue)
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
			Cols("unit_leader_skill_extra_m.effect_value,unit_leader_skill_extra_m.member_tag_id").
			Get(&tomoSubEffectValue, &tomoMemberTagId)
		if ss.CheckErr(err) {
			return
		}

		exists, err = ss.MainEng.Table("unit_type_member_tag_m").
			Where("unit_type_id = ? AND member_tag_id = ?", unitData.UnitTypeID, tomoMemberTagId).Exist()
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
		totalHp += unitData.MaxHp
		// fmt.Println("smileMax, myCenterSmile, mySubSmile, tomoCenterSmile, tomoSubSmile, totalSmile",
		// 	smileMax, myCenterSmile, mySubSmile, tomoCenterSmile, tomoSubSmile,
		// 	smileMax+myCenterSmile+mySubSmile+tomoCenterSmile+tomoSubSmile)
		// fmt.Println("pureMax, myCenterPure, mySubPure, tomoCenterPure, tomoSubPure, totalPure",
		// 	pureMax, myCenterPure, mySubPure, tomoCenterPure, tomoSubPure,
		// 	pureMax+myCenterPure+mySubPure+tomoCenterPure+tomoSubPure)
		// fmt.Println("coolMax, myCenterCool, mySubCool, tomoCenterCool, tomoSubCool, totalPure",
		// 	coolMax, myCenterCool, mySubCool, tomoCenterCool, tomoSubCool,
		// 	coolMax+myCenterCool+mySubCool+tomoCenterCool+tomoSubCool)

		// 单卡属性计算结果取上取整
		fixedSmileMax := int(smileMax)
		fixedPureMax := int(pureMax)
		fixedCoolMax := int(coolMax)
		// fmt.Println("单卡属性:", fixedSmileMax, fixedPureMax, fixedCoolMax)

		unitList = append(unitList, liveschema.UnitList{
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

	lives := []liveschema.PlayLiveList{}
	lives = append(lives, liveschema.PlayLiveList{
		LiveInfo: liveschema.LiveInfo{
			LiveDifficultyID: difficultyID,
			IsRandom:         false,
			AcFlag:           liveSetting.AcFlag,
			SwingFlag:        liveSetting.SwingFlag,
			NotesList:        notes,
		},
		DeckInfo: liveschema.DeckInfo{
			UnitDeckID:       playReq.UnitDeckID,
			TotalSmile:       fixedTotalSmile,
			TotalCute:        fixedTotalPure,
			TotalCool:        fixedTotalCool,
			TotalHp:          totalHp,
			PreparedHpDamage: 0,
			UnitList:         unitList,
		},
	})

	playResp := liveschema.PlayResp{
		ResponseData: liveschema.PlayData{
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
