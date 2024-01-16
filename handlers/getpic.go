package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetPicHandler(c *gin.Context) {
	if !ok {
		log.Fatalln("请正确填写下载图片的根目录")
	}

	// 获取一个目录下所有目录
	entries, err := os.ReadDir(download_root_folder)
	if err != nil {
		log.Panic(err.Error())
	}
	total_len := len(entries)
	if total_len == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有图片",
		})
	} else {
		// 随机在一个目录下获取随机一张图片
		random_dir := entries[rand.Intn(total_len)]
		entries, err = os.ReadDir(filepath.Join(download_root_folder, random_dir.Name()))
		if err != nil {
			log.Panic(err.Error())
		}
		total_len = len(entries)
		if total_len == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "没有图片",
			})
			return
		} else {
			random_pic := entries[rand.Intn(total_len)]
			c.File(filepath.Join(download_root_folder, random_dir.Name(), random_pic.Name()))
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "没有图片",
	})
}
