package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *TimelineHandler) {
	router.GET("/timeline", handler.GetTimeline)
}
