package Model

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
)
var Connect redis.Conn
func InitRedis() {
	conn,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		log.Fatal("redis init failed")
		return
	}
	Connect = conn
}

func RedisSet(key string,value []string) error {
	for _,str := range value {
		_, err := Connect.Do("rpush", key, str)
		if err != nil {
			log.Fatal(str," err ",err)
			return err
		}
	}
	return nil
}

func RedisGet(key string) ([]string,error) {
	value, err := redis.Strings(Connect.Do("lrange", key, "0", "-1"))
	if len(value) == 0 {
		return nil,errors.New("redis is empty")
	}
	return value,err
}
