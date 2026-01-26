package session

import (
	"encoding/json"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type Session struct {
	Ctx *gin.Context

	MainEng *xorm.Session
	UserEng *xorm.Session

	UserID   int
	UserPref usermodel.UserPref
}

func New(ctx *gin.Context) *Session {
	ss := &Session{
		Ctx: ctx,
	}

	ss.MainEng = db.MainEng.NewSession()
	ss.UserEng = db.UserEng.NewSession()
	ss.UserEng.Begin()

	userID := ctx.GetString("userid")
	if userID != "" {
		ss.UserPref = ss.GetUserPref(userID)
		ss.UserID = ss.UserPref.UserID
	}

	return ss
}

func Get(ctx *gin.Context) *Session {
	return ctx.MustGet("session").(*Session)
}

func (ss *Session) Finalize() {
	ss.MainEng.Close()
	if ss.CheckErr(ss.UserEng.Commit()) {
		return
	}
	ss.UserEng.Close()
}

func (ss *Session) Abort(err error) {
	ss.MainEng.Close()
	ss.UserEng.Rollback()
	ss.UserEng.Close()

	ss.Ctx.JSON(500, gin.H{"error": err.Error()})
	ss.Ctx.Abort()
}

func (ss *Session) CheckErr(err error) bool {
	if err != nil {
		log.Println(err.Error())
		ss.Abort(err)
		return true
	}
	return false
}

func (ss *Session) Respond(resp any) {
	data, err := json.Marshal(resp)
	if ss.CheckErr(err) {
		return
	}

	ss.Ctx.Header("X-Message-Sign", utils.GenXMS(data))
	ss.Ctx.String(200, string(data))
}
