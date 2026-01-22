package download

type BatchReq struct {
	ClientVersion      string `json:"client_version"`
	Os                 string `json:"os"`
	PackageType        int    `json:"package_type"`
	ExcludedPackageIds []int  `json:"excluded_package_ids"`
	CommandNum         string `json:"commandNum"`
}

type BatchData struct {
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type BatchResp struct {
	ResponseData []BatchData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
