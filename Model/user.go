package Model

import (
	"fmt"
	"github.com/pkg/errors"
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
	log.Printf("Register name=%s,password=%s",name,password)
	queryStr := fmt.Sprintf("select password from %s where name=\"%s\"",utils.DBUsers,name)
	isExist := false
	rows, err := DB.Query(queryStr)
	if err != nil {
		log.Printf("db query failed,str=%s,%s",queryStr,err)
		return errors.New("网站打瞌睡了，请稍后再试~")
	}
	for rows.Next() {
		isExist = true

	}
	if isExist == true {
		log.Printf("username is exist,queryStr=%s",queryStr)
		return errors.New("该用户已经存在，请换一个名字~")
	} else {
		queryStr = fmt.Sprintf("INSERT INTO %s(name,password)VALUES (?,?)",utils.DBUsers)
		_,err := DB.Exec(queryStr,name,password)
		if err != nil {
			log.Printf("create user failed,err=%v",err)
			return err
		}
	}
	return nil
}