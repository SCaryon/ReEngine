package utils

import (
	"fmt"
	"github.com/huichen/sego"
	"os"
)

var Segmenter sego.Segmenter
func InitSegment() {
	environ := os.Environ()
	for i := range environ {
		fmt.Println(environ[i])
	}
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/github.com/huichen/sego/data/dictionary.txt",goPath)
	// 载入词典
	Segmenter.LoadDictionary(path)

}
