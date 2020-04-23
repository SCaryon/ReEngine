package views

import (
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/Search"
	utils "my_go/ReEngine/util"
	"net/http"
	"strconv"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/test",func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html",gin.H{})
	})
	// 管理员相关的操作
	admin:=r.Group("/admin")
	{
		admin.GET("/login", func(c *gin.Context) {
			userLogIn(r,c)
		})
		admin.GET("/register", func(c *gin.Context) {
			userRegister(r,c)
		})
		admin.GET("/submit", func(c *gin.Context) {
			submitDoc(r,c)
		})
		admin.GET("/delete", func(c *gin.Context) {
			deleteDoc(r,c)
		})
	}
	// homepage
	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html",gin.H{
			"title"		: "HomePage",
		})
	})
	// 搜索相关操作
	search := r.Group("/s")
	{
		search.GET("/", func(c *gin.Context) {
			searchContent(r,c)
		})
	}
}

func userLogIn(r *gin.Engine, c *gin.Context) {
	c.HTML(http.StatusOK,"login.html",gin.H{
		"title"			:	"Login",
		"warning"		:	"no data",
	})
}
func userRegister(r *gin.Engine, c *gin.Context) {

}
func submitDoc(r *gin.Engine, c *gin.Context) {

}
func deleteDoc(r *gin.Engine, c *gin.Context) {

}
func searchContent(r *gin.Engine, c *gin.Context) {
	// search result s?context=xxxx&offset=x
	content := c.Query("content")
	offset := c.Query("offset")
	if content == "" {
		toHomePage(r,c)
		return
	}
	log.Println("search content:%s",content)
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
	})
}

func toHomePage(r *gin.Engine, c *gin.Context) {
	c.Request.URL.Path = "/"
	r.HandleContext(c)
}