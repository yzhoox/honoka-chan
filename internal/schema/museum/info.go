package museum

type Parameter struct {
	Smile int `json:"smile"`
	Pure  int `json:"pure"`
	Cool  int `json:"cool"`
}

type Museum struct {
	Parameter      Parameter `json:"parameter"`
	ContentsIDList []int     `json:"contents_id_list"`
}

type InfoData struct {
	MuseumInfo      Museum `json:"museum_info"`
	ServerTimestamp int64  `json:"server_timestamp"`
}

type InfoResp struct {
	ResponseData InfoData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
