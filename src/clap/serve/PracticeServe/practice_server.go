package PracticeServe

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"clap/staging/feedback"
	. "clap/staging/TBLogger"
	. "clap/model/practiceModel"
	"html/template"
)

//提交做题记录
func PrarecordHandle(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)

	postdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error(err,"ioutil.ReadAll Fail")
		fb.SendErr(err, "ioutil.ReadAll Fail")
		return
	}
	var prarec Prarecord
	err = json.Unmarshal(postdata, &prarec)
	if err != nil {
		TbLogger.Error(err,"Unmarshal Fail")
		fb.SendErr(err, "Unmarshal Fail")
		return
	}

	err = SubmitRecord(prarec)
	if err!=nil {
		fb.SendErr(err,"提交记录失败")
		return
	}

	TbLogger.Info("提交记录成功")
	fb.SendStatus(200, "提交记录成功")
}

//获取做题记录
func GetallrecHandle(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)

	reterr := []Retprorec{{Chapter_num: 0, Chapter_rec: 0}}

	postdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error(err,"获取记录失败")
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}
	var cluser Cluser
	err = json.Unmarshal(postdata, &cluser)
	if err != nil {
		TbLogger.Error(err,"获取记录失败")
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}

	err,retrec := Getallrec(cluser)
	if err!=nil{
		TbLogger.Error(err,"获取记录失败")
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}

	TbLogger.Info("成功获取记录")
	fb.SendData(200, "成功获取记录", retrec)
}

//清楚指定账号做题记录
func ClearHandle(w http.ResponseWriter, r *http.Request){
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, err := template.ParseFiles("login.html")
		if err != nil {
			TbLogger.Error(err)
			fmt.Println(err)
		}
		t.Execute(w, nil)
		w.Header().Set("Content-type", "text/html")

	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		err := r.ParseForm()
		if err != nil {
			TbLogger.Error(err)
			fmt.Println(err)
		}
		userName := template.HTMLEscapeString(r.Form.Get("username"))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("username:", userName) //输出到服务器端
		if userName == ""{
			TbLogger.Error("账号为空")
			template.HTMLEscape(w,[]byte("不能为空"))
		}
		err = ClearRecord(w,userName)
		if err!=nil{
			TbLogger.Error(err)
			template.HTMLEscape(w,[]byte("清除失败"))
			return
		}
		template.HTMLEscape(w, []byte(r.Form.Get("username")+"记录清楚成功")) //输出到客户端
		http.Redirect(w, r, "/", 302)
	}
}
