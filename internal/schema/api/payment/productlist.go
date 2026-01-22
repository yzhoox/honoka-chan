package payment

type RestrictionInfo struct {
	Restricted bool `json:"restricted"`
}

type UnderAgeInfo struct {
	BirthSet    bool `json:"birth_set"`
	HasLimit    bool `json:"has_limit"`
	LimitAmount any  `json:"limit_amount"`
	MonthUsed   int  `json:"month_used"`
}

type SnsProductItem struct {
	ItemID    int  `json:"item_id"`
	AddType   int  `json:"add_type"`
	Amount    int  `json:"amount"`
	IsFreebie bool `json:"is_freebie"`
}

type SnsProduct struct {
	ProductID   string           `json:"product_id"`
	Name        string           `json:"name"`
	Price       int              `json:"price"`
	CanBuy      bool             `json:"can_buy"`
	ProductType int              `json:"product_type"`
	ItemList    []SnsProductItem `json:"item_list"`
}

type ProductItem struct {
	ItemID    int  `json:"item_id"`
	AddType   int  `json:"add_type"`
	Amount    int  `json:"amount"`
	IsFreebie bool `json:"is_freebie"`
	IsRankMax bool `json:"is_rank_max,omitempty"`
}

type LimitStatus struct {
	TermStartDate  string `json:"term_start_date"`
	RemainingTime  string `json:"remaining_time"`
	RemainingCount int    `json:"remaining_count"`
}

type Product struct {
	ProductID      string        `json:"product_id"`
	Name           string        `json:"name"`
	BannerImgAsset string        `json:"banner_img_asset"`
	Price          int           `json:"price"`
	CanBuy         bool          `json:"can_buy"`
	ProductType    int           `json:"product_type"`
	AnnounceURL    string        `json:"announce_url"`
	ConfirmURL     string        `json:"confirm_url"`
	ItemList       []ProductItem `json:"item_list"`
	LimitStatus    LimitStatus   `json:"limit_status"`
}

type SubscriptionItem struct {
	ItemID    int  `json:"item_id"`
	AddType   int  `json:"add_type"`
	Amount    int  `json:"amount"`
	IsFreebie bool `json:"is_freebie"`
}

type LicenseReward struct {
	ItemID  int `json:"item_id"`
	AddType int `json:"add_type"`
	Amount  int `json:"amount"`
}

type LicenseItem struct {
	Seq        int             `json:"seq"`
	RewardList []LicenseReward `json:"reward_list"`
}

type LicenseInfo struct {
	Name  string        `json:"name"`
	Items []LicenseItem `json:"items"`
}

type UserStatus struct {
	IsLicensed bool `json:"is_licensed"`
}

type SubscriptionStatus struct {
	LicenseID     int         `json:"license_id"`
	LicenseType   int         `json:"license_type"`
	LicenseInfo   LicenseInfo `json:"license_info,omitempty"`
	UserStatus    UserStatus  `json:"user_status"`
	PurchaseCount int         `json:"purchase_count"`
	BadgeFlag     bool        `json:"badge_flag"`
}

type Subscription struct {
	ProductID          string             `json:"product_id"`
	Name               string             `json:"name"`
	BannerImgAsset     string             `json:"banner_img_asset"`
	Price              int                `json:"price"`
	CanBuy             bool               `json:"can_buy"`
	ProductType        int                `json:"product_type"`
	ProductURL         string             `json:"product_url"`
	ItemList           []SubscriptionItem `json:"item_list"`
	LimitStatus        LimitStatus        `json:"limit_status"`
	SubscriptionStatus SubscriptionStatus `json:"subscription_status"`
}

type ProductListData struct {
	RestrictionInfo  RestrictionInfo `json:"restriction_info"`
	UnderAgeInfo     UnderAgeInfo    `json:"under_age_info"`
	SnsProductList   []SnsProduct    `json:"sns_product_list"`
	ProductList      []Product       `json:"product_list"`
	SubscriptionList []Subscription  `json:"subscription_list"`
	ShowPointShop    bool            `json:"show_point_shop"`
}

type ProductListResp struct {
	Result     ProductListData `json:"result"`
	Status     int             `json:"status"`
	CommandNum bool            `json:"commandNum"`
	TimeStamp  int64           `json:"timeStamp"`
}
