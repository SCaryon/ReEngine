package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

/*
	Token使用方式
	一个username对应一个token
	在bigcache里面存储着双向关系，即：
	key:username	value:token
	key:token		value:name
*/
func CreateToken(user ,password string) []byte {
	timeNow := time.Now().Unix()
	str := fmt.Sprintf("%d%s%s+%s",timeNow,user,password,ToKenKey)
	enCodeStr := fmt.Sprintf("%x",md5.Sum([]byte(str)))
	return []byte(enCodeStr)
}

func IsToKenLegal(token []byte) bool {
	if token == nil || len(token) == 0{
		return false
	}
	name, _ := BigCache.Get(string(token))
	tmpToken, _ := BigCache.Get(string(name))
	if bytes.Equal(tmpToken,token) {
		return true
	}
	queryStr := fmt.Sprintf("select token from %s where name=%s",DBUsers,name)
	rows, err := DB.Query(queryStr)
	if err != nil {
		return false
	}
	for rows.Next() {
		err := rows.Scan(&tmpToken)
		if err != nil {
			log.Fatalf("get token failed %s",err)
			return false
		}
	}
	if bytes.Equal(tmpToken,token) {
		return true
	}
	return false
}

func SetToken(name string,token []byte) {
	db := DB
	updateStr := fmt.Sprintf("UPDATE %s SET token=? where name=?",DBUsers)
	db.Exec(updateStr,token,name)
	cache := BigCache
	cache.Set(name,token)
	cache.Set(string(token), []byte(name))
}

func GetUsername(c *gin.Context) string {
	token, _ := c.Cookie(CookieKey)
	name, _ := BigCache.Get(token)
	return string(name)
}
