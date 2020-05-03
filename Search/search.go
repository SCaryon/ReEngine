package Search

import (
	"ReEngine/Model"
	"ReEngine/util"
	"log"
)
// 处理过程：分词，将结果分别在数据库中寻找对应的doc_id，然后求交
func GetInvert(content string) ([]int,[]string,map[string]int,error){
	// 分词
	seg := utils.SegmentContent(content)
	log.Printf("SegmentContent Result:%v",seg)

	var docId = make(map[int]bool)
	var invert = make(map[string]int)
	for _,word := range seg {
		res,isExist,err := Model.SearchInvertDB(word)
		if err != nil {
			continue
		}
		if isExist == false {
			continue
		}
		idSlice := utils.StringToSlice(res.NumDocs)
		for _,id := range idSlice {
			docId[id] = true
		}
		invert[word] = len(idSlice)

	}
	var result []int
	for key,_ := range docId {
		result = append(result, key)
	}
	return result, seg, invert,nil
}
