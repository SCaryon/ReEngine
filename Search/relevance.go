package Search

import (
	"fmt"
	"log"
	"my_go/ReEngine/util"
)

// 给拿到的docId进行相关性排序
func RelevanceSort(docId []int,seg []string) []utils.Relevance {
	var resp []utils.Relevance
	db := utils.DB
	for _,it := range docId {
		// 拿到文档数据
		queryStr := fmt.Sprintf("select * from %s where id=%d",utils.DBDocment,it)
		rows,err := db.Query(queryStr)
		if err == nil || rows == nil {
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
			resp = append(resp,tmp)
		}
	}

	// 对于文档数据计算每个词的频率，方法是TF-IDF
	for _,it := range resp {
		_ = it
	}

	return resp
}
