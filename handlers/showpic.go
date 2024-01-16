package handlers

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RandomShowPicHandler(c *gin.Context) {
	entries, error := os.ReadDir("./imgs")
	if error != nil {
		log.Panic(error.Error())
	}
	if len(entries) == 0 {
		c.JSON(404, gin.H{"message": "没有图片"})
	}
	random_dir := entries[rand.Intn(len(entries))]
	abs_random_dir := filepath.Join("./imgs", random_dir.Name())
	if err := PicExist(abs_random_dir); err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}

	// 随机获取一张图片
	randomPicURL, err := RandomFile(abs_random_dir)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}
	c.Redirect(302, randomPicURL)
	// c.Redirect(302,strings.Join([]string{"/api/img",url.PathEscape(random_dir.Name()),url.PathEscape(random_img)},"/") )
}

func TypeShowPicHandler(c *gin.Context) {
	_, error := os.ReadDir("./imgs")
	if error != nil {
		log.Panic(error.Error())
	}
	t := c.Param("type")
	var navURL string
	if t == "wild" {
		wildBooru := filepath.Join("./imgs", "wildbooru")
		err := PicExist(wildBooru)
		if err != nil {
			c.JSON(404, gin.H{"message": "没有这个类型的图片"})
		}
		navURL = filepath.Join("./img", "wildbooru")
	} else {
		safeBooru := filepath.Join("./imgs", "wildbooru")
		err := PicExist(safeBooru)
		if err != nil {
			c.JSON(404, gin.H{"message": "没有这个类型的图片"})
		}
		navURL = filepath.Join("./img", "safebooru")
	}

	randomPicURL, err := RandomFile(navURL)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	}
	c.Redirect(302, randomPicURL)
	// c.Redirect(302,strings.Join([]string{"/api/img",url.PathEscape(random_dir.Name()),url.PathEscape(random_img)},"/") )
}
