package Search

import (
	"fmt"
	"log"
	"math"
	"my_go/ReEngine/Model"
	"my_go/ReEngine/util"
	"sort"
)

// 给拿到的docId进行相关性排序
func RelevanceSort(docId []int,segs []string,invert map[string]int) []Model.Relevance {
	var resp []Model.Relevance
	db := utils.DB
	for _,it := range docId {
		// 拿到文档数据
		queryStr := fmt.Sprintf("select id,title,auth,context,create_time from %s where id=%d",utils.DBDocment,it)
		rows,err := db.Query(queryStr)
		if err != nil || rows == nil {
			log.Printf("use %s table ,query = %s failed\n",utils.DBDocment,queryStr)
			continue
		}
		log.Println("queryStr",queryStr)
		for rows.Next() {
			//定义变量接收查询数据
			tmpId := 0
			tmpTitle := ""
			tmpAuth := ""
			tmpContent := ""
			tmpTime := 0
			err := rows.Scan(&tmpId,&tmpTitle,&tmpAuth,&tmpContent,&tmpTime)
			if err != nil {
				log.Printf("get data failed, error:[%v]\n", err.Error())
			}
			tmpArticle := Model.Article{Id:tmpId,Title:tmpTitle,Content:tmpContent,CreateTime:tmpTime}
			tmp := Model.Relevance{Article: &tmpArticle}
			log.Printf("Article info:%v",tmp)
			// 存储分词结果
			// todo 从redis里面拿到文章的分词的信息
			tmp.TitleSegs = utils.SegmentContent(tmp.Title)
			tmp.ContentSegs = utils.SegmentContent(tmp.Content)
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