package rankingschema

import (
	"encoding/json"
	profileapischema "honoka-chan/internal/schema/api/profile"
	"strconv"
)

type OptionalPage int

func (p *OptionalPage) UnmarshalJSON(data []byte) error {
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*p = OptionalPage(intValue)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return err
	}
	if stringValue == "" {
		*p = 0
		return nil
	}

	parsed, err := strconv.Atoi(stringValue)
	if err != nil {
		return err
	}
	*p = OptionalPage(parsed)
	return nil
}

type LiveReq struct {
	Module           string       `json:"module"`
	Action           string       `json:"action"`
	Mgd              int          `json:"mgd"`
	Page             OptionalPage `json:"page"`
	LiveDifficultyID string       `json:"live_difficulty_id"`
	TimeStamp        int64        `json:"timeStamp"`
	CommandNum       string       `json:"commandNum"`
}

type UserData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Level  int    `json:"level"`
}

type LiveItem struct {
	Rank           int                             `json:"rank"`
	Score          int                             `json:"score"`
	UserData       UserData                        `json:"user_data"`
	CenterUnitInfo profileapischema.CenterUnitInfo `json:"center_unit_info"`
	SettingAwardID int                             `json:"setting_award_id"`
}

type LiveData struct {
	Rank       *int       `json:"rank"`
	Items      []LiveItem `json:"items"`
	TotalCnt   int64      `json:"total_cnt"`
	PresentCnt int64      `json:"present_cnt"`
	Page       int        `json:"page"`
}

type LiveResp struct {
	ResponseData LiveData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
