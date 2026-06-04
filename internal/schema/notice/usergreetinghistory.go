package noticeschema

type UserGreetingData struct {
	ItemCount       int                  `json:"item_count"`
	HasNext         bool                 `json:"has_next"`
	NoticeList      []UserGreetingNotice `json:"notice_list"`
	ServerTimestamp int64                `json:"server_timestamp"`
}

type UserGreetingNotice struct {
	NoticeID       int          `json:"notice_id"`
	ReferenceTable int          `json:"reference_table"`
	Message        string       `json:"message"`
	ListMessage    string       `json:"list_message"`
	InsertDate     string       `json:"insert_date"`
	Receiver       GreetingPeer `json:"receiver"`
	ReplyFlag      bool         `json:"reply_flag"`
	Readed         bool         `json:"readed"`
}

type UserGreetingResp struct {
	ResponseData UserGreetingData `json:"response_data"`
	ReleaseInfo  []any            `json:"release_info"`
	StatusCode   int              `json:"status_code"`
}
