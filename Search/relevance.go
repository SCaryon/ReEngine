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
	var resp Model.DocSlice
	docs,err := Model.GetDocByIds(docId)
	if err != nil {
		log.Fatal(err)
	}
	for index := range docs {
		log.Printf("index=%d,doc id=%d,title=%s",index,docs[index].Id,docs[index].Title)
	}
	for _,doc := range docs {
		tmp := Model.Relevance{}
		tmp.Article = doc

		tmp.TitleSegs = utils.SegmentContent(tmp.Title)
		tmp.ContentSegs = utils.SegmentContent(tmp.Content)
		resp = append(resp,tmp)
	}
	dataNum, _ := Model.CountDocs()
	// 计算TF-IDF
	var weight float64
	for index,doc := range resp {
		weight = 0
		for _,seg := range segs {
			numTF := getTF(seg,doc)
			numIDF := getIDF(seg, invert, dataNum)
			log.Printf("doc=%d,numTF=%f,numIDF=%f",doc.Id,numTF,numIDF)
			weight += numTF*numIDF
		}
		resp[index].Weight = weight
		log.Printf("doc id=%d,weight=%f",doc.Id,weight)
	}
	// 按权重大小对于结果进行排序
	sort.Sort(resp)
	for index,doc := range resp {
		resp[index].Content = utils.CutString(doc.Content,0,100)
	}
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
	res = utils.IntToFloat64(times) / utils.IntToFloat64(all)
	return res
}

func getIDF(seg string, invert map[string]int, dataNum int) float64 {
	var res float64
	res = math.Log(utils.IntToFloat64(dataNum) / utils.IntToFloat64(invert[seg] + 1))
	return res
}