package utils

import (
	json "github.com/json-iterator/go"
)

const (
	DBDocment = "docment"	// 存文档数据的
	DBInvertDoc = "invert_index"	// 倒排索引
)



func StringToSlice(context string) []int {
	var resp []int
	err := json.Unmarshal([]byte(context),resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func SliceToString(context []int) string {
	resp,err := json.Marshal(context)
	if err != nil {
		panic(err)
	}
	return string(resp)
}