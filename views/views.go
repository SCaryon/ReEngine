package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
			rsp := "nil"
			c.HTML(http.StatusOK,"search.html",gin.H{
				"title"		: context,
				"context"	: context,
				"result"	: rsp,
			})
		}
	})
}
