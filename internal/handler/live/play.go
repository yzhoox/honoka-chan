package live

import (
	"errors"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func play(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playReq := liveschema.PlayReq{}
	err := honokautils.ParseRequestData(ctx, &playReq)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println(ctx.MustGet("request_data").(string))

	if err := ss.ResetRandomLiveInProgress(); ss.CheckErr(err) {
		return
	}

	ss.RegisterLiveInProgress(playReq.UnitDeckID)

	playResp, err := BuildPlayResp(ss, playReq, false)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(playResp)
}

func BuildPlayResp(ss *session.Session, playReq liveschema.PlayReq, isRandom bool) (liveschema.PlayResp, error) {
	difficultyID, err := strconv.Atoi(playReq.LiveDifficultyID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	liveSetting, err := loadLiveSetting(ss, difficultyID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	notes, err := loadLiveNotes(liveSetting.NotesSettingAsset)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	ranks := buildRankInfo(liveSetting)

	owningIDList := []int{}
	err = ss.UserEng.Table("user_deck_unit").
		Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_deck.user_id = ? AND user_deck.deck_id = ?", ss.UserID, playReq.UnitDeckID).
		Cols("unit_owning_user_id").OrderBy("user_deck_unit.position ASC").Find(&owningIDList)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	if len(owningIDList) == 0 {
		return liveschema.PlayResp{}, errors.New("deck has no units")
	}

	museumBuff, err := getMuseumBuff(ss)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	myCenterUnitID, err := getDeckCenterUnitID(ss, playReq.UnitDeckID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	myLeaderSkill, err := getLeaderSkillData(ss, myCenterUnitID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	tomoUnitID, err := getSupportCenterUnitID(ss, int(playReq.PartyUserID), myCenterUnitID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	tomoLeaderSkill, err := getLeaderSkillData(ss, tomoUnitID)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	memberTagCache := map[memberTagMatchKey]bool{}
	unitDataMap, err := loadDeckUnitDataMap(ss, owningIDList)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	accessoryBonusMap, err := loadAccessoryBonusMap(ss, ss.UserID, owningIDList)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	removableSkillMap, err := loadDeckRemovableSkillMap(ss, ss.UserID, owningIDList)
	if err != nil {
		return liveschema.PlayResp{}, err
	}
	allRemovableSkillIDs := make([]int, 0)
	for _, skillIDs := range removableSkillMap {
		allRemovableSkillIDs = append(allRemovableSkillIDs, skillIDs...)
	}
	removableSkillMetaMap, err := loadRemovableSkillMetaMap(ss, allRemovableSkillIDs)
	if err != nil {
		return liveschema.PlayResp{}, err
	}

	unitList := make([]liveschema.UnitList, 0, len(owningIDList))
	var totalSmile, totalPure, totalCool float64
	var totalHP int
	for _, owningID := range owningIDList {
		// 卡片基础属性
		var baseSmile, basePure, baseCool, smileMax, pureMax, coolMax float64
		unitData, ok := unitDataMap[owningID]
		if !ok {
			return liveschema.PlayResp{}, errors.New("unit data not found in deck")
		}
		baseSmile = float64(unitData.Smile)
		basePure = float64(unitData.Cute)
		baseCool = float64(unitData.Cool)
		// fmt.Println("================================")
		// fmt.Println("基础属性:", baseSmile, basePure, baseCool)

		// 饰品属性加成（满级）
		if bonus, ok := accessoryBonusMap[owningID]; ok {
			// 饰品属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
			baseSmile += bonus.Smile
			basePure += bonus.Pure
			baseCool += bonus.Cool
			// fmt.Println("饰品属性加成:", smileAccessory, pureAccessory, coolAccessory)
			// fmt.Println("饰品属性加成后的基础属性:", baseSmile, basePure, baseCool)
		}

		// 回忆画廊属性加成（该加成会影响个宝等百分比宝石属性加成的计算，故先计算。）
		baseSmile += museumBuff.Smile
		basePure += museumBuff.Pure
		baseCool += museumBuff.Cool
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
		removableSkillIDs := removableSkillMap[owningID]
		for _, skillID := range removableSkillIDs {
			// 判断宝石效果类型（效果范围、效果类型、效果值、是否固定数值）
			meta, ok := removableSkillMetaMap[skillID]
			if !ok {
				return liveschema.PlayResp{}, errors.New("removable skill meta not found")
			}

			if meta.FixedValueFlag == 1 {
				// 吻、眼神属性加成（固定数值）
				switch meta.EffectType {
				case 1:
					kissSmile += meta.EffectValue
				case 2:
					kissPure += meta.EffectValue
				case 3:
					kissCool += meta.EffectValue
				}
				// fmt.Println("吻、眼神属性加成:", kissSmile, kissPure, kissCool)
				continue
			}

			// 仅效果类型为1、2、3的有属性加成
			if meta.EffectType != 1 && meta.EffectType != 2 && meta.EffectType != 3 {
				continue
			}

			// 加成范围：2:全员 1:非全员
			if meta.EffectRange == 2 || meta.TargetReferenceType == 1 || meta.TargetReferenceType == 2 {
				// refType: 1 -> 年级类加成, target_type -> 指定年级（这里不需要使用，因为能装上宝石肯定是符合的）
				// refType: 2 -> 个宝
				// refType: 3 -> 爆分、奶、判宝石, 0 -> 竞技场宝石
				switch meta.EffectType {
				case 1:
					skillSmile += math.Ceil(baseSmile * (meta.EffectValue / 100))
				case 2:
					skillPure += math.Ceil(basePure * (meta.EffectValue / 100))
				case 3:
					skillCool += math.Ceil(baseCool * (meta.EffectValue / 100))
				}
				// fmt.Println("年级类宝石、个宝属性加成:", skillSmile, skillPure, skillCool)
			}
		}

		// 单卡属性
		smileMax = baseSmile + kissSmile + skillSmile
		pureMax = basePure + kissPure + skillPure
		coolMax = baseCool + kissCool + skillCool

		// 主唱技能加成：主C技能（这里不使用新C技能，即以某属性的百分比提升另一属性）
		myCenterSmile, myCenterPure, myCenterCool := applyAttributeBonus(
			myLeaderSkill.AttributeID,
			myLeaderSkill.MainEffectValue,
			smileMax, pureMax, coolMax,
		)
		// fmt.Println("主C技能属性加成:", myCenterSmile, myCenterPure, myCenterCool)

		// 主唱技能加成：副C技能
		matched, err := matchMemberTag(ss, memberTagCache, unitData.UnitTypeID, myLeaderSkill.MemberTagID)
		if err != nil {
			return liveschema.PlayResp{}, err
		}

		var mySubSmile, mySubPure, mySubCool float64
		if matched {
			mySubSmile, mySubPure, mySubCool = applyAttributeBonus(
				myLeaderSkill.AttributeID,
				myLeaderSkill.ExtraEffectValue,
				smileMax, pureMax, coolMax,
			)
			// fmt.Println("副C技能属性加成:", mySubSmile, mySubPure, mySubCool)
		}

		// 好友主唱技能加成
		tomoCenterSmile, tomoCenterPure, tomoCenterCool := applyAttributeBonus(
			myLeaderSkill.AttributeID,
			tomoLeaderSkill.MainEffectValue,
			smileMax, pureMax, coolMax,
		)
		// fmt.Println("好友主C技能属性加成:", tomoCenterSmile, tomoCenterPure, tomoCenterCool)

		// 好友主唱技能加成：副C技能
		matched, err = matchMemberTag(ss, memberTagCache, unitData.UnitTypeID, tomoLeaderSkill.MemberTagID)
		if err != nil {
			return liveschema.PlayResp{}, err
		}

		var tomoSubSmile, tomoSubPure, tomoSubCool float64
		if matched {
			tomoSubSmile, tomoSubPure, tomoSubCool = applyAttributeBonus(
				myLeaderSkill.AttributeID,
				tomoLeaderSkill.ExtraEffectValue,
				smileMax, pureMax, coolMax,
			)
			// fmt.Println("好友副C技能属性加成:", tomoSubSmile, tomoSubPure, tomoSubCool)
		}

		// 全部卡属性
		totalSmile += smileMax + myCenterSmile + mySubSmile + tomoCenterSmile + tomoSubSmile
		totalPure += pureMax + myCenterPure + mySubPure + tomoCenterPure + tomoSubPure
		totalCool += coolMax + myCenterCool + mySubCool + tomoCenterCool + tomoSubCool
		totalHP += unitData.MaxHp
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
		unitList = append(unitList, liveschema.UnitList{
			Smile: int(smileMax),
			Cute:  int(pureMax),
			Cool:  int(coolMax),
		})
		// fmt.Println("单卡属性:", int(smileMax), int(pureMax), int(coolMax))
	}

	// 全部卡属性计算结果取上取整
	lives := []liveschema.PlayLiveList{
		{
			LiveInfo: liveschema.LiveInfo{
				LiveDifficultyID: difficultyID,
				IsRandom:         isRandom,
				AcFlag:           liveSetting.AcFlag,
				SwingFlag:        liveSetting.SwingFlag,
				NotesList:        notes,
			},
			DeckInfo: liveschema.DeckInfo{
				UnitDeckID:       playReq.UnitDeckID,
				TotalSmile:       int(math.Ceil(totalSmile)),
				TotalCute:        int(math.Ceil(totalPure)),
				TotalCool:        int(math.Ceil(totalCool)),
				TotalHp:          totalHP,
				PreparedHpDamage: 0,
				UnitList:         unitList,
			},
		},
	}
	// fmt.Println("全卡组属性:", int(math.Ceil(totalSmile)), int(math.Ceil(totalPure)), int(math.Ceil(totalCool)))

	return liveschema.PlayResp{
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
	}, nil
}

func init() {
	router.AddHandler("main.php", "POST", "/live/play", middleware.Common, play)
}
