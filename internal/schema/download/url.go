package downloadschema

type UrlReq struct {
	Module   string   `json:"module"`
	Os       string   `json:"os"`
	Mgd      int      `json:"mgd"`
	PathList []string `json:"path_list"`
	Action   string   `json:"action"`
}

type UrlData struct {
	UrlList []string `json:"url_list"`
}

type UrlResp struct {
	ResponseData UrlData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
