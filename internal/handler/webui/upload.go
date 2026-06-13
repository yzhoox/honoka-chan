package webui

import (
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	"honoka-chan/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrMsg struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
}

func upload(ctx *gin.Context) {
	content := strings.TrimSpace(ctx.PostForm("content"))
	if content == "" {
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "请输入导入内容！"})
		return
	}

	session := db.UserEng.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: err.Error()})
		return
	}

	lines := strings.SplitSeq(content, "\n")
	for rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		line = strings.ReplaceAll(line, "，", ",")

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文本解析失败，请检查每行格式是否为 卡片ID,技能等级"})
			return
		}

		unitNumber := strings.TrimSpace(parts[0])
		skillLevelRaw := strings.TrimSpace(parts[1])
		if unitNumber == "" || skillLevelRaw == "" {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文本解析失败，请检查每行格式是否为 卡片ID,技能等级"})
			return
		}

		skLv, err := strconv.Atoi(skillLevelRaw)
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文本解析失败，请检查每行格式是否为 卡片ID,技能等级"})
			return
		}

		uData := unitmodel.CommonUnitData{}
		exists, err := session.Table(new(unitmodel.CommonUnitData)).Where("unit_number = ?", unitNumber).Get(&uData)
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片查询失败！"})
			return
		}
		if !exists {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片不存在！"})
			return
		}

		if uData.Rarity != 4 {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "仅支持导入UR卡片！"})
			return
		}

		if skLv < 0 || skLv > 8 {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "技能等级设置有误！"})
			return
		}

		unitData := usermodel.UserUnitData{
			UnitID:       uData.UnitID,
			FavoriteFlag: false,
			DisplayRank:  uData.MaxRank,
			UserID:       ctx.GetInt("userid"),
			InsertDate:   time.Now().Unix(),
		}

		_, err = session.Table(new(usermodel.UserUnitData)).Insert(&unitData)
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片添加失败！"})
			return
		}
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片添加失败！"})
		return
	}

	if err := usermodel.MarkUserForceRelogin(db.UserEng, ctx.GetInt("userid")); err != nil {
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片添加成功，但登录状态刷新失败！"})
		return
	}

	ctx.JSON(http.StatusOK, ErrMsg{Error: 0, Msg: "导入成功，请重新打开游戏！"})
}

func init() {
	router.AddHandler("admin", "GET", "/upload", middleware.WebAuth, func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin/upload.html", gin.H{
			"menu": 1,
			"url":  strings.Split(ctx.Request.URL.String(), "?")[0],
		})
	})
	router.AddHandler("admin", "POST", "/upload", middleware.WebAuth, upload)
}
