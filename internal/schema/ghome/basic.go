package ghome

type GetCodeResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data map[string]any `json:"data"`
}

type GetProductListResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data GetProductListData `json:"data"`
}

type GetProductListData struct {
	Message []string `json:"message"`
	Result  int      `json:"result"`
}

type HandshakeResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type LoginAreaData struct {
	UserID string `json:"userid"`
}

type LoginAreaResp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data LoginAreaData `json:"data"`
}

type PublicKeyData struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
	Method  string `json:"method,omitempty"`
}

type PublicKeyResp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data PublicKeyData `json:"data"`
}
