package views

import (
	"github.com/gin-gonic/gin"
	"my_go/ReEngine/Search"
	utils "my_go/ReEngine/util"
	"net/http"
	"log"
	"strconv"
)

func InitRoutes(r *gin.Engine) {
	// homepage
	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html",gin.H{
			"title"		: "HomePage",
		})
	})
	// search result s?context=xxxx&offset=x
	r.GET("/s",func(c *gin.Context) {
		context := c.Query("context")
		offset := c.Query("offset")
		if context == "" {
			c.Request.URL.Path = "/"
			r.HandleContext(c)
		} else {
			log.Println("search context:%s",context)
			// 查找倒排索引
			docId,seg := Search.Search(context)
			// 相关性排序
			resp := Search.RelevanceSort(docId,seg)
			// todo redis缓存搜索数据，分页用
			// 分页
			offsetTmp,err := strconv.Atoi(offset)
			if err != nil {
				// 非法参数
				c.Request.URL.Path = "/"
				r.HandleContext(c)
			} else {
				index := offsetTmp * utils.DocLimit
				c.HTML(http.StatusOK,"search.html",gin.H{
					"title"		: context,
					"context"	: context,
					"result"	: resp[index:index+utils.DocLimit],
				})
			}
		}
	})
}
