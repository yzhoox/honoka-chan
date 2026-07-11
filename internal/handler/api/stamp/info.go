package stamp

import (
	stampapischema "honoka-chan/internal/schema/api/stamp"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"
)

func stampInfo() (res any, err error) {
	stampData, err := honokautils.LoadServerData[stampapischema.InfoData]("stamp_data.json")
	if err != nil {
		return nil, err
	}

	res = stampapischema.InfoResp{
		Result:     stampData,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
