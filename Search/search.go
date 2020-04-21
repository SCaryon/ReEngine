package Search

import (
	"fmt"
	"log"
	"my_go/ReEngine/util"
)
// 处理过程：分词，将结果分别在数据库中寻找对应的doc_id，然后求交
func Search(content string) []int{
	// 分词
	seg := utils.SegmentContent(content)
	log.Printf("%v",seg)

	// 连接数据库
	db := utils.DB
	if db == nil {
		log.Println("connect db failed")
		return nil
	}
	var docId = make(map[int]bool)
	for _,tmp := range seg {
		queryStr := fmt.Sprintf("select id,doc_id from %s where key_word=%s\n",utils.DBInvertDoc,tmp)
		rows,err := db.Query(queryStr)
		if err == nil || rows == nil {
			log.Printf("use %s table ,query = %s failed\n",utils.DBInvertDoc,tmp)
			continue
		}
		for rows.Next() {
			//定义变量接收查询数据
			var id int
			var tmpId string
			err := rows.Scan(&id, &tmpId)
			if err != nil {
				log.Printf("get data failed, error:[%v]\n", err.Error())
			}
			log.Println(id, tmpId)
			idSlice := utils.StringToSlice(tmpId)
			for _,id := range idSlice {
				docId[id] = true
			}
		}
	}
	var result []int
	for key,_ := range docId {
		result = append(result, key)
	}
	return result
}
