package routes

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/handlers"
)

func MygoRoutes(router *gin.RouterGroup) {
	router.GET("/mygo", handlers.MyGoHandler)
}
