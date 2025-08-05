package session

import (
	"encoding/json"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type Session struct {
	Ctx     *gin.Context
	MainEng *xorm.Session
	UserEng *xorm.Session
}

func New(ctx *gin.Context) *Session {
	ss := &Session{
		Ctx: ctx,
	}

	ss.MainEng = db.MainEng.NewSession()

	ss.UserEng = db.UserEng.NewSession()
	ss.UserEng.Begin()

	return ss
}

func (ss *Session) Finalize() {
	ss.MainEng.Close()

	ss.UserEng.Commit()
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
		ss.Abort(err)
		return true
	}
	return false
}

func (ss *Session) Respond(resp any) {
	data, err := json.Marshal(resp)
	if err != nil {
		ss.Abort(err)
		return
	}

	ss.Ctx.Header("X-Message-Sign", utils.GenXMS(data))
	ss.Ctx.String(200, string(data))
}
