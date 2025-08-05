package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"honoka-chan/internal/model"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/encrypt"
	honokautils "honoka-chan/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SubScenarioStartup(ctx *gin.Context) {
	startReq := model.SubScenarioReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &startReq)
	utils.CheckErr(err)

	startResp := model.SubScenarioResp{
		ResponseData: model.SubScenarioRes{
			SubscenarioID:      startReq.SubscenarioID,
			ScenarioAdjustment: 50,
			ServerTimestamp:    time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(startResp)
	utils.CheckErr(err)

	nonce := ctx.GetInt("nonce")
	nonce++

	ctx.Header("user_id", ctx.GetString("userid"))
	ctx.Header("authorize", fmt.Sprintf("consumerKey=lovelive_test&timeStamp=%d&version=1.1&token=%s&nonce=%d&user_id=%s&requestTimeStamp=%d", time.Now().Unix(), ctx.GetString("token"), nonce, ctx.GetString("userid"), ctx.GetInt64("req_time")))
	ctx.Header("X-Message-Sign", base64.StdEncoding.EncodeToString(encrypt.RSASignSHA1(resp)))

	ctx.String(http.StatusOK, string(resp))
}

func SubScenarioReward(ctx *gin.Context) {
	resp := honokautils.ReadAllText("assets/serverdata/subreward.json")

	ctx.Header("X-Message-Sign", utils.GenXMS([]byte(resp)))
	ctx.String(http.StatusOK, resp)
}
