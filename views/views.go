package views

import (
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/MiddleWare"
	utils "my_go/ReEngine/util"
	"my_go/ReEngine/views/search"
	"my_go/ReEngine/views/user"
	"net/http"
)

func InitRoutes(r *gin.Engine) {
	// 计时
	r.Use(MiddleWare.Timer)
	// 用户权限判断
	r.Use(MiddleWare.AuthMiddleWare)
	r.GET("/test",func(c *gin.Context) {
		key := c.Keys[utils.IsLogin]
		c.HTML(http.StatusOK, "test.html",gin.H{
			"login"		: key,
		})
	})

	// homepage
	r.Any("/",func(c *gin.Context) {
		key := c.Keys[utils.IsLogin]
		log.Printf("log status:%v",key)
		c.HTML(http.StatusOK, "index.html",gin.H{
			"title"		: "HomePage",
			"login"		: key,
		})
	})

	// 管理员相关的操作
	v1:=r.Group("/admin")
	{
		v1.GET("/",func(c *gin.Context){
			user.Manage(r,c)
		})
		v1.Any("/login", func(c *gin.Context) {
			user.LogIn(r,c)
		})
		v1.GET("/logout", func(c *gin.Context) {
			user.LogOut(r,c)
		})
		v1.Any("/register", func(c *gin.Context) {
			user.Register(r,c)
		})
		v1.POST("/submit", func(c *gin.Context) {
			user.SubmitDoc(r,c)
		})
		v1.POST("/delete", func(c *gin.Context) {
			user.DeleteDoc(r,c)
		})
		v1.GET("/update_index",func(c *gin.Context) {
			user.UpdateIndex(r,c)
		})
	}

	// 搜索相关操作
	v2 := r.Group("/s")
	{
		v2.GET("/", func(c *gin.Context) {
			search.SearchContent(r,c)
		})
		v2.GET("/doc", func(c *gin.Context) {
			search.DocmentInfo(r,c)
		})
	}
}

