package utils

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"my_go/ReEngine/Model"
	"os"
	"syscall"
)

const (
	DBDocment = "docment"						// 文档数据
	DBInvertDoc = "invert_index"				// 倒排索引
	DBUsers = "users"							// 用户表
	DocPath = "tmp/docment/"					// 临时文档路径
	StopWordPath = "static/StopWord.txt"		// 停用词路径
	DictionaryPath = "static/Dictionary.txt"	// 停用词路径
	DocLimit = 10								// 一页的文档数量限制
	ToKenKey = "4kk1HgVV3koDM1L0"				// ToKen
	CookieKey = "ReEngine_token"				// cookie key
	IsLogin = "ReEngine_login"					// 判断是否登陆
)

func StringToSlice(context string) []int {
	var resp []int
	err := json.Unmarshal([]byte(context),resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func SliceToString(context []int) string {
	resp,err := json.Marshal(context)
	if err != nil {
		panic(err)
	}
	return string(resp)
}

func GetAndReadFiles(filePath string) []Model.Article{
	var articles []Model.Article
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,filePath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		txt, err := ioutil.ReadFile(path+file.Name())
		if err != nil {
			panic(err)
		}
		// 将字节流转换为字符串
		content := string(txt)
		createTime := GetFileCreateTime(file)
		articles = append(articles,Model.Article{Title: file.Name(),Content:content,CreateTime: createTime})

	}
	fmt.Printf("files info:%v",articles)
	return articles
}

func GetFileCreateTime(file os.FileInfo) int {
	statT := file.Sys().(*syscall.Stat_t)
	tCreate := statT.Ctimespec.Sec
	return int(tCreate)
}

