package Model

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
)
var pool *redis.Pool
func InitRedis() {
	pool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		Wait: true,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func RedisSet(key string,value []string) error {
	connect := pool.Get()
	defer connect.Close()
	for _,str := range value {
		_, err := connect.Do("rpush", key, str)
		if err != nil {
			log.Fatal(str," err ",err)
			return err
		}
	}
	return nil
}

func RedisGet(key string) ([]string,error) {
	connect := pool.Get()
	defer connect.Close()
	value, err := redis.Strings(connect.Do("lrange", key, "0", "-1"))
	if len(value) == 0 {
		return nil,errors.New("redis is empty")
	}
	return value,err
}
