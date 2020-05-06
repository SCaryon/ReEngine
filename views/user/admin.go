package user

import (
	"ReEngine/Engine"
	"ReEngine/Model"
	utils "ReEngine/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Manage(r *gin.Engine, c *gin.Context)  {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not manage web site")
		toHomePage(c)
		return
	}
	offset := c.Query("offset")
	docs, err := Model.GetAllDocs()
	for index,doc := range docs {
		docs[index].Content = utils.CutString(doc.Content,0,100)
	}
	if err != nil {
		log.Println(err)
		toHomePage(c)
	} else {
		// 分页
		offsetTmp,err := strconv.Atoi(offset)
		if err != nil {
			// 非法参数
			log.Println("get offset failed",err)
			offsetTmp = 0
		}
		downIndex := offsetTmp * utils.DocPageLimit
		if downIndex >= len(docs) {
			downIndex = 0
		}
		upIndex := utils.Min(len(docs), downIndex+utils.DocPageLimit)
		docJson, _ := json.Marshal(docs[downIndex:upIndex])
		username := Model.GetUsername(c)
		c.HTML(http.StatusOK,"manage.html",gin.H{
			"username"	:	username,
			"login"		: key,
			"numberDoc"	: len(docs),
			"docs"		: string(docJson),
			"upload"	: c.Query("upload"),
			"index"		: c.Query("index"),
			"regSucc"	: c.Query("regSucc"),
			"pageNum"	: offsetTmp+1,
		})
	}
}

func SubmitDoc(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not submit doc")
		toHomePage(c)
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
		fileName := file.Filename
		fileTmp := fileName
		for _,suffix := range utils.FileSuffixes {
			if strings.HasSuffix(fileName,suffix) {
				fileTmp = strings.TrimRight(fileName,suffix)
				break
			}
		}
		user := Model.GetUsername(c)
		ok:=c.SaveUploadedFile(file,utils.DocPath+fileTmp+"."+user)
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
		toHomePage(c)
		return
	}
	id := c.Query("doc_id")
	if id == "" {
		toHomePage(c)
		return
	}
	docId,_ := strconv.Atoi(id)
	err := Model.DeleteDoc(docId)
	if err != nil {
		log.Printf("Delete doc failed,%v",err)
		webStr := fmt.Sprintf("/s/doc?id=%d",docId)
		c.Redirect(http.StatusFound,webStr)
	} else {
		c.Redirect(http.StatusFound,"/admin")
	}
}

func UpdateIndex(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not update index")
		toHomePage(c)
		return
	}
	go Engine.UpdateIndex()
	c.Redirect(http.StatusFound,"/admin?index=1")

}

// 对于文档的修改，采用删除旧文档然后添加新文档的方式进行
func EditDocument(r *gin.Engine,c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not edit document")
		toHomePage(c)
		return
	}
	id := c.Query("doc_id")
	if id == "" {
		toHomePage(c)
		return
	}
	docId,_ := strconv.Atoi(id)
	isEditError := c.Query("status")
	title := c.DefaultPostForm("title","")
	auth := c.DefaultPostForm("auth","")
	content := c.DefaultPostForm("content","")
	// 访问修改页面
	if content == "" && title == "" && auth == "" {
		docs,_ := Model.GetDocByIds([]int{docId})
		warn := ""
		if isEditError == "1" {
			warn = "网站打瞌睡了，请稍后再试"
		}
		c.HTML(http.StatusOK,"edit.html",gin.H{
			"title"			: "Edit",
			"login"			: key,
			"doc_id"		: docId,
			"doc_title"		: docs[0].Title,
			"doc_content"	: docs[0].Content,
			"doc_auth"		: docs[0].Auth,
			"warn"			: warn,
		})
		return
	}
	// 创建修改之后的文章对象
	var afterArticle Model.Article
	afterArticle.Title = title
	afterArticle.Auth = auth
	afterArticle.Content = content
	afterArticle.CreateTime = int(time.Now().Unix())
	log.Printf("title:%s,auth:%s,content:%s",title,auth,content)
	// 提交修改请求
	newId,err := Engine.UpdateDoc(docId, afterArticle)
	if err != nil {
		var webStr string
		if newId == -1 {
			webStr = fmt.Sprintf("/admin/doc_edit/?doc_id=%d&&status=1",docId)
		} else {
			webStr = fmt.Sprintf("/admin/doc_edit/?doc_id=%d&&status=1",newId)
		}
		c.Redirect(http.StatusFound,webStr)
	}
	webStr := fmt.Sprintf("/s/doc/?id=%d",newId)
	c.Redirect(http.StatusFound,webStr)

}

func LogIn(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == true {
		log.Println("the user is already login,can not login again")
		toHomePage(c)
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
		toHomePage(c)
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
		toHomePage(c)
		return
	}
	// clear cookie
	token, _ := c.Cookie(utils.CookieKey)
	c.SetCookie(utils.CookieKey,token,-1,"/","localhost",false,true)
	toHomePage(c)
}

func Register(r *gin.Engine, c *gin.Context) {
	key := c.Keys[utils.IsLogin]
	if key == false {
		log.Println("the user is not admin,can not register user")
		toHomePage(c)
		return
	}
	username := c.DefaultPostForm("username","")
	password := c.DefaultPostForm("password","")
	if username == "" || password == "" {
		c.HTML(http.StatusOK,"register.html",gin.H{
			"title"			:	"Register",
			"login"			: key,
		})
		return
	}
	err := Model.AddUser(username,password)
	if err == nil {
		c.Redirect(http.StatusFound,"/admin?regSucc=1")
	} else {
		c.HTML(http.StatusOK,"register.html",gin.H{
			"title"			:	"Register",
			"login"			: key,
			"extra"			: err.Error(),
			"username"		: username,
			"password"		: password,
		})
	}
}


func toHomePage(c *gin.Context) {
	c.Redirect(http.StatusFound,"/")
}
