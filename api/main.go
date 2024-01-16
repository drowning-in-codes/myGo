package main

import (
	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/routes"
)

func main() {

	router := gin.Default()

	r := router.Group("/api")
	router.Static("/img", "./imgs")
	routes.DownloadPicRoutes(r)
	routes.ShowPicRoutes(r)
	routes.GetPicRoutes(r)

	router.Run(":8080")

}
