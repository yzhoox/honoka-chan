package user

import "honoka-chan/internal/schema/api/user"

type Birth struct {
	BirthMonth int `json:"birth_month"`
	BirthDay   int `json:"birth_day"`
}

type UserInfoData struct {
	User            user.UserInfo `json:"user"`
	Birth           Birth         `json:"birth"`
	ServerTimestamp int64         `json:"server_timestamp"`
}

type UserInfoResp struct {
	ResponseData UserInfoData `json:"response_data"`
	ReleaseInfo  []any        `json:"release_info"`
	StatusCode   int          `json:"status_code"`
}
