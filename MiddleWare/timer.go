package MiddleWare

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Timer(c *gin.Context) {
	t := time.Now()
	//在gin上下文中定义一个变量
	c.Set("timer", "CustomRouterMiddle1")
	//请求之前
	c.Next()
	//请求之后
	//计算整个请求过程耗时
	t2 := time.Since(t)
	log.Println("the request cost time is ",t2)
}