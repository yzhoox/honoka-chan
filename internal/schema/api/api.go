package apischema

type ApiReq struct {
	Module      string `json:"module"`
	UserID      int    `json:"user_id"`
	Action      string `json:"action"`
	ID          int    `json:"id"`
	SecretBoxID int    `json:"secret_box_id"`
	Timestamp   int64  `json:"timeStamp"`
}

type ApiResp struct {
	ResponseData any `json:"response_data"`
	ReleaseInfo  any `json:"release_info"`
	StatusCode   int `json:"status_code"`
}
