package login

import (
	. "clap/db"
	"clap/feedback"
	."clap/TBLogger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"database/sql"
)

func CheckUser(data UserInfo) (bool, error) {

	//account = data.account

	sqlStat := "SElECT from cluser WHERE account = $1"
	stmt, err := Db.Prepare(sqlStat)
	defer stmt.Close()
	if err != nil {
		//fmt.Println("数据库语句准备失败")
		TbLogger.Error("数据库语句准备失败", err)
		return false, err
	}

	rows, err := stmt.Query(data.Account)
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}

	return false, nil

}
var userInfo UserInfo
func Registered(data UserInfo) (bool, error) {

	sqlStatement := "INSERT INTO cluser(account, password) VALUES ($1, $2);"
	stmt, err := Db.Prepare(sqlStatement)

	if err != nil {
		TbLogger.Error("插入数据语句准备失败", err)
		return false, err
	}

	_, err = stmt.Exec(data.Account,data.Password)
	defer stmt.Close()

	if err != nil {
		TbLogger.Error("插入数据失败", err)
		return false, err
	}

	return true, nil
}

func RegisteredHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb := feedback.NewFeedBack(w)
	detail, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error("读取数据失败", err)
		fb.SendData(503, "读取数据失败", "null")
		return
	}

	var data UserInfo
	err = json.Unmarshal(detail, &data)
	if err != nil {
		TbLogger.Error("解析数据失败", err)
		fb.SendData(503, "解析数据失败", "")
		return
	}

	exist, err := CheckUser(data)
	if err != nil {
		TbLogger.Error("CheckUser失败", err)
		fb.SendData(503, "解析数据失败", "")
		return
	}
	if exist {
		TbLogger.Error("账号已存在", nil)
		fb.SendData(503, "账号已存在", "")
		return
	}

	ok, err := Registered(data)
	if err != nil {
		TbLogger.Error("Registered失败", err)
		fb.SendData(503, "注册失败", "")
		return
	}
	if ok {
		fb.SendData(200, "注册成功", "")
		return
	} else {
		TbLogger.Error("result失败", err)
		fb.SendData(503, "注册失败", "")
		return
	}

}
