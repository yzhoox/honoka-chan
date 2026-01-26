package costume

import (
	costumeapischema "honoka-chan/internal/schema/api/costume"
	"time"
)

func costumeList() (res any, err error) {
	res = costumeapischema.ListResp{
		Result: costumeapischema.ListData{
			CostumeList: []costumeapischema.CostumeList{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
