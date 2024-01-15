package main

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/routes"
)

func main() {
	router := gin.Default()
	routes.DownloadPicRoutes(router)
	routes.ShowPicRoutes(router)

	router.Run(":8080")
}
