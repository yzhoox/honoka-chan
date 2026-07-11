package challenge

import (
	challengeapischema "honoka-chan/internal/schema/api/challenge"
	"net/http"
	"time"
)

func challengeInfo() (res any, err error) {
	res = challengeapischema.InfoResp{
		Result:     []any{},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}
	return res, err
}
