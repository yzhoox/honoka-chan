package friend

import (
	"errors"
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	friendschema "honoka-chan/internal/schema/friend"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type friendSearchRow struct {
	UserID                 int    `xorm:"user_id"`
	InviteCode             string `xorm:"invite_code"`
	UserName               string `xorm:"user_name"`
	UserLevel              int    `xorm:"user_level"`
	UserDesc               string `xorm:"user_desc"`
	EnergyMax              int    `xorm:"energy_max"`
	AwardID                int    `xorm:"award_id"`
	LastLoginTime          int64  `xorm:"last_login_time"`
	CenterUnitOwningUserID int    `xorm:"center_unit_owning_user_id"`
	NaviUnitOwningUserID   int    `xorm:"navi_unit_owning_user_id"`
}

func search(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := friendschema.SearchReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	targetUserID, err := resolveUserIDByInviteCode(ss, req.InviteCode)
	if ss.CheckErr(err) {
		return
	}

	row := friendSearchRow{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).Alias("up").
		Join("LEFT", "users u", "u.user_id = up.user_id").
		Join("LEFT", "user_deck ud", "ud.user_id = up.user_id AND ud.main_flag = 1").
		Join("LEFT", "user_deck_unit udu", "udu.user_deck_id = ud.id AND udu.position = 5").
		Select(`
			up.user_id,
			up.invite_code,
			up.user_name,
			up.user_level,
			up.user_desc,
			up.energy_max,
			up.award_id,
			u.last_login_time,
			COALESCE(udu.unit_owning_user_id, up.unit_owning_user_id) AS center_unit_owning_user_id,
			up.unit_owning_user_id AS navi_unit_owning_user_id
		`).
		Where("up.user_id = ?", targetUserID).
		Get(&row)
	if ss.CheckErr(err) {
		return
	}
	if !has {
		ss.CheckErr(errors.New("user not found"))
		return
	}

	unitCount, err := ss.UserEng.Table(new(usermodel.UserUnitData)).
		Where("user_id = ?", targetUserID).
		Count()
	if ss.CheckErr(err) {
		return
	}

	centerUnitInfo, err := buildCenterUnitInfo(ss, targetUserID, row.CenterUnitOwningUserID, row.AwardID)
	if ss.CheckErr(err) {
		return
	}

	costume, err := getSearchCostume(ss, targetUserID, row.NaviUnitOwningUserID, centerUnitInfo)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(friendschema.SearchResp{
		ResponseData: friendschema.SearchData{
			UserInfo: friendschema.SearchUserInfo{
				UserID:               row.UserID,
				Name:                 row.UserName,
				Level:                row.UserLevel,
				CostMax:              100,
				UnitMax:              5000,
				EnergyMax:            positiveIntOrDefault(row.EnergyMax, usermodel.DefaultUserEnergyMax),
				FriendMax:            99,
				UnitCnt:              int(unitCount),
				ElapsedTimeFromLogin: formatElapsedTime(row.LastLoginTime),
				Comment:              row.UserDesc,
			},
			CenterUnitInfo:  toSearchCenterUnitInfo(centerUnitInfo, costume),
			SettingAwardID:  row.AwardID,
			IsAlliance:      false,
			FriendStatus:    0,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/search", middleware.Common, search)
}

func positiveIntOrDefault(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func toSearchCenterUnitInfo(info friendschema.CenterUnitInfo, costume friendschema.SearchCostume) friendschema.SearchCenterUnitInfo {
	return friendschema.SearchCenterUnitInfo{
		UnitOwningUserID:           info.UnitOwningUserID,
		UnitID:                     info.UnitID,
		Exp:                        info.Exp,
		NextExp:                    info.NextExp,
		Level:                      info.Level,
		LevelLimitID:               info.LevelLimitID,
		MaxLevel:                   info.MaxLevel,
		Rank:                       info.Rank,
		MaxRank:                    info.MaxRank,
		Love:                       info.Love,
		MaxLove:                    info.MaxLove,
		UnitSkillLevel:             info.UnitSkillLevel,
		MaxHp:                      info.MaxHp,
		FavoriteFlag:               info.FavoriteFlag,
		DisplayRank:                info.DisplayRank,
		UnitSkillExp:               info.UnitSkillExp,
		UnitRemovableSkillCapacity: info.UnitRemovableSkillCapacity,
		Attribute:                  info.Attribute,
		Smile:                      info.Smile,
		Cute:                       info.Cute,
		Cool:                       info.Cool,
		IsLoveMax:                  info.IsLoveMax,
		IsLevelMax:                 info.IsLevelMax,
		IsRankMax:                  info.IsRankMax,
		IsSigned:                   info.IsSigned,
		IsSkillLevelMax:            info.IsSkillLevelMax,
		SettingAwardID:             info.SettingAwardID,
		RemovableSkillIds:          info.RemovableSkillIds,
		AccessoryInfo:              info.AccessoryInfo,
		Costume:                    costume,
	}
}

func getSearchCostume(ss *session.Session, targetUserID, unitOwningUserID int, fallback friendschema.CenterUnitInfo) (friendschema.SearchCostume, error) {
	costume := friendschema.SearchCostume{
		UnitID:    fallback.UnitID,
		IsRankMax: fallback.IsRankMax,
		IsSigned:  fallback.IsSigned,
	}
	if unitOwningUserID <= 0 {
		return costume, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", targetUserID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return costume, err
	}
	if !has {
		return costume, nil
	}

	return friendschema.SearchCostume{
		UnitID:    unitData.UnitID,
		IsRankMax: unitData.IsRankMax,
		IsSigned:  unitData.IsSigned,
	}, nil
}
