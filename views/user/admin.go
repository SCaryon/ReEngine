package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"my_go/ReEngine/Engine"
	"my_go/ReEngine/Model"
	utils "my_go/ReEngine/util"
	"net/http"
)

func Manage(r *gin.Engine, c *gin.Context)  {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not manage web site")
		toHomePage(r,c)
		return
	}
	var res []Model.Article
	res = append(res,Model.Article{Title:"test1",Auth:"SCaryon",Content:"testContext1"})
	res = append(res,Model.Article{Title:"test2",Auth:"SCaryon",Content:"testContext2"})
	docJson, _ := json.Marshal(res)
	log.Println("docjson : ",string(docJson))
	username := Model.GetUsername(c)
	c.HTML(http.StatusOK,"manage.html",gin.H{
		"username"	:	username,
		"login"		: key,
		"numberDoc"	: len(res),
		"docs"		: string(docJson),
		"upload"	: c.Query("upload"),
		"index"		: c.Query("index"),
	})
}

func SubmitDoc(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not submit doc")
		toHomePage(r,c)
		return
	}
	form,err := c.MultipartForm()
	if err != nil {
		log.Printf("上传文件，获取数据出错,%v",err)
		c.Redirect(http.StatusFound,"/admin?upload=0")
		return
	}
	auth := form.Value["auth"]
	files := form.File["upload_file"]
	for _,file:=range files{
		ok:=c.SaveUploadedFile(file,utils.DocPath+file.Filename)
		if ok!=nil{
			fmt.Println("保存的时候出错了 ",ok)
			c.Redirect(http.StatusFound,"/admin?upload=0")
			return
		}
		fmt.Printf("file name :%s,file size :%v, auth : %s",file.Filename,file.Size,auth)
	}
	c.Redirect(http.StatusFound,"/admin?upload=1")
}

func DeleteDoc(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not delete doc")
		toHomePage(r,c)
		return
	}
}

func UpdateIndex(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not update index")
		toHomePage(r,c)
		return
	}
	go Engine.UpdateIndex()
	c.Redirect(http.StatusFound,"/admin?index=1")

}

func LogIn(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == true {
		log.Println("the user is already login,can not login again")
		toHomePage(r,c)
		return
	}
	username := c.DefaultPostForm("username","")
	password := c.DefaultPostForm("password","")
	var warning string
	var check bool
	if username == "" || password == "" {
		check = false
	} else {
		// check password
		check,warning = Model.CheckPassWord(username,password)
	}

	if check {
		log.Println("log success")
		token := Model.CreateToken(username,password)
		Model.SetToken(username,token)
		// set cookie
		c.SetCookie(utils.CookieKey, string(token),0,"/","localhost",false,true)
		// 跳转主页
		toHomePage(r,c)
	} else {
		log.Println("log failed")
		log.Printf("username=%s,password=%s",username,password)
		c.HTML(http.StatusOK,"login.html",gin.H{
			"title"			:	"Login",
			"username"		:	username,
			"password"		:	password,
			"warning"		:	warning,
			"login"			: key,
		})
	}
}

func LogOut(r *gin.Engine,c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not login,can not logout")
		toHomePage(r,c)
		return
	}
	// clear cookie
	token, _ := c.Cookie(utils.CookieKey)
	c.SetCookie(utils.CookieKey,token,-1,"/","localhost",false,true)
	toHomePage(r,c)
}

func Register(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not register user")
		toHomePage(r,c)
		return
	}
	name := c.Query("username")
	password := c.Query("password")
	err := Model.AddUser(name,password)
	if err == nil {

	} else {

	}
}


func toHomePage(r *gin.Engine, c *gin.Context) {
	c.Redirect(http.StatusFound,"/")
}
