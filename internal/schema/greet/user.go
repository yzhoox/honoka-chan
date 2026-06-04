package greetschema

type UserReq struct {
	Module          string `json:"module"`
	Mgd             int    `json:"mgd"`
	Action          string `json:"action"`
	TimeStamp       int64  `json:"timeStamp"`
	ToUserID        int    `json:"to_user_id"`
	Message         string `json:"message"`
	RepliedNoticeID int    `json:"replied_notice_id"`
	CommandNum      string `json:"commandNum"`
}
