package Engine

import (
	"fmt"
	"log"
	"my_go/ReEngine/Model"
	utils "my_go/ReEngine/util"
	"os"
)

func AddDocment() {

}

func DelDocment() {

}

// 将存放在tmp文件中的临时文档读取并进行存储和生成倒排索引然后删除文档文件
func UpdateIndex() bool {
	articles := utils.GetAndReadFiles(utils.DocPath)
	db := utils.DB
	if db == nil {
		log.Println("connect db failed")
		return false
	}
	deleteMap,err := InsertDoc(articles)
	if err != nil {
		log.Fatal(err)
		return false
	}
	err = createInvert(articles)
	if err != nil {
		log.Fatal(err)
		return false
	}
	errs := deleteTmpDoc(deleteMap,articles)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return true
}

func InsertDoc(articles []Model.Article) (map[int]bool,error) {
	db := utils.DB
	deleteMap := make(map[int]bool)
	// 将文档放入数据库
	for i,article := range articles {
		queryStr := fmt.Sprintf("INSERT INTO %s(title,auth,context,create_time)VALUES (?,?,?,?)",utils.DBDocment)
		result,err := db.Exec(queryStr,article.Title,article.Auth,article.Content,article.CreateTime)
		if err != nil {
			log.Printf("insert failed,err=%v",err)
			continue
		}
		id,err := result.LastInsertId()
		articles[i].Id = int(id)
		log.Printf("insert file,Title=%s,Auth=%s,Id=%d",article.Title,article.Auth,int(id))
		// 对于未能成功放入数据库的文档暂时不删除
		deleteMap[i] = true
	}
	return deleteMap,nil
}

func createInvert(articles []Model.Article) error{
	// 创建倒排索引
	db := utils.DB
	for _,article := range articles {
		dictWord := make(map[string] bool)
		seg := utils.SegmentContent(fmt.Sprintf("%s %s %s",article.Title,article.Content,article.Auth))
		for _,word := range seg {
			// 去重，使得每一个单词
			if dictWord[word] == true {
				continue
			} else {
				dictWord[word] = true
			}
			queryStr := fmt.Sprintf("select id,doc_id from %s where key_word=\"%s\"",utils.DBInvertDoc, word)
			rows,err := db.Query(queryStr)
			if err != nil || rows == nil {
				log.Printf("select id,doc_id failed,err=%v",err)
				continue
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
				log.Println(id, tmpId)
				idSlice := utils.StringToSlice(tmpId)
				idSlice = append(idSlice,article.Id)
				updateStr := fmt.Sprintf("UPDATE %s SET doc_id=? where id=?",utils.DBInvertDoc)
				_, _ = db.Exec(updateStr, utils.SliceToString(idSlice), id)
			}
			if isExist == false {
				// 倒排索引为空，建立倒排索引
				queryStr := fmt.Sprintf("INSERT INTO %s(key_word,doc_id)VALUES (?,?)",utils.DBInvertDoc)
				log.Printf("key_word:%s invert_index is empty,insert str:%s",word,queryStr)
				_,err := db.Exec(queryStr,word,utils.SliceToString([]int{article.Id}))
				if err != nil {
					log.Printf("insert invert_index failed,err=%v,queryStr=%s",err,queryStr)
					return err
				}
			}
		}
	}
	return nil
}

func deleteTmpDoc(deleteMap map[int]bool,articles []Model.Article) []error {
	// 删除临时文件
	var res []error
	for i,article := range articles {
		if deleteMap[i] == false {
			continue
		}
		docPath := fmt.Sprintf("%s%s",utils.DocPath,article.Title)
		err := os.Remove(docPath)
		if err != nil {
			log.Printf("%s delete failed,err=%s",docPath,err)
			res = append(res,err)
		} else {
			log.Printf("delete file %s",docPath)
		}
	}
	return res
}