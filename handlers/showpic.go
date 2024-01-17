package handlers

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RandomShowPicHandler(c *gin.Context) {
	if !ok {
		log.Fatalln("请正确填写下载图片的根目录")
	}

	entries, error := os.ReadDir(download_root_folder)
	if error != nil {
		log.Default().Println("下载目录读取错误")
	}
	if len(entries) == 0 {
		c.JSON(404, gin.H{"message": "没有图片"})
	}
	random_dir := entries[rand.Intn(len(entries))]
	abs_random_dir := filepath.Join(download_root_folder, random_dir.Name())
	os.MkdirAll(abs_random_dir, os.ModePerm)
	if err := PicExist(abs_random_dir); err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}

	// 随机获取一张图片
	randomPicURL, err := RandomFiletoURL(abs_random_dir)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}
	c.Redirect(302, randomPicURL)
	// c.Redirect(302,strings.Join([]string{"/api/img",url.PathEscape(random_dir.Name()),url.PathEscape(random_img)},"/") )
}

func TypeShowPicHandler(c *gin.Context) {
	if !ok {
		log.Fatalln("请正确填写下载图片的根目录")
	}
	_, error := os.ReadDir(download_root_folder)
	if error != nil {
		log.Default().Println("下载目录读取错误")
	}
	t := c.Param("type")
	var navURL string
	if t == "wild" {
		wildBooru := filepath.Join(download_root_folder, "wildbooru")
		err := PicExist(wildBooru)
		if err != nil {
			c.JSON(404, gin.H{"message": "没有这个类型的图片"})
		}
		os.MkdirAll(wildBooru, os.ModePerm)
		navURL = filepath.Join("./img", "wildbooru")
	} else {
		safeBooru := filepath.Join(download_root_folder, "wildbooru")
		err := PicExist(safeBooru)
		if err != nil {
			c.JSON(404, gin.H{"message": "没有这个类型的图片"})
		}
		os.MkdirAll(safeBooru, os.ModePerm)
		navURL = filepath.Join("./img", "safebooru")
	}

	randomPicURL, err := RandomFiletoURL(navURL)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}
	c.Redirect(302, randomPicURL)
	// c.Redirect(302,strings.Join([]string{"/api/img",url.PathEscape(random_dir.Name()),url.PathEscape(random_img)},"/") )
}
