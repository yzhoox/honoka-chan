package banner

import (
	bannerapischema "honoka-chan/internal/schema/api/banner"
	"time"
)

func bannerList() (res any, err error) {
	res = bannerapischema.ListResp{
		Result: bannerapischema.ListData{
			TimeLimit: "2037-12-31 23:59:59",
			BannerList: []bannerapischema.BannerList{
				{
					BannerType:       1,
					TargetID:         1743,
					AssetPath:        "assets/image/secretbox/icon/s_ba_1718_1.png",
					FixedFlag:        false,
					BackSide:         false,
					BannerID:         101151,
					StartDate:        "2013-04-15 00:00:00",
					EndDate:          "2037-12-31 23:59:59",
					AddUnitStartDate: "2022-01-01 00:00:00",
				},
				{
					BannerType:       1,
					TargetID:         1741,
					AssetPath:        "assets/image/secretbox/icon/s_ba_1719_1.png",
					FixedFlag:        false,
					BackSide:         false,
					BannerID:         101150,
					StartDate:        "2013-04-15 00:00:00",
					EndDate:          "2037-12-31 23:59:59",
					AddUnitStartDate: "2022-01-01 00:00:00",
				},
				{
					BannerType:       1,
					TargetID:         1740,
					AssetPath:        "assets/image/secretbox/icon/s_ba_1720_1.png",
					FixedFlag:        false,
					BackSide:         false,
					BannerID:         101149,
					StartDate:        "2013-04-15 00:00:00",
					EndDate:          "2037-12-31 23:59:59",
					AddUnitStartDate: "2022-01-01 00:00:00",
				},
				{
					BannerType:       1,
					TargetID:         1739,
					AssetPath:        "assets/image/secretbox/icon/s_ba_1721_1.png",
					FixedFlag:        false,
					BackSide:         false,
					BannerID:         101144,
					StartDate:        "2013-04-15 00:00:00",
					EndDate:          "2037-12-31 23:59:59",
					AddUnitStartDate: "2022-01-01 00:00:00",
				},
				{
					BannerType: 2,
					TargetID:   1,
					AssetPath:  "assets/image/webview/wv_ba_01.png",
					WebviewURL: "/manga",
					FixedFlag:  false,
					BackSide:   true,
					BannerID:   200001,
					StartDate:  "2016-10-15 15:00:00",
					EndDate:    "2037-12-31 23:59:59",
				},
			},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
