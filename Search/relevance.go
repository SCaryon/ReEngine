package Search

import (
	"log"
	"math"
	"my_go/ReEngine/Model"
	"my_go/ReEngine/util"
	"sort"
)

// 给拿到的docId进行相关性排序
func RelevanceSort(docId []int,segs []string,invert map[string]int) []Model.Relevance {
	var resp []Model.Relevance
	docs,err := Model.GetDocByIds(docId)
	if err != nil {
		log.Fatal(err)
	}
	for _,doc := range docs {
		tmp := Model.Relevance{}
		tmp.Article = &doc
		tmp.TitleSegs = utils.SegmentContent(tmp.Title)
		tmp.ContentSegs = utils.SegmentContent(tmp.Content)
		resp = append(resp,tmp)
	}
	dataNum, _ := Model.CountDocs()
	// 计算TF-IDF
	var weight float64
	for _,it := range resp {
		weight = 0
		_ = it
		for _,seg := range segs {
			weight += getTF(seg,it)*getIDF(seg, invert, dataNum)
		}
		it.Weight = weight
	}
	// 按权重大小对于结果进行排序
	sort.Sort(Model.DocSlice(resp))
	return resp
}

func getTF(seg string,doc Model.Relevance) float64 {
	var res float64
	var times,all int
	all = len(doc.TitleSegs) + len(doc.ContentSegs)
	times = 0
	for _,it := range doc.TitleSegs {
		if it == seg {
			times += 1
		}
	}
	for _,it := range doc.ContentSegs {
		if it == seg {
			times += 1
		}
	}
	res = float64(times / all)
	return res
}

func getIDF(seg string, invert map[string]int, dataNum int) float64 {
	var res float64
	res = math.Log(float64(dataNum / (invert[seg] + 1)))
	return res
}