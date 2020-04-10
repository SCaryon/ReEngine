package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/util"
	"my_go/ReEngine/views"
)

func main(){
	log.Println("Engine start")
	r := gin.Default()
	// use midleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")

	// 初始化分词器
	utils.InitSegment()
	// routes
	views.InitRoutes(r)

	r.Run()
}

