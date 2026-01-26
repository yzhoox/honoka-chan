package liveapischema

type LiveList struct {
	LiveDifficultyID int    `json:"live_difficulty_id"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	IsRandom         bool   `json:"is_random"`
}

type LimitedBonusCommonList struct {
	LiveType          int    `json:"live_type"`
	LimitedBonusType  int    `json:"limited_bonus_type"`
	LimitedBonusValue int    `json:"limited_bonus_value"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
}

type RandomLiveList struct {
	AttributeID int    `json:"attribute_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type TrainingLiveList struct {
	LiveDifficultyID int    `json:"live_difficulty_id"`
	StartDate        string `json:"start_date"`
	IsRandom         bool   `json:"is_random"`
}

type ScheduleData struct {
	EventList              []any                    `json:"event_list"`
	LiveList               []LiveList               `json:"live_list"`
	LimitedBonusList       []any                    `json:"limited_bonus_list"`
	LimitedBonusCommonList []LimitedBonusCommonList `json:"limited_bonus_common_list"`
	RandomLiveList         []RandomLiveList         `json:"random_live_list"`
	FreeLiveList           []any                    `json:"free_live_list"`
	TrainingLiveList       []TrainingLiveList       `json:"training_live_list"`
}

type ScheduleResp struct {
	Result     ScheduleData `json:"result"`
	Status     int          `json:"status"`
	CommandNum bool         `json:"commandNum"`
	TimeStamp  int64        `json:"timeStamp"`
}
