package eventscenario

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func EventScenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "status":
		res, err = eventScenarioStatus(ctx)
	default:
		err = fmt.Errorf("unimplemented action: costume: %s", action)
	}
	return res, err
}
