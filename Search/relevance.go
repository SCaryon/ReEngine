package Search

import (
	"ReEngine/Model"
	"ReEngine/util"
	"fmt"
	"log"
	"math"
	"sort"
)

// 给拿到的docId进行相关性排序
func RelevanceSort(docId []int,segs []string,invert map[string]int) []Model.Relevance {
	var resp Model.DocSlice
	docs,err := Model.GetDocByIds(docId)
	if err != nil {
		log.Fatal(err)
	}
	for _,doc := range docs {
		if doc.Content == "" && doc.Title == "" {
			continue
		}
		tmp := Model.Relevance{}
		tmp.Article = doc
		// 优先从redis里面拿到分词
		titleKey := fmt.Sprintf(utils.RedisDoctitleSeg,doc.Id)
		res,err :=  Model.RedisGet(titleKey)
		if err != nil {
			tmpSeg := utils.SegmentContent(tmp.Title)
			_ = Model.RedisSet(titleKey, tmpSeg)
			tmp.TitleSegs = tmpSeg
		} else {
			log.Printf("Get doc Title Seg use Redis,key:%s",titleKey)
			tmp.TitleSegs = res
		}
		contentKey := fmt.Sprintf(utils.RedisDocContentSeg,doc.Id)
		res,err =  Model.RedisGet(contentKey)
		if err != nil {
			tmpSeg := utils.SegmentContent(tmp.Content)
			_ = Model.RedisSet(contentKey, tmpSeg)
			tmp.ContentSegs = tmpSeg
		} else {
			log.Printf("Get doc Content Seg use Redis,key:%s",contentKey)
			tmp.ContentSegs = res
		}
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
			weight += numTF*numIDF
		}
		resp[index].Weight = weight
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