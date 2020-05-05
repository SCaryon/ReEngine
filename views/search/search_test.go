package search

import (
	"ReEngine/Model"
	"ReEngine/Search"
	utils "ReEngine/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func BenchmarkSearchContent(b *testing.B) {
	content := "自动机"
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.DictionaryPath)
	Model.InitDB()
	utils.InitSegment(path)
	path = fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.StopWordPath)
	err := utils.LoadStopWord(path)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i <= b.N;i++ {
		mokeSearch(content)
	}
}

func mokeSearch(content string) {
	// 查找倒排索引
	docId,seg,invert,err := Search.GetInvert(content)
	if err != nil {
		log.Println("search invert failed",err)
		return
	}
	// 相关性排序
	docs := Search.RelevanceSort(docId,seg,invert)
	_,err = json.Marshal(docs)
	if err != nil {
		log.Fatal(err)
	}
}