package routes

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/handlers"
)

func GetPicRoutes(router *gin.RouterGroup) {
	router.GET("/get", handlers.GetPicHandler)
}
