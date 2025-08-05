package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

type MuseumContent struct {
	MuseumContentsId int `xorm:"museum_contents_id"`
	SmileBuff        int `xorm:"smile_buff"`
	PureBuff         int `xorm:"pure_buff"`
	CoolBuff         int `xorm:"cool_buff"`
}

func MuseumInfo(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	var contents []MuseumContent
	err := ss.MainEng.Table("museum_contents_m").Cols("museum_contents_id,smile_buff,pure_buff,cool_buff").Find(&contents)
	if ss.CheckErr(err) {
		return
	}

	var smileBuff, pureBuff, coolBuff int
	var contentsList []int
	for _, content := range contents {
		smileBuff += content.SmileBuff
		pureBuff += content.PureBuff
		coolBuff += content.CoolBuff
		contentsList = append(contentsList, content.MuseumContentsId)
	}
	museumResp := model.MuseumResp{
		ResponseData: model.MuseumRes{
			MuseumInfo: model.Museum{
				Parameter: model.MuseumParameter{
					Smile: smileBuff,
					Pure:  pureBuff,
					Cool:  coolBuff,
				},
				ContentsIDList: contentsList,
			},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(museumResp)
}
