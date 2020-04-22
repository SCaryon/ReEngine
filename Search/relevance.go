package Search

import (
	"fmt"
	"log"
	"math"
	"my_go/ReEngine/util"
	"sort"
)

// 给拿到的docId进行相关性排序
func RelevanceSort(docId []int,segs []string,invert map[string]int) []utils.Relevance {
	var resp []utils.Relevance
	db := utils.DB
	for _,it := range docId {
		// 拿到文档数据
		queryStr := fmt.Sprintf("select * from %s where id=%d",utils.DBDocment,it)
		rows,err := db.Query(queryStr)
		if err != nil || rows == nil {
			log.Printf("use %s table ,query = %s failed\n",utils.DBDocment,queryStr)
			continue
		}
		for rows.Next() {
			//定义变量接收查询数据
			var tmp utils.Relevance
			err := rows.Scan(&tmp)
			if err != nil {
				log.Printf("get data failed, error:[%v]\n", err.Error())
			}
			// 存储分词结果
			// todo 从redis里面拿到文章的分词的信息
			tmp.TitleSegs = utils.SegmentContent(tmp.Title)
			tmp.ContentSegs = utils.SegmentContent(tmp.Context)
			resp = append(resp,tmp)
		}
	}

	queryStr := fmt.Sprintf("select count(id) from %s",utils.DBDocment)
	rows,err := db.Query(queryStr)
	if err != nil || rows == nil {
		log.Printf("count lines failed ,query:%s",queryStr)
		return nil
	}
	var dataNum int
	for rows.Next() {
		err := rows.Scan(&dataNum)
		if err != nil {
			log.Printf("get data failed, error:[%v]\n", err.Error())
			return nil
		}
	}


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
	sort.Sort(utils.DocSlice(resp))
	return resp
}

func getTF(seg string,doc utils.Relevance) float64 {
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