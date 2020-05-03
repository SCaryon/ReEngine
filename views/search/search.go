package search

import (
	"ReEngine/Model"
	"ReEngine/Search"
	utils "ReEngine/util"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"log"
	"net/http"
	"strconv"
)


func DocumentInfo(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	id := c.Query("id")
	docId,_ := strconv.Atoi(id)
	doc,err := Model.GetDocByIds([]int{docId})
	docJson, _ := json.Marshal(doc[0])
	if err != nil {
		toHomePage(c)
		return
	}
	c.HTML(http.StatusOK,"document.html",gin.H{
		"title"		: "document",
		"login"		: key,
		"doc_id"	: docId,
		"doc"		: string(docJson),
	})
}

func SearchContent(r *gin.Engine, c *gin.Context) {
	// search result s?context=xxxx&offset=x
	content := c.Query("content")
	offset := c.Query("offset")
	log.Println("offset",offset)
	key := c.Keys[utils.IsLogin]
	if content == "" {
		toHomePage(c)
		return
	}
	var docs []Model.Relevance
	log.Printf("search content:%s",content)
	// todo redis缓存搜索数据，分页用 还未测试
	tmpRes,err := utils.BigCache.Get(content)
	if  err == nil && tmpRes != nil {
		_ =json.Unmarshal(tmpRes,&docs)
		log.Printf("search %s,use bigcache",content)

	} else {
		// 查找倒排索引
		docId,seg,invert,err := Search.GetInvert(content)
		if err != nil {
			log.Println("search invert failed",err)
			toHomePage(c)
			return
		}
		log.Printf("GetInvert result:%v",docId)
		// 相关性排序
		docs = Search.RelevanceSort(docId,seg,invert)
		docsJson,err := json.Marshal(docs)
		if err != nil {
			log.Fatal(err)
		}
		err = utils.BigCache.Set(content, docsJson)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 分页
	offsetTmp,err := strconv.Atoi(offset)
	if err != nil {
		// 非法参数
		log.Println("get offset failed",err)
		offsetTmp = 0
	}
	downIndex := offsetTmp * utils.DocPageLimit
	if downIndex >= len(docs) {
		downIndex = 0
	}
	upIndex := utils.Min(len(docs), downIndex+utils.DocPageLimit)
	docJson, _ := json.Marshal(docs[downIndex:upIndex])
	log.Println(downIndex,upIndex)
	c.HTML(http.StatusOK,"search.html",gin.H{
		"title"		: content,
		"content"	: content,
		"docs"		: string(docJson),
		"numberDoc"	: len(docs),
		"login"		: key,
		"pageNum"	: offsetTmp+1,
	})
}

func toHomePage(c *gin.Context) {
	c.Redirect(http.StatusFound,"/")
}
