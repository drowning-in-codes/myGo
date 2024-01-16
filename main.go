package main

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/routes"
)


func main() {
	router := gin.Default()
	r := router.Group("/api")
	r.Static("/img","./imgs")
	// routes.DownloadPicRoutes(router)
	routes.ShowPicRoutes(r)

	router.Run(":8080")
}