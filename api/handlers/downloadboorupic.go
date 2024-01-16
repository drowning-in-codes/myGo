package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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
	img_site := booru_site
	if ctx.Param("type") != "wild" {
		log.Default().Println(ctx.Param("type"))
		img_site = safebooru_site
	}
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	conf := loadConf("./configure.yaml")
	proxy := conf["proxy"].(map[string]interface{})
	allow_img_site := make([]string, 0)
	value := reflect.ValueOf(conf["allow_site"])
	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		// implement
		for i := 0; i < value.Len(); i++ {
			// fmt.Println(value.Index[i])
			allow_img_site = append(allow_img_site, value.Index(i).Interface().(string))
		}

	} else {
		log.Panic("请填写允许爬取网站域名字符串")
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
			err := c.SetProxy(fmt.Sprintf("http:%s", proxy["http"].(string)))
			if err != nil {
				log.Panic(err.Error())
			}
		}
		if proxy["socks5"] != nil {
			err := c.SetProxy(fmt.Sprintf("socks5:%s", proxy["socks5"].(string)))
			if err != nil {
				log.Panic(err.Error())
			}
		}
	}
	limit_page, ok := conf["limit_page"].(int)
	if !ok {
		log.Panic("爬取图片目录数配置出错")
	}
	download_root_folder, ok := conf["root_folder"].(string)
	if !ok {
		log.Panic("下载路径配置出错")
	}
	if ctx.Param("type") == "wild" {
		download_root_folder = filepath.Join(download_root_folder, "wildbooru")
	} else {
		download_root_folder = filepath.Join(download_root_folder, "safebooru")

	}
	os.MkdirAll(download_root_folder, os.ModePerm)
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
			log.Default().Printf("download success! %s\n", filename)
			r.Save(filepath.Join(download_root_folder, filename))
		}
	})
	for i := 0; i < limit_page; i++ {
		c.Visit(img_site)
	}
	c.Wait()
	ctx.JSON(http.StatusOK, fmt.Sprintf("success!,%d directory image", limit_page))
}
