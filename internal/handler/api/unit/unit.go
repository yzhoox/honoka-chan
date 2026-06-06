package unit

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func UnitApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "accessoryAll":
		res, err = unitAccessoryAll(ctx)
	case "accessoryMaterialAll":
		res, err = unitAccessoryMaterialAll()
	case "accessoryTab":
		res, err = unitAccessoryTab()
	case "deckInfo":
		res, err = unitDeckInfo(ctx)
	case "removableSkillInfo":
		res, err = unitRemovableSkillInfo(ctx)
	case "supporterAll":
		res, err = unitSupporterAll()
	case "unitAll":
		res, err = unitAll(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("unit", action)
	}
	return res, err
}
