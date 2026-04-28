package account

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitapischema "honoka-chan/internal/schema/api/unit"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/utils"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func login(ctx *gin.Context) {
	ss := session.New(ctx)
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

	queryStr, _ := url.QueryUnescape(string(decryptedData))
	params, _ := url.ParseQuery(queryStr)
	phone, password := params.Get("phone"), params.Get("password")

	var userID int
	var pass, autoKey, ticket string
	_, err = ss.UserEng.Table("users").Cols("password,autokey,ticket,user_id").
		Where("phone = ?", phone).Get(&pass, &autoKey, &ticket, &userID)
	if ss.CheckErr(err) {
		return
	}

	loginData := ghomeschema.LoginData{}
	loginCode := 0
	loginMsg := "ok"
	loginTime := time.Now().Unix()
	if pass == "" {
		// 未注册 - 自动注册
		// 检查是否 userID 已经注册
		userID, _ = checkUserID(ss, int(loginTime))

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
		_, err = ss.UserEng.Table("users").Insert(&userData)
		if ss.CheckErr(err) {
			return
		}

		// 方便起见初始化 userid 和 key 一样
		// 注意：user_key 表中的 key 是上文生成的用于登录的 userid，而 userid 则是用于 Authorize Token 生成用的
		userKey := usermodel.UserKey{
			UserID: userID,
			Key:    userID,
		}
		_, err = ss.UserEng.Table("user_key").Insert(&userKey)
		if ss.CheckErr(err) {
			return
		}

		// 检查用户配置
		exists, err := ss.UserEng.Table("user_pref").Where("user_id = ?", userID).Exist()
		if ss.CheckErr(err) {
			return
		}

		if !exists {
			// 生成用于卡片
			var unitData []unitmodel.CommonUnitData
			err = ss.UserEng.Table(new(unitmodel.CommonUnitData)).Find(&unitData)
			if ss.CheckErr(err) {
				return
			}

			checked := false
			for _, u := range unitData {
				userUnit := usermodel.UserUnitData{
					UnitID:       u.UnitID,
					FavoriteFlag: false,
					DisplayRank:  u.MaxRank,
					UserID:       userID,
					InsertDate:   time.Now().Unix(),
					UpdateDate:   time.Now().Unix(),
				}

				// 检查表里是否已经有数据
				if !checked {
					ct, err := ss.UserEng.Table(new(usermodel.UserUnitData)).Count()
					if ss.CheckErr(err) {
						return
					}

					if ct == 0 {
						userUnit.UnitOwningUserID = 38383
					}

					checked = true
				}

				_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).Insert(&userUnit)
				if ss.CheckErr(err) {
					return
				}
			}

			// 默认中心成员
			var unitOwningUserID int
			_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).
				Cols("unit_owning_user_id").Where("unit_id = ?", 31).Get(&unitOwningUserID)
			if ss.CheckErr(err) {
				return
			}

			userPref := usermodel.UserPref{
				UserID:           userID,
				AwardID:          1, // 音乃木坂学生
				BackgroundID:     1, // 初始背景
				UnitOwningUserID: unitOwningUserID,
				UserName:         "音乃木坂学生",
				UserLevel:        1028,
				UserDesc:         "你好。",
				InviteCode:       strconv.Itoa(userID),
				UpdateTime:       time.Now().Unix(),
			}
			_, err = ss.UserEng.Table("user_pref").Insert(&userPref)
			if ss.CheckErr(err) {
				return
			}
		}

		// 检查用户卡组配置
		exists, err = ss.UserEng.Table("user_deck").Where("user_id = ?", userID).Exist()
		if ss.CheckErr(err) {
			return
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
			_, err = ss.UserEng.Table("user_deck").Insert(&userDeck)
			userDeckID := userDeck.ID

			// 默认卡组 - 仆光
			unitID := []int{3465, 3466, 3467, 3468, 3469, 3470, 3471, 3472, 3473}
			var unitData []unitmodel.UnitDataMap
			err = ss.GetBasicUnitInfo().In("a.unit_id", unitID).Find(&unitData)
			if ss.CheckErr(err) {
				return
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
				_, err = ss.UserEng.Table("user_deck_unit").Insert(&unitDeckData)
				if ss.CheckErr(err) {
					return
				}
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
			_, err = ss.UserEng.Table("users").Where("user_id = ?", userID).Update(&userData)
			if ss.CheckErr(err) {
				return
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

func checkUserID(ss *session.Session, userID int) (int, bool) {
	exist, err := ss.UserEng.Table("users").Where("user_id = ?", userID).Exist()
	if ss.CheckErr(err) {
		return 0, false
	}

	if exist {
		userID = int(time.Now().Unix())
		return checkUserID(ss, userID)
	}

	return userID, true
}

func init() {
	router.AddHandler("v1", "POST", "/account/login", login)
}
