package Engine

import (
	"ReEngine/Model"
	utils "ReEngine/util"
	"fmt"
	"log"
	"os"
	"testing"
)

func BenchmarkUpdateIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0;i < b.N ; i++ {
		articles := Model.GetAndReadFiles(utils.DocPath)
		_,err := InsertDoc(articles)
		if err != nil {
			b.Error(err)
		}
		err = CreateInvert(articles)
		if err != nil {
			b.Error(err)
		}
	}
}


func TestInsertDoc(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.DictionaryPath)
	Model.InitDB()
	utils.InitSegment(path)
	path = fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.StopWordPath)
	err := utils.LoadStopWord(path)
	if err != nil {
		log.Fatal(err)
	}
	articles := Model.GetAndReadFiles(utils.DocPath)
	deleteMap,err := InsertDoc(articles)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("deleteMap :",deleteMap)
	for _,it := range articles {
		log.Println("article.id=",it.Id)
	}
	err = CreateInvert(articles)
	if err != nil {
		t.Fatal(err)
	}
}
