package rliveschema

type TokenReq struct {
	Module     string `json:"module"`
	Action     string `json:"action"`
	Mgd        int    `json:"mgd"`
	Token      string `json:"token"`
	TimeStamp  int64  `json:"timeStamp"`
	CommandNum string `json:"commandNum"`
}
