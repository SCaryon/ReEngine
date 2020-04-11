package views

import (
	"github.com/gin-gonic/gin"
	"my_go/ReEngine/Search"
	"net/http"
	"log"
)

func InitRoutes(r *gin.Engine) {
	// homepage
	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html",gin.H{
			"title"		: "HomePage",
		})
	})
	// search result s?context=xxxx
	r.GET("/s",func(c *gin.Context) {
		context := c.Query("context")
		if context == "" {
			c.Request.URL.Path = "/"
			r.HandleContext(c)
		} else {
			log.Println("search context:%s",context)
			// 查找倒排索引
			docId := Search.Search(context)
			// 相关性排序
			resp := Search.RelevanceSort(docId)
			c.HTML(http.StatusOK,"search.html",gin.H{
				"title"		: context,
				"context"	: context,
				"result"	: resp,
			})
		}
	})
}
