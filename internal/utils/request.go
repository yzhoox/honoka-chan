package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func ParseRequestData(ctx *gin.Context, dst any) error {
	return json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), dst)
}
