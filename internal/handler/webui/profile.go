package webui

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	webuischema "honoka-chan/internal/schema/webui"
	"honoka-chan/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func profilePage(ctx *gin.Context) {
	pref := usermodel.UserPref{}
	exists, err := db.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", ctx.GetInt("userid")).Get(&pref)
	if err != nil || !exists {
		ctx.String(http.StatusInternalServerError, "用户资料读取失败")
		return
	}

	ctx.HTML(http.StatusOK, "admin/profile.html", gin.H{
		"menu": 2,
		"url":  strings.Split(ctx.Request.URL.String(), "?")[0],
		"pref": pref,
	})
}

func updateProfile(ctx *gin.Context) {
	parseInt := func(key string, min int) (int, error) {
		value := strings.TrimSpace(ctx.PostForm(key))
		number, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		if number < min {
			return 0, strconv.ErrSyntax
		}
		return number, nil
	}

	userName := strings.TrimSpace(ctx.PostForm("user_name"))
	inviteCode := strings.TrimSpace(ctx.PostForm("invite_code"))
	userDesc := strings.TrimSpace(ctx.PostForm("user_desc"))
	if userName == "" || inviteCode == "" {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "用户名和邀请码不能为空！",
		})
		return
	}

	userLevel, err := parseInt("user_level", 1)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "等级格式有误！"})
		return
	}
	userExp, err := parseInt("user_exp", 0)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "当前经验格式有误！"})
		return
	}
	nextExp, err := parseInt("next_exp", 0)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "升级经验格式有误！"})
		return
	}
	gameCoin, err := parseInt("game_coin", 0)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "金币格式有误！"})
		return
	}
	snsCoin, err := parseInt("sns_coin", 0)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "爱心格式有误！"})
		return
	}
	energyMax, err := parseInt("energy_max", 1)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "体力上限格式有误！"})
		return
	}
	overMaxEnergy, err := parseInt("over_max_energy", 0)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{Code: 1, Message: "当前体力格式有误！"})
		return
	}

	pref := usermodel.UserPref{
		UserName:       userName,
		UserLevel:      userLevel,
		UserDesc:       userDesc,
		InviteCode:     inviteCode,
		UserExp:        userExp,
		NextExp:        nextExp,
		GameCoin:       gameCoin,
		SnsCoin:        snsCoin,
		EnergyMax:      energyMax,
		OverMaxEnergy:  overMaxEnergy,
		ProfileVersion: usermodel.CurrentUserPrefProfileVersion,
		UpdateTime:     time.Now().Unix(),
	}
	_, err = db.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", ctx.GetInt("userid")).
		Cols(
			"user_name",
			"user_level",
			"user_desc",
			"invite_code",
			"user_exp",
			"next_exp",
			"game_coin",
			"sns_coin",
			"energy_max",
			"over_max_energy",
			"profile_version",
			"update_time",
		).
		Update(&pref)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "保存失败！",
		})
		return
	}

	ctx.JSON(http.StatusOK, webuischema.Msg{
		Code:    0,
		Message: "保存成功！",
	})
}

func resetProfile(ctx *gin.Context) {
	pref := usermodel.UserPref{}
	exists, err := db.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", ctx.GetInt("userid")).Get(&pref)
	if err != nil || !exists {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "用户资料读取失败！",
		})
		return
	}

	pref.ResetProfileDefaults()
	pref.UpdateTime = time.Now().Unix()

	_, err = db.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", ctx.GetInt("userid")).
		Cols(
			"user_name",
			"user_level",
			"user_desc",
			"invite_code",
			"user_exp",
			"next_exp",
			"game_coin",
			"sns_coin",
			"energy_max",
			"over_max_energy",
			"profile_version",
			"update_time",
		).
		Update(&pref)
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "重置失败！",
		})
		return
	}

	ctx.JSON(http.StatusOK, webuischema.Msg{
		Code:     0,
		Message:  "已恢复默认值！",
		Redirect: "/admin/profile",
	})
}

func init() {
	router.AddHandler("admin", "GET", "/profile", middleware.WebAuth, profilePage)
	router.AddHandler("admin", "POST", "/profile", middleware.WebAuth, updateProfile)
	router.AddHandler("admin", "POST", "/profile/reset", middleware.WebAuth, resetProfile)
}
