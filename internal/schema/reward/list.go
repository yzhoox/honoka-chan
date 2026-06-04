package rewardschema

type ListData struct {
	ItemCount int   `json:"item_count"`
	Limit     int   `json:"limit"`
	Order     int   `json:"order"`
	Items     []any `json:"items"`
}

type ListResp struct {
	ResponseData ListData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
