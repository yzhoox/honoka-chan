package download

type UpdateReq struct {
	Module          string `json:"module"`
	TargetOs        string `json:"target_os"`
	InstallVersion  string `json:"install_version"`
	TimeStamp       int    `json:"timeStamp"`
	Action          string `json:"action"`
	PackageList     []any  `json:"package_list"`
	CommandNum      string `json:"commandNum"`
	ExternalVersion string `json:"external_version"`
}

type UpdateData struct {
	Size    int    `json:"size"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

type UpdateResp struct {
	ResponseData []UpdateData `json:"response_data"`
	ReleaseInfo  []any        `json:"release_info"`
	StatusCode   int          `json:"status_code"`
}
