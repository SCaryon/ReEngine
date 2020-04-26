package Engine

import (
	"fmt"
	"log"
	utils "my_go/ReEngine/util"
	"os"
)

func AddDocment() {

}

func DelDocment() {

}

// 将存放在tmp文件中的临时文档读取并进行存储和生成倒排索引然后删除文档文件
func UpdateIndex() bool {
	files := utils.GetAndReadFiles(utils.DocPath)
	db := utils.DB
	if db == nil {
		log.Println("connect db failed")
		return false
	}
	delete := make(map[int]bool)
	// 将文档放入数据库
	for i,file := range files {
		log.Printf("insert file,Title=%s,Auth=%s,Content=%s",file.Title,file.Auth,file.Content)
		queryStr := fmt.Sprintf("INSERT INTO %s(title,auth,context,create_time)VALUES (?,?,?)",utils.DBDocment)
		result,err := db.Exec(queryStr,file.Title,file.Auth,file.Content,file.CreateTime)
		if err != nil {
			log.Printf("insert failed,err=%v",err)
			continue
		}
		id,err := result.LastInsertId()
		file.Id = int(id)
		// 对于未能成功放入数据库的文档暂时不删除
		delete[i] = true
	}

	// 创建倒排索引
	for _,it := range files {
		seg := utils.SegmentContent(fmt.Sprintf("%s %s %s",it.Title,it.Content,it.Auth))
		for _,word := range seg {
			queryStr := fmt.Sprintf("select id,doc_id from %s where key_word=%s\n",utils.DBInvertDoc,word)
			rows,err := db.Query(queryStr)
			if err == nil || rows == nil {
				log.Printf("use %s table ,query = %s failed\n",utils.DBInvertDoc,word)
				continue
			}
			for rows.Next() {
				var id int
				var tmpId string
				err := rows.Scan(&id, &tmpId)
				if err != nil {
					log.Printf("get data failed, error:[%v]\n", err.Error())
				}
				log.Println(id, tmpId)
				idSlice := utils.StringToSlice(tmpId)
				idSlice = append(idSlice,it.Id)
				updateStr := fmt.Sprintf("UPDATE %s SET doc_id=? where id=?",utils.DBInvertDoc)
				db.Exec(updateStr,utils.SliceToString(idSlice),id)
			}

		}
	}

	// 删除临时文件
	for i,it := range files {
		if delete[i] == false {
			continue
		}
		docPath := fmt.Sprintf("%s%s",utils.DocPath,it.Title)
		err := os.Remove(docPath)
		if err != nil {
			log.Printf("%s delete failed,err=%s",docPath,err)
		} else {
			log.Printf("delete file %s",docPath)
		}
	}
	return true
}