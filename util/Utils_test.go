package utils

import (
	"testing"
)

func TestGetAndReadFiles(t *testing.T) {
	articles := GetAndReadFiles(DocPath)
	for _,article := range articles {
		t.Logf("Title:%s,Auth:%s,CreateTime:%d",article.Title,article.Auth,article.CreateTime)
	}

}