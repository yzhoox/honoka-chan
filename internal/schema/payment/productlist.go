package paymentschema

import paymentapischema "honoka-chan/internal/schema/api/payment"

type ProductListData struct {
	RestrictionInfo  paymentapischema.RestrictionInfo `json:"restriction_info"`
	UnderAgeInfo     paymentapischema.UnderAgeInfo    `json:"under_age_info"`
	SnsProductList   []any                            `json:"sns_product_list"`
	ProductList      []any                            `json:"product_list"`
	SubscriptionList []any                            `json:"subscription_list"`
	ShowPointShop    bool                             `json:"show_point_shop"`
	ServerTimestamp  int64                            `json:"server_timestamp"`
}

type ProductListResp struct {
	ResponseData ProductListData `json:"response_data"`
	ReleaseInfo  []any           `json:"release_info"`
	StatusCode   int             `json:"status_code"`
}
