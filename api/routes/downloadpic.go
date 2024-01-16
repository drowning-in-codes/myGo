package routes

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/handlers"
)

func DownloadPicRoutes(router *gin.RouterGroup) {
	router.GET("/pix", handlers.DownloadPixvisionPicHandler)
	router.GET("/booru/:type", handlers.DownloadBooruPicHandler)

}
