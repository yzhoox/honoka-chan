package scenario

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ScenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "scenarioStatus":
		res, err = scenarioStatus(ctx)
	default:
		err = fmt.Errorf("unimplemented action: scenario: %s", action)
	}
	return res, err
}
