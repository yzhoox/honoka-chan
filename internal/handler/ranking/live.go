package ranking

import (
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	profileapischema "honoka-chan/internal/schema/api/profile"
	rankingschema "honoka-chan/internal/schema/ranking"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

const livePageSize = 20

type liveRankingRow struct {
	UserID                 int    `xorm:"user_id"`
	HiScore                int    `xorm:"hi_score"`
	UserName               string `xorm:"user_name"`
	UserLevel              int    `xorm:"user_level"`
	AwardID                int    `xorm:"award_id"`
	CenterUnitOwningUserID int    `xorm:"center_unit_owning_user_id"`
}

func live(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := rankingschema.LiveReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	page := max(int(req.Page), 0)

	liveDifficultyID, err := strconv.Atoi(req.LiveDifficultyID)
	if ss.CheckErr(err) {
		return
	}

	totalCnt, err := ss.UserEng.Table(new(usermodel.UserLiveStatus)).
		Where("live_difficulty_id = ?", liveDifficultyID).
		Where("clear_cnt > 0").
		Count()
	if ss.CheckErr(err) {
		return
	}

	rows := []liveRankingRow{}
	if totalCnt > 0 {
		err = ss.UserEng.Table(new(usermodel.UserLiveStatus)).Alias("uls").
			Join("LEFT", "user_pref up", "up.user_id = uls.user_id").
			Join("LEFT", "user_deck ud", "ud.user_id = uls.user_id AND ud.main_flag = 1").
			Join("LEFT", "user_deck_unit udu", "udu.user_deck_id = ud.id AND udu.position = 5").
			Select(`
				uls.user_id,
				uls.hi_score,
				up.user_name,
				up.user_level,
				up.award_id,
				COALESCE(udu.unit_owning_user_id, up.unit_owning_user_id) AS center_unit_owning_user_id
			`).
			Where("uls.live_difficulty_id = ?", liveDifficultyID).
			Where("uls.clear_cnt > 0").
			OrderBy("uls.hi_score DESC, uls.user_id ASC").
			Limit(livePageSize, page*livePageSize).
			Find(&rows)
		if ss.CheckErr(err) {
			return
		}
	}

	items := make([]rankingschema.LiveItem, 0, len(rows))
	for i, row := range rows {
		centerUnitInfo, err := buildLiveCenterUnitInfo(ss, row.UserID, row.CenterUnitOwningUserID)
		if ss.CheckErr(err) {
			return
		}

		items = append(items, rankingschema.LiveItem{
			Rank:  page*livePageSize + i + 1,
			Score: row.HiScore,
			UserData: rankingschema.UserData{
				UserID: row.UserID,
				Name:   row.UserName,
				Level:  row.UserLevel,
			},
			CenterUnitInfo: centerUnitInfo,
			SettingAwardID: row.AwardID,
		})
	}

	ss.Respond(rankingschema.LiveResp{
		ResponseData: rankingschema.LiveData{
			Rank:       nil,
			Items:      items,
			TotalCnt:   totalCnt,
			PresentCnt: 0,
			Page:       page,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func buildLiveCenterUnitInfo(ss *session.Session, targetUserID, unitOwningUserID int) (profileapischema.CenterUnitInfo, error) {
	info := profileapischema.CenterUnitInfo{
		RemovableSkillIds: []int{},
	}
	if unitOwningUserID <= 0 {
		return info, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", targetUserID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return info, err
	}
	if !has {
		return info, nil
	}

	accessoryInfo, err := getRankingAccessoryInfo(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return info, err
	}

	removableSkillIDs := []int{}
	err = ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
		Where("user_id = ? AND unit_owning_user_id = ?", targetUserID, unitOwningUserID).
		Cols("unit_removable_skill_id").
		Find(&removableSkillIDs)
	if err != nil {
		return info, err
	}

	return profileapischema.CenterUnitInfo{
		UnitOwningUserID:           unitData.UnitOwningUserID,
		UnitID:                     unitData.UnitID,
		Exp:                        unitData.Exp,
		NextExp:                    0,
		Level:                      unitData.Level,
		LevelLimitID:               unitData.LevelLimitID,
		MaxLevel:                   unitData.MaxLevel,
		Love:                       unitData.Love,
		Rank:                       unitData.Rank,
		MaxRank:                    unitData.MaxRank,
		MaxLove:                    unitData.MaxLove,
		UnitSkillLevel:             unitData.UnitSkillLevel,
		MaxHp:                      unitData.MaxHp,
		FavoriteFlag:               unitData.FavoriteFlag,
		DisplayRank:                unitData.DisplayRank,
		Attribute:                  unitData.Attribute,
		Smile:                      unitData.Smile,
		Cute:                       unitData.Cute,
		Cool:                       unitData.Cool,
		IsLoveMax:                  unitData.IsLoveMax,
		IsRankMax:                  unitData.IsRankMax,
		IsLevelMax:                 unitData.IsLevelMax,
		IsSigned:                   unitData.IsSigned,
		IsSkillLevelMax:            unitData.IsSkillLevelMax,
		UnitSkillExp:               unitData.UnitSkillExp,
		RemovableSkillIds:          removableSkillIDs,
		UnitRemovableSkillCapacity: unitData.UnitRemovableSkillCapacity,
		AccessoryInfo:              accessoryInfo,
		Costume: profileapischema.Costume{
			UnitID:    unitData.UnitID,
			IsRankMax: unitData.IsRankMax,
			IsSigned:  unitData.IsSigned,
		},
		TotalSmile: unitData.Smile,
		TotalCute:  unitData.Cute,
		TotalCool:  unitData.Cool,
		TotalHp:    unitData.MaxHp,
	}, nil
}

func getRankingAccessoryInfo(ss *session.Session, targetUserID, unitOwningUserID int) (*profileapischema.AccessoryInfo, error) {
	if unitOwningUserID <= 0 {
		return nil, nil
	}

	var accessoryOwningID int
	has, err := ss.UserEng.Table("user_accessory_wear").
		Where("unit_owning_user_id = ? AND user_id = ?", unitOwningUserID, targetUserID).
		Cols("accessory_owning_user_id").
		Get(&accessoryOwningID)
	if err != nil {
		return nil, err
	}
	if !has || accessoryOwningID <= 0 {
		return nil, nil
	}

	has, accessoryData := ss.GetAccessoryByAccessoryOwningUserID(accessoryOwningID)
	if !has {
		return nil, nil
	}

	return &profileapischema.AccessoryInfo{
		AccessoryOwningUserID: accessoryOwningID,
		AccessoryID:           accessoryData.AccessoryID,
		Exp:                   accessoryData.Exp,
		NextExp:               0,
		Level:                 accessoryData.MaxLevel,
		MaxLevel:              accessoryData.MaxLevel,
		RankUpCount:           accessoryData.MaxLevel - accessoryData.DefaultMaxLevel,
		FavoriteFlag:          true,
	}, nil
}

func init() {
	router.AddHandler("main.php", "POST", "/ranking/live", middleware.Common, live)
}
