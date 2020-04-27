package Engine

import (
	"github.com/robfig/cron"
	"log"
	utils "my_go/ReEngine/util"
	"time"
)

var CronUpdateIndex *cron.Cron

func newCronWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func CreateCron() {
	log.Println("Init Cron",time.Now())
	CronUpdateIndex = newCronWithSecond()
	_, err := CronUpdateIndex.AddFunc(utils.UpdateIndexSpec, AutoUpdateIndex)
	if err != nil {
		log.Println("AddFunc Failed")
	}
	CronUpdateIndex.Start()
}

func AutoUpdateIndex() {
	log.Printf("update doc and invert,time:%v",time.Now())
	UpdateIndex()
}

