package utils

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"time"
)

const (
	DBDocment = "docment"			// 文档数据
	DBInvertDoc = "invert_index"	// 倒排索引
	DocPath = "tmp/docment/"		// 临时文档路径
)

type Article struct {
	id int
	title string
	auth string
	context string
	createTime int
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

func GetAndReadFiles() []Article{
	var resp []Article
	environ := os.Environ()
	for i := range environ {
		fmt.Println(environ[i])
	}
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/%s",goPath,DocPath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	// 获取文件，并输出它们的名字
	for _, file := range files {
		txt, err := ioutil.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}
		// 将字节流转换为字符串
		content := string(txt)
		resp = append(resp,Article{context:content,createTime:int(time.Now().Unix())})
	}
	return resp
}