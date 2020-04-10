package main

import (
	"github.com/gin-gonic/gin"
	"my_go/ReEngine/views"
)

func main(){
	r := gin.Default()
	// use midleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")
	// routes
	views.InitRoutes(r)


	r.Run()
}