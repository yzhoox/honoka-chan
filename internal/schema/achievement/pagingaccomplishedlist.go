package achievementschema

import achievementapischema "honoka-chan/internal/schema/api/achievement"

type PagingAccomplishedListReq struct {
	Module           string `json:"module"`
	Action           string `json:"action"`
	TimeStamp        int64  `json:"timeStamp"`
	Mgd              int    `json:"mgd"`
	FilterCategoryID int    `json:"filter_category_id"`
	CommandNum       string `json:"commandNum"`
	FromCount        int    `json:"from_count"`
}

type PagingAccomplishedListResp struct {
	ResponseData []achievementapischema.AchievementListItem `json:"response_data"`
	ReleaseInfo  []any                                      `json:"release_info"`
	StatusCode   int                                        `json:"status_code"`
}
