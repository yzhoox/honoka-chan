package loginapischema

type Notification struct {
	Push       bool `json:"push"`
	Lp         bool `json:"lp"`
	UpdateInfo bool `json:"update_info"`
	Campaign   bool `json:"campaign"`
	Live       bool `json:"live"`
	Lbonus     bool `json:"lbonus"`
	Event      bool `json:"event"`
	Secretbox  bool `json:"secretbox"`
	Birthday   bool `json:"birthday"`
}

type TopInfoOnceData struct {
	NewAchievementCnt            int          `json:"new_achievement_cnt"`
	UnaccomplishedAchievementCnt int          `json:"unaccomplished_achievement_cnt"`
	LiveDailyRewardExist         bool         `json:"live_daily_reward_exist"`
	TrainingEnergy               int          `json:"training_energy"`
	TrainingEnergyMax            int          `json:"training_energy_max"`
	Notification                 Notification `json:"notification"`
	OpenArena                    bool         `json:"open_arena"`
	CostumeStatus                bool         `json:"costume_status"`
	OpenAccessory                bool         `json:"open_accessory"`
	ArenaSiSkillUniqueCheck      bool         `json:"arena_si_skill_unique_check"`
	OpenV98                      bool         `json:"open_v98"`
}

type TopInfoOnceResp struct {
	Result     TopInfoOnceData `json:"result"`
	Status     int             `json:"status"`
	CommandNum bool            `json:"commandNum"`
	TimeStamp  int64           `json:"timeStamp"`
}
