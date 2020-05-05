package main

import (
	"ReEngine/Engine"
	"ReEngine/Model"
	utils "ReEngine/util"
	"fmt"
	"os"
	"testing"
)

func BenchmarkUpdateIndex(b *testing.B) {
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
	for i := 0;i < b.N ; i++ {
		articles := Model.GetAndReadFiles(utils.DocPath)
		_,err := Engine.InsertDoc(articles)
		if err != nil {
			b.Error(err)
		}
		err = Engine.CreateInvert(articles)
		if err != nil {
			b.Error(err)
		}
	}
}
