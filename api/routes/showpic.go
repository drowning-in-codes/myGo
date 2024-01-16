package routes

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/handlers"
)

func ShowPicRoutes(router *gin.RouterGroup) {
	router.GET("/show", handlers.RandomShowPicHandler)
	router.GET("/show/:type", handlers.TypeShowPicHandler)
}
