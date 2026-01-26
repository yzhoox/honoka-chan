package downloadschema

type AdditionalReq struct {
	Module      string `json:"module"`
	Mgd         int    `json:"mgd"`
	Action      string `json:"action"`
	TimeStamp   int    `json:"timeStamp"`
	PackageID   int    `json:"package_id"`
	TargetOs    string `json:"target_os"`
	PackageType int    `json:"package_type"`
	CommandNum  string `json:"commandNum"`
}

type AdditionalData struct {
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type AdditionalResp struct {
	ResponseData []AdditionalData `json:"response_data"`
	ReleaseInfo  []any            `json:"release_info"`
	StatusCode   int              `json:"status_code"`
}
