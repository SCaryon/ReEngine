package search

import (
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/Search"
	utils "my_go/ReEngine/util"
	"net/http"
	"strconv"
)


func Docment(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	c.HTML(http.StatusOK,"docment.html",gin.H{
		"login"		: key,
	})
}

func SearchContent(r *gin.Engine, c *gin.Context) {
	// search result s?context=xxxx&offset=x
	content := c.Query("content")
	offset := c.Query("offset")
	key := c.Keys[utils.IsLogin]
	if content == "" {
		toHomePage(r,c)
		return
	}
	log.Printf("search content:%s",content)
	// 查找倒排索引
	docId,seg,invert,err := Search.Search(content)
	if err != nil {
		toHomePage(r,c)
		return
	}
	// 相关性排序
	resp := Search.RelevanceSort(docId,seg,invert)
	// todo redis缓存搜索数据，分页用
	// 分页
	offsetTmp,err := strconv.Atoi(offset)
	if err != nil {
		// 非法参数
		toHomePage(r,c)
		return
	}
	index := offsetTmp * utils.DocLimit
	c.HTML(http.StatusOK,"search.html",gin.H{
		"title"		: content,
		"content"	: content,
		"result"	: resp[index:index+utils.DocLimit],
		"login"		: key,
	})
}

func toHomePage(r *gin.Engine, c *gin.Context) {
	c.Redirect(http.StatusFound,"/")
}
