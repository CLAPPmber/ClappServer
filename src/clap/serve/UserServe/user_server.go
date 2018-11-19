package UserServe

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"clap/staging/feedback"
	. "clap/staging/TBLogger"
	. "clap/model/userModel"
)

//登录
func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb := feedback.NewFeedBack(w)

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error("ioutil失败", err)

	}
	var userInfo UserInfo
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		TbLogger.Error("读取数据失败", err)
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

//注册
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

//修改密码
func ChangePasswordHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb := feedback.NewFeedBack(w)

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error("ioutil失败", err)
		fb.SendData(501, "ioutil失败", "null")
		return
	}

	var userInfo UserInfo
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		TbLogger.Error("读取数据失败", err)
		fb.SendData(501, "读取数据失败", "null")
		return
	}

	err = ChangePasswor(userInfo)
	if err!=nil {
		fb.SendData(503, "修改密码失败", "null")
		return
	}

	TbLogger.Info("修改密码成功")
	fb.SendStatus(200, "修改密码成功")
}


