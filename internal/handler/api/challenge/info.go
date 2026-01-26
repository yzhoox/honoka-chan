package challenge

import (
	challengeapischema "honoka-chan/internal/schema/api/challenge"
	"time"
)

func challengeInfo() (res any, err error) {
	res = challengeapischema.InfoResp{
		Result:     []any{},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}
	return res, err
}
