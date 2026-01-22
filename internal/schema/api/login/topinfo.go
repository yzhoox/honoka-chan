package login

type LicenseInfo struct {
	LicenseList  []any `json:"license_list"`
	LicensedInfo []any `json:"licensed_info"`
	ExpiredInfo  []any `json:"expired_info"`
	BadgeFlag    bool  `json:"badge_flag"`
}

type TopInfoData struct {
	FriendActionCnt        int         `json:"friend_action_cnt"`
	FriendGreetCnt         int         `json:"friend_greet_cnt"`
	FriendVarietyCnt       int         `json:"friend_variety_cnt"`
	FriendNewCnt           int         `json:"friend_new_cnt"`
	PresentCnt             int         `json:"present_cnt"`
	SecretBoxBadgeFlag     bool        `json:"secret_box_badge_flag"`
	ServerDatetime         string      `json:"server_datetime"`
	ServerTimestamp        int64       `json:"server_timestamp"`
	NoticeFriendDatetime   string      `json:"notice_friend_datetime"`
	NoticeMailDatetime     string      `json:"notice_mail_datetime"`
	FriendsApprovalWaitCnt int         `json:"friends_approval_wait_cnt"`
	FriendsRequestCnt      int         `json:"friends_request_cnt"`
	IsTodayBirthday        bool        `json:"is_today_birthday"`
	LicenseInfo            LicenseInfo `json:"license_info"`
	UsingBuffInfo          []any       `json:"using_buff_info"`
	IsKlabIDTaskFlag       bool        `json:"is_klab_id_task_flag"`
	KlabIDTaskCanSync      bool        `json:"klab_id_task_can_sync"`
	HasUnreadAnnounce      bool        `json:"has_unread_announce"`
	ExchangeBadgeCnt       []int       `json:"exchange_badge_cnt"`
	AdFlag                 bool        `json:"ad_flag"`
	HasAdReward            bool        `json:"has_ad_reward"`
}

type TopInfoResp struct {
	Result     TopInfoData `json:"result"`
	Status     int         `json:"status"`
	CommandNum bool        `json:"commandNum"`
	TimeStamp  int64       `json:"timeStamp"`
}
