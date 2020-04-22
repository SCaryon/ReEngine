package utils

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"syscall"
)

const (
	DBDocment = "docment"			// 文档数据
	DBInvertDoc = "invert_index"	// 倒排索引
	DocPath = "tmp/docment/"		// 临时文档路径
	StopWordPath = "util/StopWord.txt"	// 停用词路径
	DocLimit = 10					// 一页的文档数量限制
)

type Article struct {
	Id         int
	Title      string
	Auth       string
	Context    string
	CreateTime int
}

type Invert struct {
	KeyWord		string
	NumDocs		int
}



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

func GetAndReadFiles(filePath string) []Article{
	var resp []Article
	environ := os.Environ()
	for i := range environ {
		fmt.Println(environ[i])
	}
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/%s",goPath,filePath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		txt, err := ioutil.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}
		// 将字节流转换为字符串
		content := string(txt)
		createTime := GetFileCreateTime(file)
		resp = append(resp,Article{Title: file.Name(),Context:content,CreateTime: createTime})
	}
	return resp
}

func GetFileCreateTime(file os.FileInfo) int {
	statT := file.Sys().(*syscall.Stat_t)
	tCreate := statT.Ctimespec.Sec
	return int(tCreate)
}