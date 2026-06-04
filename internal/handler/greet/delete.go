package greet

import (
	"errors"

	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	greetschema "honoka-chan/internal/schema/greet"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func deleteMail(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := greetschema.DeleteReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	notice := usermodel.UserGreet{}
	has, err := ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("notice_id = ?", req.MailNoticeID).
		Get(&notice)
	if ss.CheckErr(err) {
		return
	}
	if !has {
		ss.CheckErr(errors.New("notice not found"))
		return
	}

	switch {
	case req.IsSendMail && notice.AffectorID == ss.UserID:
		_, err = ss.UserEng.Table(new(usermodel.UserGreet)).
			Where("notice_id = ?", req.MailNoticeID).
			Cols("deleted_from_affector").
			Update(&usermodel.UserGreet{DeletedFromAffector: true})
	case !req.IsSendMail && notice.ReceiverID == ss.UserID:
		_, err = ss.UserEng.Table(new(usermodel.UserGreet)).
			Where("notice_id = ?", req.MailNoticeID).
			Cols("deleted_from_receiver").
			Update(&usermodel.UserGreet{DeletedFromReceiver: true})
	default:
		ss.CheckErr(errors.New("cannot delete this notice"))
		return
	}
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(greetschema.EmptyResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/greet/delete", middleware.Common, deleteMail)
}
