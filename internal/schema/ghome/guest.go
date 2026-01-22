package ghome

type GuestStatusData struct {
	Disablead   int    `json:"disablead"`
	Loginswitch int    `json:"loginswitch"`
	Message     string `json:"message"`
	Result      int    `json:"result"`
}

type GuestStatusResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data GuestStatusData `json:"data"`
}
