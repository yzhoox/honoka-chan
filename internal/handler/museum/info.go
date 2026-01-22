package museum

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/museum"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func info(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	var contents []struct {
		MuseumContentsId int `xorm:"museum_contents_id"`
		SmileBuff        int `xorm:"smile_buff"`
		PureBuff         int `xorm:"pure_buff"`
		CoolBuff         int `xorm:"cool_buff"`
	}
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
	museumResp := museum.InfoResp{
		ResponseData: museum.InfoData{
			MuseumInfo: museum.Museum{
				Parameter: museum.Parameter{
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

func init() {
	router.AddHandler("main.php", "POST", "/museum/info", middleware.Common, info)
}
