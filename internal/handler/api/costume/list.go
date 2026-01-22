package costume

import (
	"honoka-chan/internal/schema/api/costume"
	"time"
)

func costumeList() (res any, err error) {
	res = costume.ListResp{
		Result: costume.ListData{
			CostumeList: []costume.CostumeList{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
