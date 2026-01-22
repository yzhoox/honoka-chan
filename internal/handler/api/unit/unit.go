package unit

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func UnitApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "accessoryAll":
		res, err = unitAccessoryAll(ctx)
	case "deckInfo":
		res, err = unitDeckInfo(ctx)
	case "removableSkillInfo":
		res, err = unitRemovableSkillInfo(ctx)
	case "supporterAll":
		res, err = unitSupporterAll()
	case "unitAll":
		res, err = unitAll(ctx)
	default:
		err = fmt.Errorf("unimplemented action: unit: %s", action)
	}
	return res, err
}
