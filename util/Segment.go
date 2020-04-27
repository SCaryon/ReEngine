package utils

import (
	"bufio"
	"github.com/huichen/sego"
	"os"
	"unicode"
)

var Segmenter sego.Segmenter
var DictStopWord map[string] bool
func InitSegment(path string) {
	// 载入词典
	Segmenter.LoadDictionary(path)
}
// 分词后去除停用词
func SegmentContent(content string) []string {
	tmpSegments := Segmenter.Segment([]byte(content))
	seg := sego.SegmentsToSlice(tmpSegments, true)
	var res []string
	for _,it := range seg {
		if DictStopWord[it] == false && isLegal(it){
			res = append(res,it)
		}
	}
	return res
}

func isLegal(content string) bool {
	if content == " " {
		return  false
	}
	for _,it := range content {
		if unicode.IsPrint(it) == false{
			return false
		}
	}
	return true
}
func LoadStopWord(path string) error {
	file,err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	var stopWord []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopWord = append(stopWord,scanner.Text())
	}
	dict := make(map[string] bool)
	for _,word := range stopWord {
		dict[word] = true
	}
	tmpPath := GetPath()
	WriteFile(tmpPath+"/tmp/tmp.txt",stopWord)
	DictStopWord = dict
	return nil
}
