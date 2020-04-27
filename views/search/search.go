package search

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
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
	log.Println("offset",offset)
	key := c.Keys[utils.IsLogin]
	if content == "" {
		toHomePage(r,c)
		return
	}
	log.Printf("search content:%s",content)
	// 查找倒排索引
	docId,seg,invert,err := Search.SearchInvert(content)
	if err != nil {
		log.Println("search invert failed",err)
		toHomePage(r,c)
		return
	}
	log.Printf("SearchInvert result:%v",docId)
	// 相关性排序
	docs := Search.RelevanceSort(docId,seg,invert)
	log.Printf("RelevanceSort result:%v",docs)
	// todo redis缓存搜索数据，分页用
	// 分页
	offsetTmp,err := strconv.Atoi(offset)
	if err != nil {
		// 非法参数
		log.Println("get offset failed",err)
		offsetTmp = 0
	}
	downIndex := offsetTmp * utils.DocPageLimit
	upIndex := utils.Min(len(docs), downIndex+utils.DocPageLimit)
	docJson, _ := json.Marshal(docs[downIndex:upIndex])
	log.Println(downIndex,upIndex)
	c.HTML(http.StatusOK,"search.html",gin.H{
		"title"		: content,
		"content"	: content,
		"docs"		: string(docJson),
		"numberDoc"	: len(docs),
		"login"		: key,
	})
}

func toHomePage(r *gin.Engine, c *gin.Context) {
	c.Redirect(http.StatusFound,"/")
}
