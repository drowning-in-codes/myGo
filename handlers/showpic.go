package handlers

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


func ShowPicHandler(c *gin.Context) {
	entries,error:=os.ReadDir("./imgs")
	if error!=nil {
		log.Panic(error.Error())
	}
	log.Default().Println(entries)
	random_dir := entries[rand.Intn(len(entries))]
	abs_random_dir := filepath.Join("./imgs",random_dir.Name())

	all_img,error := os.ReadDir(abs_random_dir)
	if error!=nil {
		log.Panic(error.Error())
	}
	random_img := all_img[rand.Intn(len(all_img))].Name()
	abs_random_img := filepath.Join(abs_random_dir,random_img)
	log.Default().Println(abs_random_img)
	c.Redirect(302,filepath.Join("img",random_dir.Name(),random_img))
	// c.Redirect(302,strings.Join([]string{"/api/img",url.PathEscape(random_dir.Name()),url.PathEscape(random_img)},"/") )

}