package Engine

import (
	"fmt"
	"log"
	utils "my_go/ReEngine/util"
	"os"
	"testing"
)

func TestUpdateIndex(t *testing.T) {

}

func TestInsertDoc(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.DictionaryPath)
	utils.InitDB()
	utils.InitSegment(path)
	path = fmt.Sprintf("%s/src/my_go/ReEngine/%s",goPath,utils.StopWordPath)
	err := utils.LoadStopWord(path)
	if err != nil {
		log.Fatal(err)
	}
	articles := utils.GetAndReadFiles(utils.DocPath)
	deleteMap,err := InsertDoc(articles)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("deleteMap :",deleteMap)
	err = createInvert(articles)
	if err != nil {
		t.Fatal(err)
	}
}
