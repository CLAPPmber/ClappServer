package login

import (
	. "clap/db"
	"clap/feedback"
	"clap/logger"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserInfo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func Login(userInfo UserInfo) (bool, string) {
	sqlStatement := "SELECT account FROM cluser WHERE account = $1 AND password = $2;"
	stmt, err := Db.Prepare(sqlStatement)
	if err != nil {
		logger.Errorln("查询出错", err)
		return false, "查询出错"
	}
	var account string
	err = stmt.QueryRow(userInfo.Account, userInfo.Password).Scan(&account)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorln("账号不存在", err)
			return false, "账号不存在"
		} else {
			logger.Errorln("查询出错", err)
			return false, "查询出错"
		}
	}
	return true, ""
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb := feedback.NewFeedBack(w)

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln("ioutil失败", err)

	}
	var userInfo UserInfo
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		logger.Errorln("读取数据失败", err)
		fb.SendData(501, "读取数据失败", "null")
		return
	}

	ok, msg := Login(userInfo)
	if ok {
		fb.SendData(200, "登录成功", "null")
		return
	} else {
		fb.SendData(501, msg, "null")
		return
	}
}
