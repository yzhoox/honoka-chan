package user

type LpRecoveryItem struct {
	ItemID int `json:"item_id"`
	Amount int `json:"amount"`
}

type UserInfo struct {
	UserID                         int              `json:"user_id"`
	Name                           string           `json:"name"`
	Level                          int              `json:"level"`
	Exp                            int              `json:"exp"`
	PreviousExp                    int              `json:"previous_exp"`
	NextExp                        int              `json:"next_exp"`
	GameCoin                       int              `json:"game_coin"`
	SnsCoin                        int              `json:"sns_coin"`
	FreeSnsCoin                    int              `json:"free_sns_coin"`
	PaidSnsCoin                    int              `json:"paid_sns_coin"`
	SocialPoint                    int              `json:"social_point"`
	UnitMax                        int              `json:"unit_max"`
	WaitingUnitMax                 int              `json:"waiting_unit_max"`
	EnergyMax                      int              `json:"energy_max"`
	EnergyFullTime                 string           `json:"energy_full_time"`
	LicenseLiveEnergyRecoverlyTime int              `json:"license_live_energy_recoverly_time"`
	EnergyFullNeedTime             int              `json:"energy_full_need_time"`
	OverMaxEnergy                  int              `json:"over_max_energy"`
	TrainingEnergy                 int              `json:"training_energy"`
	TrainingEnergyMax              int              `json:"training_energy_max"`
	FriendMax                      int              `json:"friend_max"`
	InviteCode                     string           `json:"invite_code"`
	InsertDate                     string           `json:"insert_date"`
	UpdateDate                     string           `json:"update_date"`
	TutorialState                  int              `json:"tutorial_state"`
	DiamondCoin                    int              `json:"diamond_coin"`
	CrystalCoin                    int              `json:"crystal_coin"`
	LpRecoveryItem                 []LpRecoveryItem `json:"lp_recovery_item"`
}

type InfoResp struct {
	Result     UserInfo `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
