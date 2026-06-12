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
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type friendListRow struct {
	FriendUserID           int    `xorm:"friend_user_id"`
	Status                 int    `xorm:"status"`
	IsNew                  bool   `xorm:"is_new"`
	InsertDate             int64  `xorm:"insert_date"`
	UpdateDate             int64  `xorm:"update_date"`
	InviteCode             string `xorm:"invite_code"`
	UserName               string `xorm:"user_name"`
	UserLevel              int    `xorm:"user_level"`
	UserDesc               string `xorm:"user_desc"`
	AwardID                int    `xorm:"award_id"`
	LastLoginTime          int64  `xorm:"last_login_time"`
	CenterUnitOwningUserID int    `xorm:"center_unit_owning_user_id"`
}

func list(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	listReq := friendschema.ListReq{}
	err := honokautils.ParseRequestData(ctx, &listReq)
	if ss.CheckErr(err) {
		return
	}

	friendStatus, err := mapFriendListTypeToStatus(listReq.Type)
	if ss.CheckErr(err) {
		return
	}

	totalCount, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("status = ?", friendStatus).
		Count()
	if ss.CheckErr(err) {
		return
	}

	page := max(listReq.Page, 0)
	offset := page * usermodel.FriendListPageSize

	rows := []friendListRow{}
	err = ss.UserEng.Table(new(usermodel.UserFriend)).Alias("uf").
		Join("LEFT", "user_pref up", "up.user_id = uf.friend_user_id").
		Join("LEFT", "users u", "u.user_id = uf.friend_user_id").
		Join("LEFT", "user_deck ud", "ud.user_id = uf.friend_user_id AND ud.main_flag = 1").
		Join("LEFT", "user_deck_unit udu", "udu.user_deck_id = ud.id AND udu.position = 5").
		Join("LEFT", "common_unit_data cud", "cud.unit_id = udu.unit_id").
		Where("uf.user_id = ?", ss.UserID).
		Where("uf.status = ?", friendStatus).
		Select(`
			uf.friend_user_id,
			uf.status,
			uf.is_new,
			uf.insert_date,
			uf.update_date,
			up.invite_code,
			up.user_name,
			up.user_level,
			up.user_desc,
			up.award_id,
			u.last_login_time,
			COALESCE(udu.unit_owning_user_id, up.unit_owning_user_id) AS center_unit_owning_user_id
		`).
		OrderBy(friendListOrderBy(listReq.Sort)).
		Limit(usermodel.FriendListPageSize, offset).
		Find(&rows)
	if ss.CheckErr(err) {
		return
	}

	friendList := make([]friendschema.FriendList, 0, len(rows))
	for _, row := range rows {
		item, err := buildFriendListItem(ss, row)
		if ss.CheckErr(err) {
			return
		}
		friendList = append(friendList, item)
	}

	if len(rows) > 0 {
		switch friendStatus {
		case usermodel.FriendStatusApproved, usermodel.FriendStatusAwaitingApproval:
			_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
				Where("user_id = ?", ss.UserID).
				Where("status = ?", friendStatus).
				Where("is_new = ?", true).
				Cols("is_new", "update_date").
				Update(&usermodel.UserFriend{
					IsNew:      false,
					UpdateDate: time.Now().Unix(),
				})
			if ss.CheckErr(err) {
				return
			}
		}
	}

	ss.Respond(friendschema.ListResp{
		ResponseData: friendschema.ListData{
			ItemCount:       totalCount,
			FriendList:      friendList,
			NewFriendList:   []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func buildFriendListItem(ss *session.Session, row friendListRow) (friendschema.FriendList, error) {
	appliedAt := max(row.UpdateDate, row.InsertDate)

	centerUnitInfo, err := buildCenterUnitInfo(ss, row.FriendUserID, row.CenterUnitOwningUserID, row.AwardID)
	if err != nil {
		return friendschema.FriendList{}, err
	}

	return friendschema.FriendList{
		UserData: friendschema.UserData{
			UserID:                 displayFriendUserID(row.InviteCode, row.FriendUserID),
			Name:                   row.UserName,
			Level:                  row.UserLevel,
			ElapsedTimeFromLogin:   formatElapsedTime(row.LastLoginTime),
			ElapsedTimeFromApplied: formatElapsedTime(appliedAt),
			Comment:                row.UserDesc,
		},
		CenterUnitInfo: centerUnitInfo,
		SettingAwardID: row.AwardID,
		IsNew:          row.IsNew,
	}, nil
}

func buildCenterUnitInfo(ss *session.Session, userID, unitOwningUserID, awardID int) (friendschema.CenterUnitInfo, error) {
	info := friendschema.CenterUnitInfo{
		SettingAwardID:    awardID,
		RemovableSkillIds: []int{},
	}
	if unitOwningUserID <= 0 {
		return info, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", userID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return info, err
	}
	if !has {
		return info, nil
	}

	var accessoryInfo friendschema.AccessoryInfo
	accessoryWear := usermodel.UserAccessoryWear{}
	has, err = ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ? AND unit_owning_user_id = ?", userID, unitOwningUserID).
		Get(&accessoryWear)
	if err != nil {
		return info, err
	}
	if has && accessoryWear.AccessoryOwningUserID > 0 {
		hasAccessory, accessoryData := ss.GetAccessoryByAccessoryOwningUserID(accessoryWear.AccessoryOwningUserID)
		if hasAccessory {
			accessoryInfo = friendschema.AccessoryInfo{
				AccessoryOwningUserID: accessoryWear.AccessoryOwningUserID,
				AccessoryID:           accessoryData.AccessoryID,
				Exp:                   accessoryData.Exp,
				NextExp:               0,
				Level:                 accessoryData.MaxLevel,
				MaxLevel:              accessoryData.MaxLevel,
				RankUpCount:           accessoryData.MaxLevel - accessoryData.DefaultMaxLevel,
				FavoriteFlag:          true,
			}
		}
	}

	removableSkillIDs := []int{}
	err = ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
		Where("user_id = ? AND unit_owning_user_id = ?", userID, unitOwningUserID).
		Cols("unit_removable_skill_id").
		Find(&removableSkillIDs)
	if err != nil {
		return info, err
	}

	info = friendschema.CenterUnitInfo{
		UnitOwningUserID:           int64(unitData.UnitOwningUserID),
		UnitID:                     unitData.UnitID,
		Exp:                        unitData.Exp,
		NextExp:                    0,
		Level:                      unitData.Level,
		LevelLimitID:               unitData.LevelLimitID,
		MaxLevel:                   unitData.MaxLevel,
		Rank:                       unitData.Rank,
		MaxRank:                    unitData.MaxRank,
		Love:                       unitData.Love,
		MaxLove:                    unitData.MaxLove,
		UnitSkillLevel:             unitData.UnitSkillLevel,
		MaxHp:                      unitData.MaxHp,
		FavoriteFlag:               unitData.FavoriteFlag,
		DisplayRank:                unitData.DisplayRank,
		UnitSkillExp:               unitData.UnitSkillExp,
		UnitRemovableSkillCapacity: unitData.UnitRemovableSkillCapacity,
		Attribute:                  unitData.Attribute,
		Smile:                      unitData.Smile,
		Cute:                       unitData.Cute,
		Cool:                       unitData.Cool,
		IsLoveMax:                  unitData.IsLoveMax,
		IsLevelMax:                 unitData.IsLevelMax,
		IsRankMax:                  unitData.IsRankMax,
		IsSigned:                   unitData.IsSigned,
		IsSkillLevelMax:            unitData.IsSkillLevelMax,
		SettingAwardID:             awardID,
		RemovableSkillIds:          removableSkillIDs,
		AccessoryInfo:              accessoryInfo,
	}

	return info, nil
}

func mapFriendListTypeToStatus(listType int) (int, error) {
	switch listType {
	case 0:
		return usermodel.FriendStatusApproved, nil
	case 1:
		return usermodel.FriendStatusPending, nil
	case 2:
		return usermodel.FriendStatusAwaitingApproval, nil
	default:
		return 0, errors.New("invalid friend list type")
	}
}

func friendListOrderBy(sort int) string {
	switch sort {
	case 1:
		return "up.user_level ASC"
	case 2:
		return "up.user_level DESC"
	case 3:
		return "cud.smile ASC"
	case 4:
		return "cud.smile DESC"
	case 5:
		return "cud.cool ASC"
	case 6:
		return "cud.cool DESC"
	case 7:
		return "cud.cute ASC"
	case 8:
		return "cud.cute DESC"
	case 9:
		return "u.last_login_time ASC"
	case 10:
		return "u.last_login_time DESC"
	case 11:
		return "uf.insert_date ASC"
	case 12:
		return "uf.insert_date DESC"
	default:
		return "uf.insert_date DESC"
	}
}

func displayFriendUserID(inviteCode string, fallback int) int {
	if id, err := strconv.Atoi(inviteCode); err == nil && id > 0 {
		return id
	}
	return fallback
}

func formatElapsedTime(ts int64) string {
	if ts <= 0 {
		return "刚刚"
	}

	delta := time.Since(time.Unix(ts, 0))
	if delta < time.Minute {
		return "刚刚"
	}
	if delta < time.Hour {
		return strconv.Itoa(int(delta.Minutes())) + "分钟前"
	}
	if delta < 24*time.Hour {
		return strconv.Itoa(int(delta.Hours())) + "小时前"
	}
	return strconv.Itoa(int(delta.Hours())/24) + "天前"
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/list", middleware.Common, list)
}
