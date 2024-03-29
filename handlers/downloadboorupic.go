package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// type Pics struct {
// 	title string
// 	pics map[string]string
// }

func DownloadBooruPicHandler(ctx *gin.Context) {
	counter := 0
	img_site := booru_site
	allow_img_site := checkAllowSite()
	if ctx.Param("type") != "wild" {
		log.Default().Println(ctx.Param("type"))
		img_site = safebooru_site
	}

	// fmt.Println(allow_img_site)
	c := colly.NewCollector(colly.UserAgent(userAgent), colly.AllowedDomains(allow_img_site...),
		colly.Async())
	c.SetRequestTimeout(10 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		//Delay:      5 * time.Second,
		RandomDelay: 500 * time.Duration(time.Millisecond),
	})

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
	limit_page, ok := conf["limit_page"].(int)
	if !ok {
		log.Panic("爬取图片目录数配置出错")
	}

	if !ok {
		log.Panic("下载路径配置出错")
	}
	var my_folder string
	if ctx.Param("type") == "wild" {
		my_folder = filepath.Join(download_root_folder, "wildbooru")
	} else {
		my_folder = filepath.Join(download_root_folder, "safebooru")

	}
	os.MkdirAll(my_folder, os.ModePerm)
	// Find and visit all links
	c.OnHTML("div.content > div:not([class])", func(e *colly.HTMLElement) {
		e.ForEach("a[id]", func(_ int, el *colly.HTMLElement) {
			el.Request.Visit(el.Attr("href"))
		})
	})
	c.OnHTML("div#note-container", func(e *colly.HTMLElement) {
		img_ele := e.DOM.SiblingsFiltered("img[alt][src]")
		img_src, _ := img_ele.Attr("src")
		if img_src != "" {
			e.Request.Visit(img_src)
		}
		img_ele = e.DOM.ChildrenFiltered("img[alt][src]")
		img_src, _ = img_ele.Attr("src")
		if img_src != "" {
			e.Request.Visit(img_src)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.Method = "GET"
		log.Default().Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		if strings.HasPrefix(r.Headers.Get("Content-Type"), "image") {
			filename := r.FileName()
			p := regexp.MustCompile(`_\d+$`)
			filename = p.ReplaceAllString(filename, "")
			r.Save(filepath.Join(my_folder, filename))
			counter += 1
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Default().Printf("visit %s error %s\n", r.Request.URL.String(), err.Error())
	})
	for i := 0; i < limit_page; i++ {
		c.Visit(img_site)
	}
	c.Wait()
	ctx.JSON(http.StatusOK, fmt.Sprintf("success! %d image go Load", counter))
}
