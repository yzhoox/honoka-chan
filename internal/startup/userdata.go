package startup

import (
	"fmt"
	"honoka-chan/internal/model/user"
	"honoka-chan/internal/schema/api/profile"
	"honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"strconv"
	"time"
)

func InitUserData(userId int) {
	userList := []user.Users{}
	if userId != 0 {
		err := db.UserEng.Table("users").Where("user_id = ?", userId).Find(&userList)
		utils.CheckErr(err)
	} else {
		err := db.UserEng.Table("users").Asc("id").Find(&userList)
		utils.CheckErr(err)
	}

	// 同步用户表
	// db.UserEng.ShowSQL(true)
	// db.UserEng.Table("user_info_m").Sync2(new(user.UserInfo))
	// db.UserEng.Table("user_pref").Sync2(new(UserPref))

	// 开始会话
	session := db.UserEng.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		panic(err)
	}

	for _, u := range userList {
		// 检查用户配置
		exists, err := db.UserEng.Table("user_pref").Where("user_id = ?", u.UserID).Exist()
		utils.CheckErr(err)

		if !exists {
			// 默认中心成员（）
			var oId int
			_, err = db.MainEng.Table("common_unit_m").Cols("unit_owning_user_id").Where("unit_id = ?", 31).Get(&oId)
			utils.CheckErr(err)
			fmt.Println("Center UnitOwningUserId:", oId)
			userPref := user.UserPref{
				UserID:           u.UserID,
				AwardID:          1, // 音乃木坂学生
				BackgroundID:     1, // 初始背景
				UnitOwningUserID: oId,
				UserName:         "音乃木坂学生",
				UserLevel:        1,
				UserDesc:         "你好。",
				InviteCode:       strconv.Itoa(u.UserID),
				UpdateTime:       time.Now().Unix(),
			}
			_, err = session.Table("user_pref").Insert(&userPref)
			utils.CheckErr(err)
			// fmt.Println("UserPref Id", userPref.Id)
		}

		// 检查用户信息
		exists, err = db.UserEng.Table("user_info_m").Where("user_id = ?", u.UserID).Exist()

		// 检查用户卡组配置
		exists, err = db.UserEng.Table("user_deck").Where("user_id = ?", u.UserID).Asc("deck_id").Exist()
		utils.CheckErr(err)
		// fmt.Println("UserDeck exists:", exists)

		if !exists {
			userDeck := unit.UserDeckData{
				DeckID:     1,
				MainFlag:   1,
				DeckName:   "队伍A",
				UserID:     u.UserID,
				InsertDate: time.Now().Unix(),
			}

			// 默认队伍
			_, err = session.Table("user_deck").Insert(&userDeck)
			if err != nil {
				session.Rollback()
				panic(err)
			}
			userDeckId := userDeck.ID
			// fmt.Println("New UserDeck:", userDeckId)

			// 默认卡组
			unitIds := []int{}
			err = db.MainEng.Table("unit_m").Cols("unit_id").Where("album_series_id = ?", 615).Find(&unitIds)
			if err != nil {
				session.Rollback()
				panic(err)
			}

			unitData := []profile.UnitData{}
			err = db.MainEng.Table("common_unit_m").In("unit_id", unitIds).Find(&unitData)
			if err != nil {
				session.Rollback()
				panic(err)
			}
			// fmt.Println(unitData)

			position := 1
			for _, u := range unitData {
				unitDeckData := unit.UnitDeckData{
					UserDeckID:       userDeckId,
					UnitOwningUserID: u.UnitOwningUserID,
					UnitID:           u.UnitID,
					Position:         position,
					Level:            100,
					LevelLimitID:     1,
					DisplayRank:      2,
					Love:             1000,
					UnitSkillLevel:   8,
					IsRankMax:        true,
					IsLoveMax:        true,
					IsLevelMax:       true,
					IsSigned:         u.IsSigned,
					BeforeLove:       1000,
					MaxLove:          1000,
					InsertData:       time.Now().Unix(),
				}
				_, err = session.Table("user_deck_unit").Insert(&unitDeckData)
				if err != nil {
					session.Rollback()
					panic(err)
				}
				// fmt.Println("New DeckUnit:", unitDeckData.Id)

				position++
			}

			// 结束会话
			if err = session.Commit(); err != nil {
				panic(err)
			}
		}
	}
}
