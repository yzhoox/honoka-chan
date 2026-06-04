package profile

import (
	profileapischema "honoka-chan/internal/schema/api/profile"
	honokautils "honoka-chan/internal/utils"
	"time"
)

func cardRanking() (res any, err error) {
	cardRankingData, err := honokautils.LoadServerData[[]profileapischema.CardRankingData]("card_ranking_data.json")
	if err != nil {
		return nil, err
	}

	res = profileapischema.CardRankingResp{
		Result:     cardRankingData,
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
