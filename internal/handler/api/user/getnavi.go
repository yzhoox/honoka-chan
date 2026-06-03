package user

import (
	userapischema "honoka-chan/internal/schema/api/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userGetNavi(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var uID, oID int
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("user_id,unit_owning_user_id").Get(&uID, &oID)
	if err != nil {
		return nil, err
	}

	res = userapischema.GetNaviResp{
		Result: userapischema.GetNaviData{
			User: userapischema.User{
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
