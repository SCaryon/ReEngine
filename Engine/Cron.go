package Engine

import (
	utils "ReEngine/util"
	"github.com/robfig/cron"
	"log"
	"time"
)

var CronUpdateIndex *cron.Cron

func newCronWithSecond() *cron.Cron {
	//secondParser := cron.NewParser(cron.Second | cron.Minute |
	//	cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New()
	//return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func CreateCron() {
	log.Println("Init Cron",time.Now())
	CronUpdateIndex = newCronWithSecond()
	CronUpdateIndex.AddFunc(utils.UpdateIndexSpec, AutoUpdateIndex)
	CronUpdateIndex.Start()
}

func AutoUpdateIndex() {
	log.Printf("update doc and invert,time:%v",time.Now())
	UpdateIndex()
}

