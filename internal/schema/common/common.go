package commonschema

import (
	"honoka-chan/internal/constant"
)

type ErrorData struct {
	ErrorCode constant.ErrorCode `json:"error_code"`
}

type ErrorResp struct {
	ResponseData ErrorData `json:"response_data"`
	ReleaseInfo  []any     `json:"release_info"`
	StatusCode   int       `json:"status_code"`
}

type ApiErrorResp struct {
	Result     ErrorData `json:"result"`
	Status     int       `json:"status"`
	CommandNum bool      `json:"commandNum"`
	TimeStamp  int64     `json:"timeStamp"`
}
