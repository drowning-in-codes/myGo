package routes

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/handlers"
)

func DownloadPicRoutes(router *gin.RouterGroup) {
	router.GET("/refresh", handlers.DownloadPicHandler)
}
