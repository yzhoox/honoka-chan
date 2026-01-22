package profile

import (
	"encoding/json"
	"honoka-chan/internal/schema/api/profile"
	honokautils "honoka-chan/pkg/utils"
	"time"
)

func cardRanking() (res any, err error) {
	var result []any
	love := honokautils.ReadAllText("assets/serverdata/love.json")
	err = json.Unmarshal([]byte(love), &result)
	if err != nil {
		return nil, err
	}

	res = profile.CardRankingResp{
		Result:     result,
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
