package notice

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	noticeschema "honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func friendGreeting(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := noticeschema.GreetingNoticeReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	query := ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("receiver_id = ?", ss.UserID).
		Where("deleted_from_receiver = ?", false)
	if req.NextID > 0 {
		query = query.Where("notice_id < ?", req.NextID)
	}

	rows := []usermodel.UserGreet{}
	err = query.
		OrderBy("notice_id DESC").
		Limit(usermodel.GreetingPageSize).
		Find(&rows)
	if ss.CheckErr(err) {
		return
	}

	noticeList := make([]noticeschema.FriendGreetingNotice, 0, len(rows))
	nextID := 0
	for _, row := range rows {
		affector, err := getGreetingPeer(ss, row.AffectorID)
		if ss.CheckErr(err) {
			return
		}

		noticeList = append(noticeList, noticeschema.FriendGreetingNotice{
			NoticeID:       row.NoticeID,
			NewFlag:        true,
			ReferenceTable: 6,
			Message:        row.Message,
			ListMessage:    row.Message,
			Readed:         true,
			InsertDate:     formatGreetingElapsedTime(row.InsertDate),
			Affector:       affector,
			ReplyFlag:      row.Reply,
		})
		nextID = row.NoticeID
	}

	if len(rows) < usermodel.GreetingPageSize {
		nextID = 0
	}

	_, err = ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("receiver_id = ?", ss.UserID).
		Where("readed = ?", false).
		Cols("readed").
		Update(&usermodel.UserGreet{Readed: true})
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(noticeschema.FriendGreetingResp{
		ResponseData: noticeschema.FriendGreetingData{
			NextId:          nextID,
			NoticeList:      noticeList,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/notice/noticeFriendGreeting", middleware.Common, friendGreeting)
}
