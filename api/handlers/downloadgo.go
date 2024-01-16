package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func LoadGoHander(c *gin.Context) {

	allow_img_site := checkAllowSite()
	c := colly.NewCollector(colly.UserAgent(userAgent), colly.AllowedDomains(allow_img_site...),colly.Async())
    




    
    c.Visit(mygo_site)

}
