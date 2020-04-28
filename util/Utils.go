package utils

import (
	"bufio"
	"fmt"
	json "github.com/json-iterator/go"
	"log"
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
	DocPageLimit = 10								// 一页的文档数量限制
	ToKenKey = "4kk1HgVV3koDM1L0"				// ToKen
	CookieKey = "ReEngine_token"				// cookie key
	IsLogin = "ReEngine_login"					// 判断是否登陆
	UpdateIndexSpec = "0 0 3 * * ?"				// 更新索引的定时任务参数，每天3天开始更新
)

func StringToSlice(context string) []int {
	var resp = make([]int,0)
	err := json.Unmarshal([]byte(context),&resp)
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

func GetFileCreateTime(file os.FileInfo) int {
	statT := file.Sys().(*syscall.Stat_t)
	tCreate := statT.Ctimespec.Sec
	return int(tCreate)
}


//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteFile(filePath string,content []string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for _,it := range content {
		_, _ = write.WriteString(it + "\n")
	}
	//Flush将缓存的文件真正写入到文件中
	_ = write.Flush()
}

func GetPath() string {
	goPath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/my_go/ReEngine",goPath)
}

func Min(a,b int) int {
	if a < b {
		return  a
	} else {
		return b
	}
}
func Max(a,b int) int {
	if a < b {
		return  b
	} else {
		return a
	}
}