package webui

import (
	"encoding/csv"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/api/profile"
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
		panic(err)
	}

	f, err := os.Open(tmpPath)
	utils.CheckErr(err)
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

		skillLv, err := strconv.Atoi(rr[1])
		if err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "文件解析失败！"})
			return
		}

		var unitId, unitExp, unitRarity, unitHp, unitSigned int
		exists, err := db.MainEng.Table("common_unit_m").Join("LEFT", "unit_m", "common_unit_m.unit_id = unit_m.unit_id").
			Where("unit_m.unit_number = ?", rr[0]).
			Cols("common_unit_m.unit_id,common_unit_m.exp,unit_m.rarity,common_unit_m.max_hp,common_unit_m.is_signed").
			Get(&unitId, &unitExp, &unitRarity, &unitHp, &unitSigned)
		utils.CheckErr(err)

		if !exists {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "卡片不存在！"})
			return
		}

		if unitRarity != 4 {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "仅支持导入UR卡片！"})
			return
		}

		if skillLv < 0 || skillLv > 8 {
			session.Rollback()
			ctx.JSON(http.StatusOK, ErrMsg{Error: 1, Msg: "技能等级设置有误！"})
			return
		}

		var diffExp, diffSmile, diffPure, diffCool int
		_, err = db.MainEng.Table("unit_level_limit_pattern_m").Where("unit_level_limit_id = 1 AND unit_level = 350").
			Cols("next_exp,smile_diff,pure_diff,cool_diff").Get(&diffExp, &diffSmile, &diffPure, &diffCool)
		utils.CheckErr(err)

		isSigned := false
		if unitSigned == 1 {
			isSigned = true
		}

		var skillExp int
		if skillLv != 8 {
			skillExp = 0
		} else {
			skillExp = 29900
		}

		unitData := profile.UnitData{
			UserID:                      ctx.GetInt("userid"),
			UnitID:                      unitId,
			Exp:                         unitExp + diffExp,
			NextExp:                     0,
			Level:                       350,
			MaxLevel:                    350,
			LevelLimitID:                1,
			Rank:                        2,
			MaxRank:                     2,
			Love:                        1000,
			MaxLove:                     1000,
			UnitSkillExp:                skillExp,
			UnitSkillLevel:              skillLv,
			MaxHp:                       unitHp,
			UnitRemovableSkillCapacity:  8,
			FavoriteFlag:                false,
			DisplayRank:                 2,
			IsRankMax:                   true,
			IsLoveMax:                   true,
			IsLevelMax:                  true,
			IsSigned:                    isSigned,
			IsSkillLevelMax:             true,
			IsRemovableSkillCapacityMax: true,
			InsertDate:                  time.Now().Format("2006-01-02 03:04:05"),
		}

		_, err = session.Table("user_unit").Insert(&unitData)
		if err != nil {
			session.Rollback()
			panic(err)
		}

		if err = session.Commit(); err != nil {
			session.Rollback()
			panic(err)
		}
	}

	ctx.JSON(http.StatusOK, ErrMsg{Error: 0, Msg: "上传成功！"})
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
