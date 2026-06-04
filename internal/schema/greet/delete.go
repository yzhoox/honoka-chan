package greetschema

type DeleteReq struct {
	Module       string `json:"module"`
	Action       string `json:"action"`
	TimeStamp    int64  `json:"timeStamp"`
	Mgd          int    `json:"mgd"`
	IsSendMail   bool   `json:"is_send_mail"`
	MailNoticeID int    `json:"mail_notice_id"`
	CommandNum   string `json:"commandNum"`
}
