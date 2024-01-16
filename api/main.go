package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sekyoro.top/Goimg/routes"
)
var srv http.Handler

func Ф(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
func init() {
	router := gin.Default()

	r := router.Group("/api")
	router.Static("/img", "./imgs")
	routes.DownloadPicRoutes(r)
	routes.ShowPicRoutes(r)
	routes.GetPicRoutes(r)

	router.Run(":8080")
	srv = router
}

