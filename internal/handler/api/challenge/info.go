package challenge

import (
	"honoka-chan/internal/schema/api/challenge"
	"time"
)

func challengeInfo() (res any, err error) {
	res = challenge.InfoResp{
		Result:     []any{},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}
	return res, err
}
