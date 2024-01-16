package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func DownloadPixvisionPicHandler(ctx *gin.Context) {
	allow_img_site := checkAllowSite()
	// fmt.Println(allow_img_site)
	c := colly.NewCollector(colly.UserAgent(userAgent), colly.AllowedDomains(allow_img_site...),
		colly.Async())
	c.SetRequestTimeout(20 * time.Second)

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
			if counter >= limit_page {
				ctx.JSON(http.StatusOK, fmt.Sprintf("success!,%d directory image", limit_page))
			}
			e.Request.Visit(e.Attr("href"))
			counter += 1
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
		if r.URL.Host == "i.pximg.net" {
			r.Headers.Set("Referer", "https://www.pixivision.net/")
		}
	})
	c.OnError(func(r *colly.Response, err error) {

		log.Default().Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err.Error())
	})
	c.Visit(pixivision_site)
	c.Wait()
	ctx.JSON(http.StatusOK, fmt.Sprintf("success!,%d directory image", limit_page))
}
