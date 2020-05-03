package Engine

import (
	"ReEngine/Model"
	utils "ReEngine/util"
	"fmt"
	"log"
	"os"
	"testing"
)


func TestInsertDoc(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/ReEngine/%s",goPath,utils.DictionaryPath)
	Model.InitDB()
	utils.InitSegment(path)
	path = fmt.Sprintf("%s/src/ReEngine/%s",goPath,utils.StopWordPath)
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
	err = createInvert(articles)
	if err != nil {
		t.Fatal(err)
	}
}
