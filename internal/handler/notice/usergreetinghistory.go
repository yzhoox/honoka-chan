package notice

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	noticeschema "honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func userGreetingHistory(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := noticeschema.GreetingNoticeReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	totalCount, err := ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("affector_id = ?", ss.UserID).
		Where("deleted_from_affector = ?", false).
		Count()
	if ss.CheckErr(err) {
		return
	}

	rows := []usermodel.UserGreet{}
	err = ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("affector_id = ?", ss.UserID).
		Where("deleted_from_affector = ?", false).
		OrderBy("notice_id DESC").
		Limit(usermodel.GreetingPageSize).
		Find(&rows)
	if ss.CheckErr(err) {
		return
	}

	noticeList := make([]noticeschema.UserGreetingNotice, 0, len(rows))
	for _, row := range rows {
		receiver, err := getGreetingPeer(ss, row.ReceiverID)
		if ss.CheckErr(err) {
			return
		}

		noticeList = append(noticeList, noticeschema.UserGreetingNotice{
			NoticeID:       row.NoticeID,
			ReferenceTable: 6,
			Message:        row.Message,
			ListMessage:    row.Message,
			InsertDate:     formatGreetingElapsedTime(row.InsertDate),
			Receiver:       receiver,
			ReplyFlag:      row.Reply,
			Readed:         row.Readed,
		})
	}

	ss.Respond(noticeschema.UserGreetingResp{
		ResponseData: noticeschema.UserGreetingData{
			ItemCount:       totalCount,
			HasNext:         totalCount > int64(len(rows)),
			NoticeList:      noticeList,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/notice/noticeUserGreetingHistory", middleware.Common, userGreetingHistory)
}
