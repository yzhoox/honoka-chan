package login

type AuthKeyData struct {
	AuthorizeToken string `json:"authorize_token"`
	DummyToken     string `json:"dummy_token"`
}

type AuthKeyResp struct {
	ResponseData AuthKeyData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
