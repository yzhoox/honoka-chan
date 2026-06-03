package session

import (
	"encoding/base64"
	"encoding/json"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type responseHooks struct {
	setHeader func(string, string)
	writeJSON func(int, any)
	writeBody func(int, string)
	abort     func()
}

type Session struct {
	MainEng *xorm.Session
	UserEng *xorm.Session

	UserID   int
	UserPref usermodel.UserPref

	deviceID string
	resp     responseHooks
	done     bool
}

func New(ctx *gin.Context) *Session {
	ss := &Session{
		deviceID: ctx.GetHeader("X-DEVICEID"),
		resp: responseHooks{
			setHeader: ctx.Header,
			writeJSON: ctx.JSON,
			writeBody: func(code int, body string) {
				ctx.String(code, body)
			},
			abort: ctx.Abort,
		},
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

func Attach(ctx *gin.Context) *Session {
	if ss, ok := ctx.Get("session"); ok {
		return ss.(*Session)
	}

	ss := New(ctx)
	ctx.Set("session", ss)
	return ss
}

func (ss *Session) Finalize() {
	if ss.done {
		return
	}
	ss.done = true

	if ss.MainEng != nil {
		ss.MainEng.Close()
		ss.MainEng = nil
	}

	if ss.UserEng == nil {
		return
	}

	if err := ss.UserEng.Commit(); err != nil {
		log.Println(err.Error())
		ss.UserEng.Close()
		ss.UserEng = nil
		ss.resp.writeJSON(500, gin.H{"error": err.Error()})
		ss.resp.abort()
		return
	}

	ss.UserEng.Close()
	ss.UserEng = nil
}

func (ss *Session) Abort(err error) {
	if ss.done {
		return
	}
	ss.done = true

	if ss.MainEng != nil {
		ss.MainEng.Close()
		ss.MainEng = nil
	}
	if ss.UserEng != nil {
		ss.UserEng.Rollback()
		ss.UserEng.Close()
		ss.UserEng = nil
	}

	ss.resp.writeJSON(500, gin.H{"error": err.Error()})
	ss.resp.abort()
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

	ss.resp.setHeader("Content-Type", "application/json")
	ss.resp.setHeader("X-Message-Sign", base64.StdEncoding.EncodeToString(encrypt.RSASignSHA1(data)))
	ss.resp.writeBody(http.StatusOK, string(data))
}
