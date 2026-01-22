package user

import (
	"honoka-chan/internal/schema/api/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userGetNavi(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var uID, oID int
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("user_id,unit_owning_user_id").Get(&uID, &oID)
	if ss.CheckErr(err) {
		return
	}

	res = user.GetNaviResp{
		Result: user.GetNaviData{
			User: user.User{
				UserID:           uID,
				UnitOwningUserID: oID,
			},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
