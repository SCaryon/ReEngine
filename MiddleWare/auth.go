package MiddleWare

import (
	"github.com/gin-gonic/gin"
	utils "my_go/ReEngine/util"
)

func AuthMiddleWare(c *gin.Context) {
	token, _ := c.Cookie(utils.CookieKey)
	if utils.IsToKenLegal([]byte(token)) {
		c.Set(utils.IsLogin,true)
	} else {
		c.Set(utils.IsLogin,false)
	}
	c.Next()
}
