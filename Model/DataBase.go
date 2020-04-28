package Model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "19971117liu"
	ip = "127.0.0.1"
	port = "3306"
	dbName = "ReEngine"
)
var DB *sql.DB

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	db, err := sql.Open("mysql", path)
	DB = db
	if err != nil{
		log.Printf("connect mysql fail :%s\n",err)
		return
	}else{
		log.Println("connect to mysql success")
	}
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(1000)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil{
		log.Printf("open database fail,%v\n",err)
		return
	}
	log.Println("connnect success")
}
