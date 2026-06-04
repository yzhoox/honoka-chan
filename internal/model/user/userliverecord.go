package usermodel

type UserLiveRecord struct {
	ID               int    `xorm:"id pk autoincr"`
	LiveDifficultyID int    `xorm:"live_difficulty_id"`
	DeckID           int    `xorm:"deck_id"`
	TotalScore       int    `xorm:"total_score"` //
	PerfectCnt       int    `xorm:"perfect_cnt"` //
	GreatCnt         int    `xorm:"great_cnt"`   //
	GoodCnt          int    `xorm:"good_cnt"`    //
	BadCnt           int    `xorm:"bad_cnt"`     //
	MissCnt          int    `xorm:"miss_cnt"`    //
	MaxCombo         int    `xorm:"max_combo"`
	IsSkillOn        bool   `xorm:"is_skill_on"`
	LiveInfoJSON     string `xorm:"live_info_json"`
	PreciseListJSON  string `xorm:"precise_list_json"`
	DeckInfoJSON     string `xorm:"deck_info_json"`
	TapAdjust        int    `xorm:"tap_adjust"`
	LiveSettingJSON  string `xorm:"live_setting_json"`
	CanReplay        bool   `xorm:"can_replay"`
	TriggerLogJSON   string `xorm:"trigger_log_json"`
	UserID           int    `xorm:"user_id"`
	UpdateDate       string `xorm:"insert_date"`
}

func (UserLiveRecord) TableName() string {
	return "user_live_record"
}
