package subscenario

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
)

func reward(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	data := honokautils.ReadAllText("assets/serverdata/subreward.json")
	var resp map[string]any
	err := json.Unmarshal([]byte(data), &resp)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(resp)
}

func init() {
	router.AddHandler("main.php", "POST", "/subscenario/reward", middleware.Common, reward)
}
