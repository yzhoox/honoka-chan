package subscenario

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SubscenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "subscenarioStatus":
		res, err = subscenarioStatus(ctx)
	default:
		err = fmt.Errorf("unimplemented action: subscenario: %s", action)
	}
	return res, err
}
