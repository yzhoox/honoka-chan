package live

type NormalLiveStatusList struct {
	LiveDifficultyID   int   `json:"live_difficulty_id"`
	Status             int   `json:"status"`
	HiScore            int   `json:"hi_score"`
	HiComboCount       int   `json:"hi_combo_count"`
	ClearCnt           int   `json:"clear_cnt"`
	AchievedGoalIDList []int `json:"achieved_goal_id_list"`
}

type SpecialLiveStatusList struct {
	LiveDifficultyID   int   `json:"live_difficulty_id"`
	Status             int   `json:"status"`
	HiScore            int   `json:"hi_score"`
	HiComboCount       int   `json:"hi_combo_count"`
	ClearCnt           int   `json:"clear_cnt"`
	AchievedGoalIDList []int `json:"achieved_goal_id_list"`
}

type TrainingLiveStatusList struct {
	LiveDifficultyID   int   `json:"live_difficulty_id"`
	Status             int   `json:"status"`
	HiScore            int   `json:"hi_score"`
	HiComboCount       int   `json:"hi_combo_count"`
	ClearCnt           int   `json:"clear_cnt"`
	AchievedGoalIDList []int `json:"achieved_goal_id_list"`
}

type StatusData struct {
	NormalLiveStatusList   []NormalLiveStatusList   `json:"normal_live_status_list"`
	SpecialLiveStatusList  []SpecialLiveStatusList  `json:"special_live_status_list"`
	TrainingLiveStatusList []TrainingLiveStatusList `json:"training_live_status_list"`
	MarathonLiveStatusList []any                    `json:"marathon_live_status_list"`
	FreeLiveStatusList     []any                    `json:"free_live_status_list"`
	CanResumeLive          bool                     `json:"can_resume_live"`
}

type StatusResp struct {
	Result     StatusData `json:"result"`
	Status     int        `json:"status"`
	CommandNum bool       `json:"commandNum"`
	TimeStamp  int64      `json:"timeStamp"`
}
