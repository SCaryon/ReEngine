package utils

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)
const (
	ValueTimeLimit	= 48 * time.Hour
)
var BigCache *bigcache.BigCache
func InitCache() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(ValueTimeLimit))
	if err != nil {
		log.Fatalf("cache init error :%s",err)
		return
	}
	BigCache = cache
}