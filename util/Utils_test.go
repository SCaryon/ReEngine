package utils

import (
	"my_go/ReEngine/Model"
	"testing"
)

func TestGetAndReadFiles(t *testing.T) {
	articles := Model.GetAndReadFiles(DocPath)
	for _,article := range articles {
		t.Logf("Title:%s,Auth:%s,CreateTime:%d",article.Title,article.Auth,article.CreateTime)
	}

}