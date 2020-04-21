package utils

import (
	"fmt"
	"github.com/huichen/sego"
	"io/ioutil"
	"os"
	"strings"
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
// 分词后去除停用词
func SegmentContent(content string) []string {
	tmpSegments := Segmenter.Segment([]byte(content))
	seg := sego.SegmentsToSlice(tmpSegments, true)
	stop, err := GetStopWord()
	dict := make(map[string] bool)
	if err != nil {
		return seg

	}
	for _,it := range stop {
		dict[it] = true
	}
	var res []string
	for _,it := range seg {
		if dict[it] != true {
			res = append(res,it)
		}
	}
	return res
}

func GetStopWord() ([]string, error) {
	file,err := os.Open(StopWord)
	if err != nil {
		return nil,err
	}
	tmp, _ := ioutil.ReadAll(file)
	content := string(tmp)
	res := strings.Split(content,"\n")
	return res, nil
}