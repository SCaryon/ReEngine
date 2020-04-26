package utils

import (
	"fmt"
	"github.com/huichen/sego"
	"io/ioutil"
	"os"
	"strings"
)

var Segmenter sego.Segmenter
var StopWord []string
func InitSegment() {
	environ := os.Environ()
	for i := range environ {
		fmt.Println(environ[i])
	}
	// 载入词典
	Segmenter.LoadDictionary(DictionaryPath)

}
// 分词后去除停用词
func SegmentContent(content string) []string {
	tmpSegments := Segmenter.Segment([]byte(content))
	seg := sego.SegmentsToSlice(tmpSegments, true)
	dict := make(map[string] bool)

	for _,it := range StopWord {
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

func LoadStopWord() {
	file,err := os.Open(StopWordPath)
	if err != nil {
		return
	}
	tmp, _ := ioutil.ReadAll(file)
	content := string(tmp)
	StopWord = strings.Split(content,"\n")
}
