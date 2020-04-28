package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/Engine"
	"my_go/ReEngine/Model"
	"my_go/ReEngine/util"
	"my_go/ReEngine/views"
)

func main(){
	log.Println("Engine start")
	r := gin.Default()
	// use middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")

	// 初始化分词器
	utils.InitSegment(utils.DictionaryPath)
	// 初始化数据库
	Model.InitDB()
	defer Model.DB.Close()
	// 加载停用词
	utils.LoadStopWord(utils.StopWordPath)
	// 初始化Cache
	utils.InitCache()
	// 路由
	views.InitRoutes(r)

	// 添加定时任务:更新索引表
	Engine.CreateCron()
	defer Engine.CronUpdateIndex.Stop()

	r.Run()
}
