package commonschema

import "honoka-chan/internal/constant"

type ErrorData struct {
	ErrorCode constant.ErrorCode `json:"error_code"`
}
