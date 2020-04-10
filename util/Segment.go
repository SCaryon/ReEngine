package utils

import (
	"github.com/huichen/sego"
)

var Segmenter sego.Segmenter
func InitSegment() {
	// 载入词典
	Segmenter.LoadDictionary("github.com/huichen/sego/data/dictionary.txt")

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	//sego.SegmentsToString(segments, false)
}
