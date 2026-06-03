package scenario

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
)

func reward(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	data := utils.ReadAllText("assets/serverdata/reward.json")
	var resp map[string]any
	err := json.Unmarshal([]byte(data), &resp)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(resp)
}

func init() {
	router.AddHandler("main.php", "POST", "/scenario/reward", middleware.Common, reward)
}
