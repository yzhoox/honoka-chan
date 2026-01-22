package live

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
)

func partyList(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	data := honokautils.ReadAllText("assets/serverdata/partylist.json")
	var partyResp map[string]any
	err := json.Unmarshal([]byte(data), &partyResp)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(partyResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/partyList", middleware.Common, partyList)
}
