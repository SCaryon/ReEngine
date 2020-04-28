package Model

import (
	"fmt"
	"log"
	utils "my_go/ReEngine/util"
)

type User struct{
	Username	string
	Password	string
	Token		[]byte
}

func CheckPassWord(name, password string) (bool,string) {
	var check bool
	var warning string
	queryStr := fmt.Sprintf("select password from %s where name=\"%s\"",utils.DBUsers,name)
	var passwordTmp string
	rows, err := DB.Query(queryStr)
	if err != nil {
		log.Printf("db query failed,str=%s,%s",queryStr,err)
		check = false
		warning = "网站打瞌睡了，请稍后再试"
	}
	for rows.Next() {
		_ = rows.Scan(&passwordTmp)
	}
	if password == passwordTmp {
		check = true
	} else {
		check = false
		warning = "用户名或密码错误"
	}
	return check,warning
}

func AddUser(name,password string) error {

	return nil
}