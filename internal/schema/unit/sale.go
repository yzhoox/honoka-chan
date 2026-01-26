package unitschema

import userschema "honoka-chan/internal/schema/user"

type SaleReq struct {
	Module           string `json:"module"`
	Action           string `json:"action"`
	TimeStamp        int    `json:"timeStamp"`
	Mgd              int    `json:"mgd"`
	UnitSupportList  []any  `json:"unit_support_list"`
	UnitOwningUserID []int  `json:"unit_owning_user_id"`
	CommandNum       string `json:"commandNum"`
}

type Detail struct {
	UnitOwningUserID int  `json:"unit_owning_user_id"`
	UnitID           int  `json:"unit_id"`
	IsSigned         bool `json:"is_signed"`
	Price            int  `json:"price"`
}

type OwningInfo struct {
	UnitRemovableSkillID int    `json:"unit_removable_skill_id"`
	TotalAmount          int    `json:"total_amount"`
	EquippedAmount       int    `json:"equipped_amount"`
	InsertDate           string `json:"insert_date"`
}

type UnitRemovableSkill struct {
	OwningInfo []OwningInfo `json:"owning_info"`
}

type SaleData struct {
	Total                int                 `json:"total"`
	Detail               []Detail            `json:"detail"`
	BeforeUserInfo       userschema.UserInfo `json:"before_user_info"`
	AfterUserInfo        userschema.UserInfo `json:"after_user_info"`
	RewardBoxFlag        bool                `json:"reward_box_flag"`
	GetExchangePointList []any               `json:"get_exchange_point_list"`
	UnitRemovableSkill   UnitRemovableSkill  `json:"unit_removable_skill"`
	ServerTimestamp      int64               `json:"server_timestamp"`
}

type SaleResp struct {
	ResponseData SaleData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
