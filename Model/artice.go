package Model

import (
	utils "ReEngine/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Article struct {
	Id         int
	Title      string
	Auth       string
	Content    string
	CreateTime int
}


func GetAllDocs() ([]Article,error) {
	var resp []Article
	queryStr := fmt.Sprintf("select id,title,auth,context,create_time from %s where is_delete=0",utils.DBDocument)
	rows,err := DB.Query(queryStr)
	if err != nil || rows == nil {
		log.Printf("use %s table ,query = %s failed\n",utils.DBDocument,queryStr)
		return nil,err
	}
	//log.Println("queryStr",queryStr)
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
		tmpArticle := Article{Id:tmpId, Auth:tmpAuth, Title:tmpTitle, Content:tmpContent, CreateTime:tmpTime}
		// 存储分词结果
		// todo 从redis里面拿到文章的分词的信息
		resp = append(resp,tmpArticle)
	}
	return  resp,nil
}

func DeleteDoc(docId int) error {
	updateStr := fmt.Sprintf("update %s set is_delete=1 where id=?",utils.DBDocument)
	_, err := DB.Exec(updateStr, docId)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}



func GetDocByIds(ids []int) ([]Article,error){
	var resp []Article
	for _,id := range ids {
		queryStr := fmt.Sprintf("select id,title,auth,context,create_time from %s where id=%d and is_delete=0",utils.DBDocument,id)
		rows,err := DB.Query(queryStr)
		if err != nil || rows == nil {
			log.Printf("use %s table ,query = %s failed\n",utils.DBDocument,queryStr)
			return nil,err
		}
		//log.Println("queryStr",queryStr)
		//定义变量接收查询数据
		tmpId := 0
		tmpTitle := ""
		tmpAuth := ""
		tmpContent := ""
		tmpTime := 0
		for rows.Next() {
			err := rows.Scan(&tmpId,&tmpTitle,&tmpAuth,&tmpContent,&tmpTime)
			if err != nil {
				log.Printf("get data failed, error:[%v]\n", err.Error())
			}
		}
		// 过滤已经被删除的文档
		if tmpId != 0 {
			tmpArticle := Article{Id:tmpId, Auth:tmpAuth, Title:tmpTitle, Content:tmpContent, CreateTime:tmpTime}
			resp = append(resp,tmpArticle)
		}
	}
	return resp,nil
}

func CountDocs() (int,error) {
	queryStr := fmt.Sprintf("select count(id) from %s where is_delete=0",utils.DBDocument)
	rows,err := DB.Query(queryStr)
	if err != nil || rows == nil {
		log.Printf("count lines failed ,query:%s",queryStr)
		return -1 ,err
	}
	var dataNum int
	for rows.Next() {
		err := rows.Scan(&dataNum)
		if err != nil {
			log.Printf("get data failed, error:[%v]\n", err.Error())
			return -1,err
		}
	}
	return dataNum,nil
}

func GetAndReadFiles(filePath string) []Article {
	var articles []Article
	goPath := os.Getenv("GOPATH")
	path := fmt.Sprintf("%s/src/my_go/ReEngine/%s", goPath, filePath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		txt, err := ioutil.ReadFile(path + file.Name())
		if err != nil {
			panic(err)
		}
		// 将字节流转换为字符串
		content := string(txt)
		createTime := utils.GetFileCreateTime(file)
		strTmp := strings.Split(file.Name(),".")
		lenTmp := len(strTmp)
		fileName := ""
		auth := strTmp[lenTmp-1]
		for index := 0;index < lenTmp-1;index++ {
			fileName = fileName + strTmp[index]
		}
		fmt.Println("get tmp file:",file.Name())
		fmt.Printf("file info title:%s,auth:%s",fileName,auth)
		articles = append(articles, Article{Title: fileName,Auth:auth, Content: content, CreateTime: createTime})

	}
	return articles
}

func InsertDoc(article Article) (int,error) {
	queryStr := fmt.Sprintf("INSERT INTO %s(title,auth,context,create_time)VALUES (?,?,?,?)",utils.DBDocument)
	result,err := DB.Exec(queryStr,article.Title,article.Auth,article.Content,article.CreateTime)
	if err != nil {
		log.Printf("insert failed,err=%v",err)
		return -1, err
	}
	id,err := result.LastInsertId()
	log.Printf("insert file,Title=%s,Auth=%s,Id=%d",article.Title,article.Auth,int(id))
	return int(id),nil
}