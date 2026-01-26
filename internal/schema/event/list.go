package eventschema

type TargetList struct {
	Position      int  `json:"position"`
	IsDisplayable bool `json:"is_displayable"`
}

type ListData struct {
	TargetList      []TargetList `json:"target_list"`
	ServerTimestamp int64        `json:"server_timestamp"`
}

type ListResp struct {
	ResponseData any   `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
