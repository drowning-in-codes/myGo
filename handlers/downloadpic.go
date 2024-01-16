package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

const pixivision_site = "https://www.pixivision.net/zh/c/illustration"

// type Pics struct {
// 	title string
// 	pics map[string]string
// }

type URL_path interface {
	isPic() bool
	picName() string
}
type url_path string

func (url url_path) isPic() bool {
	ext := url[strings.LastIndex(string(url), ".")+1:]
	if ext == "jpg" || ext == "png" || ext == "gif" || ext == "jpeg" || ext == "webp" {
		return true
	}
	return false
}

func (url url_path) picName() string {
	img_name := strings.Split(string(url), "/")[len(strings.Split(string(url), "/"))-1]
	return img_name
}

func DownloadPicHandler(ctx *gin.Context) {
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
	c := colly.NewCollector(colly.AllowedDomains(allow_img_site...),
		colly.Async())
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*pximg.*",
		Parallelism: 5,
		//Delay:      5 * time.Second,
		RandomDelay: 500 * time.Duration(time.Millisecond),
	})
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*pixivision.*",
		Parallelism: 5,
		Delay:       200 * time.Duration(time.Millisecond),
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
	counter := 0

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.DOM.Parent().HasClass("arc__title") {
			log.Default().Println("Link found:", e.Attr("href"))
			e.Request.Visit(e.Attr("href"))
		}
	})
	c.OnHTML("div[class='_article-main']", func(e *colly.HTMLElement) {

		title := e.ChildText("h1[class='am__title']")
		// log.Default().Println("title:", title)
		// p := Pics{title: title, pics: make(map[string]string)}
		err := os.MkdirAll(filepath.Join(download_root_folder, title), os.ModePerm)
		if err != nil {
			log.Default().Println(err.Error())
		}
		e.ForEach("div.article-item:not(._feature-article-body__paragraph) div.am__work__main", func(i int, h *colly.HTMLElement) {
			log.Default().Println("pic:", h.ChildAttr("img", "src"))

			img_src := h.ChildAttr("img", "src")
			h.Request.Visit(img_src)
			h.Request.Ctx.Put("title", title)
		})
		counter += 1
		if counter >= limit_page {
			return
		}

	})
	c.OnResponse(func(r *colly.Response) {
		var img_url URL_path = url_path(r.Request.URL.Path)
		// log.Default().Println("img_name:", img_name)
		if img_url.isPic() {
			img_path := filepath.Join(r.Ctx.Get("title"), string(img_url.(url_path)))
			log.Default().Println("save path", img_path)
			r.Save(img_path)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)

	})
	c.OnError(func(r *colly.Response, err error) {
		log.Default().Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(pixivision_site)
	c.Wait()
	ctx.JSON(http.StatusOK, fmt.Sprintf("success!,%d directory image", limit_page))
}
