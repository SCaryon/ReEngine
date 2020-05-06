package main

import (
	"ReEngine/Engine"
	"ReEngine/Model"
	utils "ReEngine/util"
	"ReEngine/views"
	"github.com/DeanThompson/ginpprof"
	_ "github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
)

func main(){
	// 写log到文件
	//f, _ := os.Create("/log/gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	// 中间件使用
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 加载html模板
	r.LoadHTMLGlob("templates/**/*")
	// 加载静态资源
	r.Static("/static", "./static")
	// 初始化数据库
	Model.InitDB()
	defer Model.DB.Close()
	// 初始化Redis
	Model.InitRedis()
	defer Model.Connect.Close()
	// 初始化分词器
	utils.InitSegment(utils.DictionaryPath)
	// 加载停用词
	utils.LoadStopWord(utils.StopWordPath)
	// 初始化Cache
	utils.InitCache()
	// 路由
	views.InitRoutes(r)

	// 添加定时任务:更新索引表
	Engine.CreateCron()
	defer Engine.CronUpdateIndex.Stop()
	// 使用net/http/pprof
	ginpprof.Wrapper(r)
	r.Run()
}
