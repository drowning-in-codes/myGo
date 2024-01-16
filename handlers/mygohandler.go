package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func MyGoHandler(c *gin.Context) {
	if !ok {
		log.Fatalln("请正确填写下载图片的根目录")
	}
	img_folder := filepath.Join(download_root_folder, "mygo")
	// 获取一个目录下所有目录
	err := PicExist(img_folder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有图片",
		})
	}
	random_image, err := RandomFiletoURL(img_folder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有图片",
		})
	}
	c.Redirect(302, random_image)

	c.JSON(http.StatusOK, gin.H{
		"message": "没有图片",
	})
}
