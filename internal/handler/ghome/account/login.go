package account

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitapischema "honoka-chan/internal/schema/api/unit"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/utils"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
	"xorm.io/xorm"
)

const (
	defaultUserID       = 377385143
	defaultAwardID      = 113 // 极推穗乃果
	defaultBackgroundID = 65  // 穗乃果的房间
	defaultUnitID       = 338 // 高坂穂乃果 直到那天来临
	defaultUserName     = "梦路"
	defaultUserDesc     = "你好。"
)

func login(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.Finalize()

	data, err := ctx.GetRawData()
	if ss.CheckErr(err) {
		return
	}

	data, err = base64.StdEncoding.DecodeString(string(data))
	if ss.CheckErr(err) {
		return
	}

	randKey, err := ss.Get3DESRandKey()
	if ss.CheckErr(err) {
		return
	}

	decryptedData, err := openssl.Des3ECBDecrypt(data, randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	queryStr, err := url.QueryUnescape(string(decryptedData))
	if ss.CheckErr(err) {
		return
	}
	params, err := url.ParseQuery(queryStr)
	if ss.CheckErr(err) {
		return
	}
	phone, password := normalizePhone(params.Get("phone")), params.Get("password")
	if phone == "" || password == "" {
		ss.Abort(errors.New("invalid login params"))
		return
	}

	loginData, loginCode, loginMsg, _, err := AddUserWithSession(ss.UserEng, phone, password, false)
	if ss.CheckErr(err) {
		return
	}

	data, err = json.Marshal(loginData)
	if ss.CheckErr(err) {
		return
	}
	encryptedData, err := openssl.Des3ECBEncrypt([]byte(data), randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.LoginResp{
		Code: loginCode,
		Msg:  loginMsg,
		Data: base64.StdEncoding.EncodeToString(encryptedData),
	})
}

func AddUser(phone, password string, isDefault bool) (ghomeschema.LoginData, int, string, bool, error) {
	return addUser(nil, phone, password, isDefault)
}

func AddUserWithSession(dbSession *xorm.Session, phone, password string, isDefault bool) (ghomeschema.LoginData, int, string, bool, error) {
	return addUser(dbSession, phone, password, isDefault)
}

func addUser(dbSession *xorm.Session, phone, password string, isDefault bool) (ghomeschema.LoginData, int, string, bool, error) {
	loginData := ghomeschema.LoginData{}
	loginCode := 0
	loginMsg := "ok"
	loginTime := time.Now().Unix()
	created := false
	phone = normalizePhone(phone)
	if phone == "" {
		return loginData, loginCode, loginMsg, created, errors.New("invalid phone")
	}

	var userID int
	var pass, autoKey, ticket string

	localSession := dbSession == nil
	if localSession {
		dbSession = db.UserEng.NewSession()
		defer dbSession.Close()
		if err := dbSession.Begin(); err != nil {
			return loginData, loginCode, loginMsg, created, err
		}
	}

	_, err := dbSession.Table("users").Cols("password,autokey,ticket,user_id").
		Where("phone = ?", phone).Get(&pass, &autoKey, &ticket, &userID)
	if err != nil {
		if localSession {
			dbSession.Rollback()
		}
		return loginData, loginCode, loginMsg, created, err
	}

	if pass == "" {
		// 未注册 - 自动注册
		awardID := 1      // 音乃木坂学生
		backgroundID := 1 // 初始背景
		unitID := 31      // 初始高坂穂乃果
		userName := usermodel.DefaultAutoUserName
		userDesc := usermodel.DefaultAutoUserDesc
		// 是否初始化
		if isDefault {
			// 初始化默认用户
			userID = defaultUserID
			awardID = defaultAwardID
			backgroundID = defaultBackgroundID
			unitID = defaultUnitID
			userName = defaultUserName
			userDesc = defaultUserDesc
		} else {
			// 检查是否 userID 已经注册
			userID, err = getAvailableUserID(int(loginTime))
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}
		}

		pass = openssl.Md5ToString(password)
		autoKey = "AUTO" + strings.ToUpper(utils.RandomStr(32))
		ticket = fmt.Sprintf("9999999%d%d", userID, userID)

		userData := usermodel.Users{
			Phone:         phone,
			Password:      pass,
			Autokey:       autoKey,
			Ticket:        ticket,
			UserID:        userID,
			LastLoginTime: loginTime,
		}
		_, err = dbSession.Table("users").Insert(&userData)
		if err != nil {
			if localSession {
				dbSession.Rollback()
			}
			return loginData, loginCode, loginMsg, created, err
		}
		created = true

		// 方便起见初始化 userid 和 key 一样
		// 注意：user_key 表中的 key 是上文生成的用于登录的 userid，而 userid 则是用于 Authorize Token 生成用的
		userKey := usermodel.UserKey{
			UserID: userID,
			Key:    userID,
		}
		_, err = dbSession.Table("user_key").Insert(&userKey)
		if err != nil {
			if localSession {
				dbSession.Rollback()
			}
			return loginData, loginCode, loginMsg, created, err
		}

		// 检查用户配置
		exists, err := dbSession.Table("user_pref").Where("user_id = ?", userID).Exist()
		if err != nil {
			if localSession {
				dbSession.Rollback()
			}
			return loginData, loginCode, loginMsg, created, err
		}

		if !exists {
			// 生成用于卡片
			var unitData []unitmodel.CommonUnitData
			err = dbSession.Table(new(unitmodel.CommonUnitData)).Find(&unitData)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}

			checked := false
			for _, u := range unitData {
				userUnit := usermodel.UserUnitData{
					UnitID:       u.UnitID,
					FavoriteFlag: false,
					DisplayRank:  u.MaxRank,
					UserID:       userID,
					InsertDate:   time.Now().Unix(),
				}

				// 检查表里是否已经有数据
				if !checked {
					ct, err := dbSession.Table(new(usermodel.UserUnitData)).Count()
					if err != nil {
						if localSession {
							dbSession.Rollback()
						}
						return loginData, loginCode, loginMsg, created, err
					}

					if ct == 0 {
						userUnit.UnitOwningUserID = 38383
					}
				}

				_, err = dbSession.Table(new(usermodel.UserUnitData)).Insert(&userUnit)
				if err != nil {
					if localSession {
						dbSession.Rollback()
					}
					return loginData, loginCode, loginMsg, created, err
				}

				checked = true
			}

			// 默认中心成员
			var unitOwningUserID int
			_, err = dbSession.Table(new(usermodel.UserUnitData)).
				Cols("unit_owning_user_id").
				Where("user_id = ?", userID).
				Where("unit_id = ?", unitID).
				Get(&unitOwningUserID)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}

			userPref := usermodel.UserPref{
				UserID:           userID,
				AwardID:          awardID,
				BackgroundID:     backgroundID,
				UnitOwningUserID: unitOwningUserID,
				UserName:         userName,
				UserLevel:        usermodel.DefaultUserLevel,
				UserDesc:         userDesc,
				InviteCode:       strconv.Itoa(userID),
				UserExp:          usermodel.DefaultUserExp,
				NextExp:          usermodel.DefaultUserNextExp,
				GameCoin:         usermodel.DefaultUserGameCoin,
				SnsCoin:          usermodel.DefaultUserSnsCoin,
				EnergyMax:        usermodel.DefaultUserEnergyMax,
				OverMaxEnergy:    usermodel.DefaultUserOverMaxEnergy,
				BirthMonth:       usermodel.DefaultBirthMonth,
				BirthDay:         usermodel.DefaultBirthDay,
				ProfileVersion:   usermodel.CurrentUserPrefProfileVersion,
				UpdateTime:       time.Now().Unix(),
			}
			_, err = dbSession.Table("user_pref").Insert(&userPref)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}
		}

		// 检查用户卡组配置
		exists, err = dbSession.Table("user_deck").Where("user_id = ?", userID).Exist()
		if err != nil {
			if localSession {
				dbSession.Rollback()
			}
			return loginData, loginCode, loginMsg, created, err
		}

		if !exists {
			// 默认队伍
			userDeck := unitapischema.UserDeckData{
				DeckID:     1,
				MainFlag:   1,
				DeckName:   "队伍A",
				UserID:     userID,
				InsertDate: time.Now().Unix(),
			}
			_, err = dbSession.Table("user_deck").Insert(&userDeck)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}
			userDeckID := userDeck.ID

			// 默认卡组 - 仆光
			unitID := []int{3465, 3466, 3467, 3468, 3469, 3470, 3471, 3472, 3473}
			var unitData []unitmodel.UnitDataMap
			err = dbSession.Table("user_unit_data").Alias("a").
				Join("LEFT", "common_unit_data", "a.unit_id = common_unit_data.unit_id").
				Cols(`a.unit_owning_user_id,a.favorite_flag,a.display_rank,common_unit_data.*`).
				Where("a.user_id = ?", userID).
				In("a.unit_id", unitID).Find(&unitData)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}

			for i, u := range unitData {
				unitDeckData := usermodel.UserDeckUnit{
					UserDeckID:       userDeckID,
					UnitOwningUserID: u.UnitOwningUserID,
					UnitID:           u.UnitID,
					Position:         i + 1,
					Level:            u.Level,
					LevelLimitID:     u.LevelLimitID,
					DisplayRank:      u.DisplayRank,
					Love:             u.MaxLove,
					UnitSkillLevel:   u.UnitSkillLevel,
					IsRankMax:        u.IsRankMax,
					IsLoveMax:        u.IsLoveMax,
					IsLevelMax:       u.IsLevelMax,
					IsSigned:         u.IsSigned,
					BeforeLove:       u.MaxLove,
					MaxLove:          u.MaxLove,
					UserID:           userID,
					InsertDate:       time.Now().Unix(),
				}
				_, err = dbSession.Table("user_deck_unit").Insert(&unitDeckData)
				if err != nil {
					if localSession {
						dbSession.Rollback()
					}
					return loginData, loginCode, loginMsg, created, err
				}
			}
		}

		if !isDefault {
			err = usermodel.EnsureDefaultFriendship(dbSession, userID)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}
		}

		loginData.Autokey = autoKey
		loginData.HasRealInfo = 1
		loginData.Message = "ok"
		loginData.RealInfoForce = 1
		loginData.Ticket = ticket
		loginData.UserAttribute = "0"
		loginData.UserID = userID
	} else {
		// 已注册 - 检查密码
		if pass != openssl.Md5ToString(password) {
			loginCode = 31
			loginMsg = "账号不存在或者密码有误！"
		} else {
			userData := usermodel.Users{
				Autokey:       autoKey,
				Ticket:        ticket,
				LastLoginTime: loginTime,
			}
			_, err = dbSession.Table("users").Where("user_id = ?", userID).Update(&userData)
			if err != nil {
				if localSession {
					dbSession.Rollback()
				}
				return loginData, loginCode, loginMsg, created, err
			}

			loginData.Autokey = autoKey // 注意：更换设备（deviceId 发生变化）应重新生成 autokey
			loginData.HasRealInfo = 1
			loginData.Message = "ok"
			loginData.RealInfoForce = 1
			loginData.Ticket = fmt.Sprintf("9999999%d%d", userID, loginTime) // 实际登录用的密码（每次登录都会重新生成新的）
			loginData.UserAttribute = "0"
			loginData.UserID = userID // 实际登录用的账号
		}
	}

	if localSession {
		if err := dbSession.Commit(); err != nil {
			dbSession.Rollback()
			return loginData, loginCode, loginMsg, created, err
		}
	}

	return loginData, loginCode, loginMsg, created, nil
}

func getAvailableUserID(seed int) (int, error) {
	const maxAttempts = 16

	userID := seed
	for attempt := range maxAttempts {
		exist, err := db.UserEng.Table("users").Where("user_id = ?", userID).Exist()
		if err != nil {
			return 0, err
		}
		if !exist {
			return userID, nil
		}

		// Move to millisecond precision to reduce repeated collisions.
		userID = int(time.Now().UnixMilli()) + attempt + 1
	}

	return 0, errors.New("failed to allocate user id")
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}

	// 客户端上传的手机号可能带国家/地区区号，如 "86-13800000000"。
	if idx := strings.Index(phone, "-"); idx >= 0 && idx+1 < len(phone) {
		return strings.TrimSpace(phone[idx+1:])
	}
	return phone
}

func init() {
	router.AddHandler("v1", "POST", "/account/login", login)
}
