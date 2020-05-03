package Model

import (
	"fmt"
	"log"
	utils "my_go/ReEngine/util"
)

type Invert struct {
	id			int
	KeyWord		string
	NumDocs		string
}

func SearchInvertDB(word string) (Invert,bool,error) {
	resp := Invert{}
	queryStr := fmt.Sprintf("select id,doc_id from %s where key_word=\"%s\"",utils.DBInvertDoc, word)
	rows,err := DB.Query(queryStr)
	if err != nil || rows == nil {
		log.Printf("select id,doc_id failed,err=%v",err)
		return resp,false,err
	}
	isExist := false
	for rows.Next() {
		isExist = true
		var id int
		var tmpId string
		err := rows.Scan(&id, &tmpId)
		if err != nil {
			log.Printf("get data failed, error:[%v]\n", err.Error())
		}
		resp.id = id
		resp.NumDocs = tmpId
		resp.KeyWord = word
	}
	return resp,isExist,nil
}

func UpdateInvert(invert Invert) error{
	updateStr := fmt.Sprintf("UPDATE %s SET doc_id=? where id=?",utils.DBInvertDoc)
	_, err := DB.Exec(updateStr, invert.NumDocs, invert.id)
	return err
}

func InsertInvert(word string,id int) error {
	queryStr := fmt.Sprintf("INSERT INTO %s(key_word,doc_id)VALUES (?,?)", utils.DBInvertDoc)
	_, err := DB.Exec(queryStr, word, utils.SliceToString([]int{id}))
	if err != nil {
		log.Printf("insert invert_index failed,err=%v,queryStr=%s", err, queryStr)
		return err
	}
	return nil
}