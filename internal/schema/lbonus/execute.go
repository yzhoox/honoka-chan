package lbonusschema

type Item struct {
	ItemID  int `json:"item_id"`
	AddType int `json:"add_type"`
	Amount  int `json:"amount"`
}

type Day struct {
	Day               int    `json:"day"`
	DayOfTheWeek      int    `json:"day_of_the_week"`
	SpecialDay        bool   `json:"special_day"`
	SpecialImageAsset string `json:"special_image_asset"`
	Received          bool   `json:"received"`
	AdReceived        bool   `json:"ad_received"`
	Item              Item   `json:"item"`
}

type Month struct {
	Year  int   `json:"year"`
	Month int   `json:"month"`
	Days  []Day `json:"days"`
}

type CalendarInfo struct {
	CurrentDate  string `json:"current_date"`
	CurrentMonth Month  `json:"current_month"`
	NextMonth    Month  `json:"next_month"`
}

type Reward struct {
	ItemID  int `json:"item_id"`
	AddType int `json:"add_type"`
	Amount  int `json:"amount"`
}

type TotalLoginInfo struct {
	LoginCount     int      `json:"login_count"`
	RemainingCount int      `json:"remaining_count"`
	Reward         []Reward `json:"reward"`
}

type RankInfo struct {
	BeforeClassRankID int    `json:"before_class_rank_id"`
	AfterClassRankID  int    `json:"after_class_rank_id"`
	RankUpDate        string `json:"rank_up_date"`
}

type ClassSystem struct {
	RankInfo     RankInfo `json:"rank_info"`
	CompleteFlag bool     `json:"complete_flag"`
	IsOpened     bool     `json:"is_opened"`
	IsVisible    bool     `json:"is_visible"`
}

type Rewards struct {
	Rarity         int    `json:"rarity"`
	ItemID         int    `json:"item_id"`
	AddType        int    `json:"add_type"`
	Amount         int    `json:"amount"`
	ItemCategoryID int    `json:"item_category_id"`
	RewardBoxFlag  bool   `json:"reward_box_flag"`
	InsertDate     string `json:"insert_date"`
}

type EffortPoint struct {
	LiveEffortPointBoxSpecID int       `json:"live_effort_point_box_spec_id"`
	Capacity                 int       `json:"capacity"`
	Before                   int       `json:"before"`
	After                    int       `json:"after"`
	Rewards                  []Rewards `json:"rewards"`
}

type Parameter struct {
	Smile int `json:"smile"`
	Pure  int `json:"pure"`
	Cool  int `json:"cool"`
}

type Museum struct {
	Parameter      Parameter `json:"parameter"`
	ContentsIDList []int     `json:"contents_id_list"`
}

type ExecuteData struct {
	Sheets            []any          `json:"sheets"`
	CalendarInfo      CalendarInfo   `json:"calendar_info"`
	TotalLoginInfo    TotalLoginInfo `json:"total_login_info"`
	LicenseLbonusList []any          `json:"license_lbonus_list"`
	ClassSystem       ClassSystem    `json:"class_system"`
	StartDashSheets   []any          `json:"start_dash_sheets"`
	EffortPoint       []EffortPoint  `json:"effort_point"`
	LimitedEffortBox  []any          `json:"limited_effort_box"`
	MuseumInfo        Museum         `json:"museum_info"`
	ServerTimestamp   int64          `json:"server_timestamp"`
	PresentCnt        int            `json:"present_cnt"`
}

type ExecuteResp struct {
	ResponseData ExecuteData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
