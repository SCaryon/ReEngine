package Engine

import (
	"ReEngine/Model"
	utils "ReEngine/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func AddDocment() {

}

func DelDocment() {

}

// 将存放在tmp文件中的临时文档读取并进行存储和生成倒排索引然后删除文档文件
func UpdateIndex() bool {
	articles := Model.GetAndReadFiles(utils.DocPath)
	deleteMap,err := InsertDoc(articles)
	if err != nil {
		log.Fatal(err)
		return false
	}
	err = CreateInvert(articles)
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
func UpdateDoc(id int,doc Model.Article) (int,error) {
	// 对于修改文章的操作，先删除原有的文档，再插入新的内容，然后插入倒排索引
	log.Printf("UpdateDoc %d,doc : %s",id,doc.Title)
	err := Model.DeleteDoc(id)
	if err != nil {
		return -1,err
	}
	newId,err := Model.InsertDoc(doc)
	if err != nil {
		return newId,err
	}
	doc.Id = newId
	err = CreateInvert([]Model.Article{doc})
	// 保证Cache数据一致性
	jsonDoc,err := json.Marshal(doc)
	if err == nil {
		utils.BigCache.Set(fmt.Sprintf(utils.CacheDocContent,newId),jsonDoc)
	}
	// 保证redis的数据一致性
	titleKey := fmt.Sprintf(utils.RedisDoctitleSeg,doc.Id)
	tmpSeg := utils.SegmentContent(doc.Title)
	_ = Model.RedisSet(titleKey, tmpSeg)
	contentKey := fmt.Sprintf(utils.RedisDocContentSeg,doc.Id)
	tmpSeg = utils.SegmentContent(doc.Content)
	_ = Model.RedisSet(contentKey, tmpSeg)

	if err != nil {
		return newId,err
	}
	return newId,nil
}
func InsertDoc(articles []Model.Article) (map[int]bool,error) {
	deleteMap := make(map[int]bool)
	// 将文档放入数据库
	for i,article := range articles {
		id,err := Model.InsertDoc(article)
		if err != nil {
			// 对于未能成功放入数据库的文档暂时不删除
			deleteMap[i] = true
		} else {
			articles[i].Id = id
			log.Printf("insert file,Title=%s,Auth=%s,Id=%d", article.Title, article.Auth, id)
		}
	}
	return deleteMap,nil
}

func CreateInvert(articles []Model.Article) error{
	// 创建倒排索引
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
			tmpInvert,isExist,err := Model.SearchInvertDB(word)
			if err != nil {
				return err
			}
			if isExist == false {
				// 倒排索引为空，建立倒排索引
				_ = Model.InsertInvert(word, article.Id)
			} else {
				idSlice := utils.StringToSlice(tmpInvert.NumDocs)
				idSlice = append(idSlice,article.Id)
				tmpInvert.NumDocs = utils.SliceToString(idSlice)
				err = Model.UpdateInvert(tmpInvert)
			}
		}
	}
	return nil
}

func deleteTmpDoc(deleteMap map[int]bool,articles []Model.Article) []error {
	// 删除临时文件
	var res []error
	for i,article := range articles {
		if deleteMap[i] == true {
			continue
		}
		docPath := fmt.Sprintf("%s%s.%s",utils.DocPath,article.Title,article.Auth)
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