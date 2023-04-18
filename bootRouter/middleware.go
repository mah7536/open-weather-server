package bootRouter

import (
	"alarm-system/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func middleware(ctx *gin.Context) {
	if !config.Power {
		ctx.JSON(http.StatusNotFound, nil)
		ctx.Abort()
		return
	}
	ctx.Next()
}
