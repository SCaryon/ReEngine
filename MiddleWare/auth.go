package MiddleWare

import (
	"ReEngine/Model"
	utils "ReEngine/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(c *gin.Context) {
	token, _ := c.Cookie(utils.CookieKey)
	if Model.IsToKenLegal([]byte(token)) {
		c.Set(utils.IsLogin,true)
	} else {
		c.Set(utils.IsLogin,false)
	}
	c.Next()
}
