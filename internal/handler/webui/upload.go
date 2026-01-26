package webui

import (
	"encoding/csv"
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"net/http"
	"os"
	"path"
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
	file, err := ctx.FormFile("file")
	utils.CheckErr(err)

	tmpPath := path.Join("./temp", file.Filename)
	err = ctx.SaveUploadedFile(file, tmpPath)
	utils.CheckErr(err)

	session := db.UserEng.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: err.Error()})
		return
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文件创建失败！"})
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	rs, err := r.ReadAll()
	if err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文件解析失败！"})
		return
	}
	for _, rr := range rs {
		if len(rr) != 2 || rr[0] == "" || rr[1] == "" {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文件解析失败！"})
			return
		}

		skLv, err := strconv.Atoi(rr[1])
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文件解析失败！"})
			return
		}

		uData := unitmodel.CommonUnitData{}
		exists, err := session.Table(new(unitmodel.CommonUnitData)).Where("unit_number = ?", rr[0]).Get(&uData)
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
			UpdateDate:   time.Now().Unix(),
		}

		_, err = session.Table(new(usermodel.UserUnitData)).Insert(&unitData)
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片添加失败！"})
			return
		}

		if err = session.Commit(); err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片添加失败！"})
			return
		}
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
