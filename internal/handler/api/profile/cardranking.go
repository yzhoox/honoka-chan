package profile

import (
	"encoding/json"
	profileapischema "honoka-chan/internal/schema/api/profile"
	"honoka-chan/pkg/utils"
	"time"
)

func cardRanking() (res any, err error) {
	var result []any
	love := utils.ReadAllText("assets/serverdata/love.json")
	err = json.Unmarshal([]byte(love), &result)
	if err != nil {
		return nil, err
	}

	res = profileapischema.CardRankingResp{
		Result:     result,
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
