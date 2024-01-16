package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func LoadGoHandler(ctx *gin.Context) {
	if !ok {
		log.Fatalln("请正确填写下载图片的根目录")
	}
	total_img := 0
	img_folder := filepath.Join(download_root_folder, "mygo")
	os.MkdirAll(img_folder, os.ModePerm)
	allow_img_site := checkAllowSite()
	c := colly.NewCollector(colly.UserAgent(userAgent), colly.AllowedDomains(allow_img_site...),
		colly.Async())

	if proxy != nil {
		if proxy["http"] != nil {
			err := c.SetProxy(fmt.Sprintf("http://%s", proxy["http"].(string)))
			if err != nil {
				log.Panic(err.Error())
			}
		}
		if proxy["socks5"] != nil {
			err := c.SetProxy(fmt.Sprintf("socks5://%s", proxy["socks5"].(string)))
			if err != nil {
				log.Panic(err.Error())
			}
		}
	}
	c.SetRequestTimeout(10 * time.Second)

	c.OnHTML("div#posts > div[class*=central-block]", func(e *colly.HTMLElement) {
		log.Default().Println("Link found:", e.ChildAttr("a", "href"))
		e.ForEach("img[alt][src]", func(_ int, e *colly.HTMLElement) {
			c.Visit(e.ChildAttr("a", "src"))
		})
	})
	c.OnHTML("div#posts ", func(e *colly.HTMLElement) {
		log.Default().Println("Link found:", e.ChildAttr("a", "href"))

	})
	c.OnResponse(func(r *colly.Response) {
		if strings.HasPrefix(r.Headers.Get("Content-Type"), "image") {
			r.Save(filepath.Join(img_folder, r.FileName()))
			total_img += 1
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Method = "GET"
		r.Headers.Set("Accept", "*/*")
		log.Default().Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Default().Println(r.Request.URL.String())
		log.Default().Println("error", err.Error())
	})
	c.Visit(mygo_site)

	ctx.JSON(http.StatusOK, fmt.Sprintf("success! %d image Go", total_img))
}
