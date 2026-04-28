package startup

import (
	ghomemodel "honoka-chan/internal/model/ghome"
	loginmodel "honoka-chan/internal/model/login"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/pkg/db"
	"log"
	"time"

	"xorm.io/xorm"
)

var (
	userEng *xorm.Session
)

func CreateTables() {
	// db.UserEng.ShowSQL(true)
	db.UserEng.Sync2(new(ghomemodel.DeviceKey))
	db.UserEng.Sync2(new(loginmodel.AuthKey))
	db.UserEng.Sync2(new(usermodel.UserAccessoryWear))
	db.UserEng.Sync2(new(usermodel.UserDeck))
	db.UserEng.Sync2(new(usermodel.UserDeckUnit))
	db.UserEng.Sync2(new(usermodel.UserKey))
	db.UserEng.Sync2(new(usermodel.UserLiveGoal))
	db.UserEng.Sync2(new(usermodel.UserLiveInProgress))
	db.UserEng.Sync2(new(usermodel.UserLiveRecord))
	db.UserEng.Sync2(new(usermodel.UserPref))
	db.UserEng.Sync2(new(usermodel.Users))
	db.UserEng.Sync2(new(usermodel.UserUnit))
	db.UserEng.Sync2(new(usermodel.UserUnitSkillEquip))
}

func LoadUnitData() {
	userEng = db.UserEng.NewSession()
	defer userEng.Close()

	err := userEng.Begin()
	CheckErr(err)

	commonUnitExist, err := userEng.IsTableExist(new(unitmodel.CommonUnitData))
	CheckErr(err)
	userUnitExist, err := userEng.IsTableExist(new(usermodel.UserUnitData))
	CheckErr(err)

	if !commonUnitExist || !userUnitExist {
		log.Println("卡片数据不存在，正在同步...")

		userEng.DropTable(new(unitmodel.CommonUnitData))
		userEng.CreateTable(new(unitmodel.CommonUnitData))
		userEng.DropTable(new(usermodel.UserUnitData))
		userEng.CreateTable(new(usermodel.UserUnitData))

		var unitData []unitmodel.UnitM
		err = db.MainEng.Table(new(unitmodel.UnitM)).OrderBy("unit_id ASC").Find(&unitData)
		CheckErr(err)

		checked := false
		for _, u := range unitData {
			// 判断卡片最大等级
			var unitMaxLevel, nextExp, sumExp int
			_, err = db.MainEng.Table("unit_level_up_pattern_m").
				Where("unit_level_up_pattern_id = ?", u.UnitLevelUpPatternId).
				Select("MAX(unit_level),next_exp").Get(&unitMaxLevel, &nextExp)
			CheckErr(err)

			// 计算突破前的经验总和
			_, err = db.MainEng.Table("unit_level_up_pattern_m").
				Where("unit_level_up_pattern_id = ?", u.UnitLevelUpPatternId).
				Where("unit_level = ?", unitMaxLevel-1).Cols("next_exp").Get(&sumExp)
			CheckErr(err)

			// 计算突破前的属性
			var smileMax, pureMax, coolMax int
			smileMax = u.SmileMax
			pureMax = u.PureMax
			coolMax = u.CoolMax

			// 如果 nexpExp 不为零，则说明卡片等级没有达到上限
			if nextExp != 0 {
				// 计算突破后的经验总和
				_, err = db.MainEng.Table("unit_level_limit_pattern_m").
					Where("unit_level_limit_id = 1 AND unit_level = 349").
					Cols("next_exp").Get(&sumExp)
				CheckErr(err)

				// 突破后最大等级
				unitMaxLevel = 350

				// 计算突破后的属性
				smileMax += 6000
				pureMax += 6000
				coolMax += 6000
			}

			// 计算绊值、技能等级、技能经验
			var maxLove, skillLevel, skillExp, removableSkillCapacity, levelLimitID int
			switch u.Rarity {
			case 1:
				maxLove = 50
				skillExp = 0
				skillLevel = 0
				removableSkillCapacity = 0
				levelLimitID = 0
			case 2:
				maxLove = 200
				skillExp = 490
				skillLevel = 8
				removableSkillCapacity = 1
				levelLimitID = 0
			case 3:
				maxLove = 500
				skillExp = 4900
				skillLevel = 8
				removableSkillCapacity = 2
				levelLimitID = 0
			case 4:
				maxLove = 1000
				skillExp = 29900
				skillLevel = 8
				removableSkillCapacity = 8
				levelLimitID = 1
			case 5:
				maxLove = 750
				skillExp = 12700
				skillLevel = 8
				removableSkillCapacity = 3
				levelLimitID = 0
			}

			// 针对技能卡等应援卡片
			if smileMax == 1 {
				maxLove = 0
				skillExp = 0
				skillLevel = 0
				removableSkillCapacity = 0
			}

			// 检查是否签名卡
			var isSigned bool
			exist, err := db.MainEng.Table("unit_sign_asset_m").Where("unit_id = ?", u.UnitId).Exist()
			CheckErr(err)
			if exist {
				isSigned = true
			}

			// 生成公共卡片
			unitCommon := unitmodel.CommonUnitData{
				UnitNumber:                  u.UnitNumber,
				UnitID:                      u.UnitId,
				UnitTypeID:                  u.UnitTypeId,
				Name:                        *u.NameEn,
				Eponym:                      u.EponymEn,
				Rarity:                      u.Rarity,
				Attribute:                   u.AttributeId,
				Smile:                       smileMax,
				Cute:                        pureMax,
				Cool:                        coolMax,
				Exp:                         sumExp,
				Level:                       unitMaxLevel,
				MaxLevel:                    unitMaxLevel,
				LevelLimitID:                levelLimitID,
				Rank:                        u.RankMin,
				MaxRank:                     u.RankMax,
				Love:                        maxLove,
				MaxLove:                     maxLove,
				UnitSkillExp:                skillExp,
				UnitSkillLevel:              skillLevel,
				MaxHp:                       u.HpMax,
				UnitRemovableSkillCapacity:  removableSkillCapacity,
				IsRankMax:                   true,
				IsLoveMax:                   true,
				IsLevelMax:                  true,
				IsSigned:                    isSigned,
				IsSkillLevelMax:             true,
				IsRemovableSkillCapacityMax: true,
				InsertDate:                  time.Now().Unix(),
			}

			_, err = userEng.Insert(&unitCommon)
			CheckErr(err)

			var userID []int
			err = db.UserEng.Table(new(usermodel.Users)).Cols("user_id").Find(&userID)
			CheckErr(err)

			for _, id := range userID {
				userUnit := usermodel.UserUnitData{
					UnitID:       u.UnitId,
					FavoriteFlag: false,
					DisplayRank:  u.RankMax,
					UserID:       id,
					InsertDate:   time.Now().Unix(),
				}

				// 检查表里是否已经有数据
				if !checked {
					ct, err := userEng.Table(new(usermodel.UserUnitData)).Count()
					CheckErr(err)

					if ct == 0 {
						userUnit.UnitOwningUserID = 38383
					}

					checked = true
				}

				_, err = userEng.Insert(&userUnit)
				CheckErr(err)
			}
		}

		err = userEng.Commit()
		CheckErr(err)

		log.Println("同步完成！")
	}
}

func CheckErr(err error) {
	if err != nil {
		userEng.Rollback()
		log.Fatalln("同步失败:", err.Error())
	}
}
