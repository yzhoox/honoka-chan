package noticeschema

type FriendGreetingData struct {
	NextId          int                    `json:"next_id"`
	NoticeList      []FriendGreetingNotice `json:"notice_list"`
	ServerTimestamp int64                  `json:"server_timestamp"`
}

type FriendGreetingNotice struct {
	NoticeID       int          `json:"notice_id"`
	NewFlag        bool         `json:"new_flag"`
	ReferenceTable int          `json:"reference_table"`
	Message        string       `json:"message"`
	ListMessage    string       `json:"list_message"`
	Readed         bool         `json:"readed"`
	InsertDate     string       `json:"insert_date"`
	Affector       GreetingPeer `json:"affector"`
	ReplyFlag      bool         `json:"reply_flag"`
}

type FriendGreetingResp struct {
	ResponseData FriendGreetingData `json:"response_data"`
	ReleaseInfo  []any              `json:"release_info"`
	StatusCode   int                `json:"status_code"`
}
