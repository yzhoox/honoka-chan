package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	usermodel "honoka-chan/internal/model/user"
	honokautils "honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	"log"
	"net/http"
	"runtime"
	"strings"

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

	if ss.UserEng == nil {
		ss.closeMainSession()
		return
	}

	if err := ss.UserEng.Commit(); err != nil {
		log.Println(err.Error())
		ss.closeAllSessions(false)
		ss.respondError(http.StatusInternalServerError, honokautils.NewInternalErrorContent())
		return
	}

	ss.closeMainSession()
	ss.UserEng.Close()
	ss.UserEng = nil
}

func (ss *Session) Abort(err error) {
	if ss.done {
		return
	}

	msg := ss.formatAbortMessage(err)
	log.Println(msg)
	ss.abortWithContent(http.StatusInternalServerError, honokautils.NewInternalErrorContent())
}

func (ss *Session) AbortWithStatus(status int, content any) {
	if ss.done {
		return
	}
	ss.abortWithContent(status, content)
}

// FinalizeOrRollback commits successful requests and rolls back before a panic reaches recovery.
func (ss *Session) FinalizeOrRollback() {
	defer func() {
		if recovered := recover(); recovered != nil {
			ss.Rollback()
			panic(recovered)
		}
	}()
	ss.Finalize()
}

func (ss *Session) Rollback() {
	if ss.done {
		return
	}
	ss.done = true
	ss.closeAllSessions(true)
}

func (ss *Session) formatAbortMessage(err error) string {
	skip := 2
	if pc, _, _, ok := runtime.Caller(2); ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil && strings.HasSuffix(fn.Name(), ".(*Session).CheckErr") {
			skip = 3
		}
	}
	loc := "unknown"
	if pc, file, line, ok := runtime.Caller(skip); ok {
		fn := runtime.FuncForPC(pc)
		funcName := ""
		if fn != nil {
			funcName = fn.Name()
		}
		loc = fmt.Sprintf("%s:%d (%s)", file, line, funcName)
	}
	return fmt.Sprintf("[%s] %s", loc, err.Error())
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

func (ss *Session) abortWithContent(status int, content any) {
	ss.done = true
	ss.closeAllSessions(true)
	ss.respondError(status, content)
}

func (ss *Session) respondError(status int, content any) {
	honokautils.WriteMaintenanceJSON(ss.resp.setHeader, ss.resp.writeJSON, status, content)
	ss.resp.abort()
}

func (ss *Session) closeAllSessions(rollbackUser bool) {
	ss.closeMainSession()
	if ss.UserEng == nil {
		return
	}
	if rollbackUser {
		_ = ss.UserEng.Rollback()
	}
	ss.UserEng.Close()
	ss.UserEng = nil
}

func (ss *Session) closeMainSession() {
	if ss.MainEng != nil {
		ss.MainEng.Close()
		ss.MainEng = nil
	}
}
