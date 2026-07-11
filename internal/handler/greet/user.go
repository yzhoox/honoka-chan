package greet

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	greetschema "honoka-chan/internal/schema/greet"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func user(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := greetschema.UserReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		ss.CheckErr(errors.New("message is required"))
		return
	}
	if len([]rune(message)) > 200 {
		ss.CheckErr(errors.New("message too long"))
		return
	}

	targetUserID, err := resolveGreetingUserID(ss, req.ToUserID)
	if ss.CheckErr(err) {
		return
	}
	if targetUserID == ss.UserID {
		ss.CheckErr(errors.New("cannot greet self"))
		return
	}

	replyFlag := false
	if req.RepliedNoticeID > 0 {
		exists, err := ss.UserEng.Table(new(usermodel.UserGreet)).
			Where("notice_id = ?", req.RepliedNoticeID).
			Exist()
		if ss.CheckErr(err) {
			return
		}
		if !exists {
			ss.CheckErr(errors.New("replied notice not found"))
			return
		}
		replyFlag = true
	}

	_, err = ss.UserEng.Table(new(usermodel.UserGreet)).Insert(&usermodel.UserGreet{
		AffectorID:          ss.UserID,
		ReceiverID:          targetUserID,
		Message:             message,
		Reply:               replyFlag,
		Readed:              false,
		DeletedFromAffector: false,
		DeletedFromReceiver: false,
		InsertDate:          time.Now().Unix(),
	})
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(greetschema.EmptyResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/greet/user", middleware.Common, user)
}
