package costume

import (
	costumeapischema "honoka-chan/internal/schema/api/costume"
	"net/http"
	"time"
)

func costumeList() (res any, err error) {
	res = costumeapischema.ListResp{
		Result: costumeapischema.ListData{
			CostumeList: []costumeapischema.CostumeList{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
