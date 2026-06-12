package achievementapischema

type AchievementRewardItem struct {
	AddType int `json:"add_type"`
	ItemID  int `json:"item_id"`
	Amount  int `json:"amount"`
}

type AchievementListItem struct {
	AchievementID  int                     `json:"achievement_id"`
	Count          int                     `json:"count"`
	IsAccomplished bool                    `json:"is_accomplished"`
	InsertDate     string                  `json:"insert_date"`
	EndDate        *string                 `json:"end_date"`
	RemainingTime  string                  `json:"remaining_time"`
	IsNew          bool                    `json:"is_new"`
	ForDisplay     bool                    `json:"for_display"`
	RewardList     []AchievementRewardItem `json:"reward_list"`
}

type UnaccomplishListData struct {
	FilterCategoryID int                   `json:"filter_category_id"`
	AchievementList  []AchievementListItem `json:"achievement_list"`
	Count            int                   `json:"count"`
	IsLast           bool                  `json:"is_last"`
}

type UnaccomplishListResp struct {
	Result     []UnaccomplishListData `json:"result"`
	Status     int                    `json:"status"`
	CommandNum bool                   `json:"commandNum"`
	TimeStamp  int64                  `json:"timeStamp"`
}
