package login

import (
	//"APPtable"
	. "clap/db"
	"clap/feedback"
	"clap/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"database/sql"
)

func CheckUser(data interface{}) (bool, error) {

	//account = data.account

	sqlStat := "SElECT from ish2b WHERE account = $1"
	stmt, err := Db.Prepare(sqlStat)

	if err != nil {
		//fmt.Println("数据库语句准备失败")
		logger.Errorln("数据库语句准备失败", err)
		return false, err
	}

	rows, err := stmt.Query(data)

	if rows.Next() {
		return true, nil
	}

	return false, nil

}

func Registered(data interface{}) (bool, error) {

	sqlStatement := "INSERT INTO ish2b(account, password) VALUES ($1, $2);"
	stmt, err := Db.Prepare(sqlStatement)

	if err != nil {
		logger.Errorln("插入数据语句准备失败", err)
		return false, err
	}

	_, err = stmt.Exec(data)
	defer stmt.Close()

	if err != nil {
		logger.Errorln("插入数据失败", err)
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
		logger.Errorln("读取数据失败", err)
		fb.SendData(503, "读取数据失败", "null")
		return
	}

	var data interface{}
	err = json.Unmarshal(detail, &data)
	if err != nil {
		logger.Errorln("解析数据失败", err)
		fb.SendData(503, "解析数据失败", "null")
		return
	}
	exist, err := CheckUser(data)
	if err != nil {
		logger.Errorln("CheckUser失败", err)
		fb.SendData(503, "解析数据失败", "null")
		return
	}

	if exist {
		fb.SendData(503, "账号已存在", "null")
		return
	} else {
		ok, err := Registered(data)
		if err != nil {
			logger.Errorln("Registered失败", err)
			fb.SendData(503, "注册失败", "null")
			return
		}
		if ok {
			fb.SendData(200, "注册成功", "null")
			return
		} else {
			logger.Errorln("result失败", err)
			fb.SendData(503, "注册失败", "null")
			return
		}
	}
}
