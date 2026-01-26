package userapischema

import userschema "honoka-chan/internal/schema/user"

type InfoResp struct {
	Result     userschema.UserInfo `json:"result"`
	Status     int                 `json:"status"`
	CommandNum bool                `json:"commandNum"`
	TimeStamp  int64               `json:"timeStamp"`
}
