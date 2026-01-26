package museum

import (
	museumapischema "honoka-chan/internal/schema/api/museum"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func museumInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var museumRes []struct {
		MuseumContentsId int `xorm:"museum_contents_id"`
		SmileBuff        int `xorm:"smile_buff"`
		PureBuff         int `xorm:"pure_buff"`
		CoolBuff         int `xorm:"cool_buff"`
	}
	var museumID []int
	var smileBuff, pureBuff, coolBuff int
	err = ss.MainEng.Table("museum_contents_m").Cols("museum_contents_id,smile_buff,pure_buff,cool_buff").
		OrderBy("museum_contents_id ASC").Find(&museumRes)
	if ss.CheckErr(err) {
		return
	}

	for _, res := range museumRes {
		smileBuff += res.SmileBuff
		pureBuff += res.PureBuff
		coolBuff += res.CoolBuff
		museumID = append(museumID, res.MuseumContentsId)
	}
	res = museumapischema.InfoResp{
		Result: museumapischema.InfoData{
			MuseumInfo: museumapischema.Info{
				Parameter: museumapischema.Parameter{
					Smile: smileBuff,
					Pure:  pureBuff,
					Cool:  coolBuff,
				},
				ContentsIDList: museumID,
			},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
