package UserServe

import (
	. "clap/DataStructure"
	"clap/model/practiceModel"
	. "clap/model/userModel"
	"clap/staging/TBCache"
	. "clap/staging/TBLogger"
	"clap/staging/db"
	"clap/staging/feedback"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type TestProgress struct {
	TestNum int `json:"test_num"`
	TestCompleteCount int `json:"test_complete_count"`
}

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
		TbLogger.Error("登录失败", err)
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

//获取用户总的做题进度
func GerAllRecord(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)
	if r.Method != "GET" {
		fb.SendData(400,"Request Method no Get",nil)
		return
	}
	value := r.URL.Query()
	account := value.Get("account")
	if account==""{
		_=fb.SendData(400,"account is empty",nil)
		return
	}

	var clu Cluser
	err := db.Db.QueryRow("SELECT * from cluser where account = $1",
		account).Scan(&clu.Account, &clu.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			TbLogger.Error(err,"账号不存在")
			_ = fb.SendData(400,"账号不存在",nil)
			return
		} else {
			TbLogger.Error(err,"查询账号错误")
			_ = fb.SendData(400,"查询账号错误",nil)
			return
		}
	}


	RetData := struct {
		Complete int `json:"completed"`
		Total int `json:"total"`
	}{}

	total:=TBCache.TbCache.GetValue("TOTAL")
	if total==nil{
		err := db.Db.QueryRow(`SELECT COUNT(*) FROM clapp_test;`).Scan(&RetData.Total)
		if err!=nil{
		    TbLogger.Error("get total fail",err)
		    _=fb.SendData(400,"get total fail",nil)
		    return
		}
		TBCache.TbCache.InsertCache("TOTAL",24*30*time.Hour)
	}else{
		RetData.Total = total.(int)
	}
	err = db.Db.QueryRow(`SELECT COUNT(*) FROM pra_record WHERE  account = $1;`,account).Scan(&RetData.Complete)
	if err!=nil{
	    TbLogger.Error("get complete fail",err)
	    _=fb.SendData(400,"get complete fail",nil)
	    return
	}

	_=fb.SendData(200,"获取数据成功",RetData)
	return
}

//获取用户做题进度
func GetUserProgress(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)
	if r.Method != "GET" {
		_=fb.SendData(400,"Request Method no Get",nil)
		return
	}

	value := r.URL.Query()
	flag := value.Get("flag")
	if flag==""{
		_=fb.SendData(400,"flag is empty",nil)
		return
	}

	flagInt,err :=  strconv.Atoi(flag)
	if err!=nil{
		TbLogger.Error("flag to int fail",err)
		_=fb.SendData(400,"get flag fail,it must be a int",nil)
		return
	}

	account := value.Get("account")
	if account==""{
		_=fb.SendData(400,"account is empty",nil)
		return
	}


	cluser :=Cluser{Account:account}
	err,retrec := practiceModel.Getallrec(cluser,flagInt)
	if err!=nil{
		TbLogger.Error(err,"获取记录失败")
		fb.SendErr(err, "获取记录失败", nil)
		return
	}
	reterr := make([]TestProgress,len(retrec),len(retrec))

	for i:=0;i<len(retrec);i++{
		reterr[i].TestNum = retrec[i].Chapter_num
		reterr[i].TestCompleteCount = retrec[i].Chapter_rec
	}
	
	_=fb.SendData(200, "成功获取记录", reterr)
	
}

//获取用户做题记录
func GetUserTestRecord(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)
	if r.Method != "GET" {
		_=fb.SendData(400,"Request Method no Get",nil)
		return
	}


	value := r.URL.Query()
	flag := value.Get("flag")
	if flag==""{
		_=fb.SendData(400,"flag is empty",nil)
		return
	}

	flagInt,err :=  strconv.Atoi(flag)
	if err!=nil{
		TbLogger.Error("flag to int fail",err)
		_=fb.SendData(400,"get flag fail,it must be a int",nil)
		return
	}

	account := value.Get("account")
	if account==""{
		_=fb.SendData(400,"account is empty",nil)
		return
	}

	testNum := value.Get("test_num")
	if testNum == ""{
		_=fb.SendData(400,"test_num is empty",nil)
		return
	}
	testNumInt,err:=strconv.Atoi(testNum)
	if err!=nil{
		TbLogger.Error("get int testNumint fail ",err)
		_=fb.SendData(400,"test_number must be a int",nil)
		return
	}

	questionNum:= value.Get("question_num")
	if questionNum==""{
		_=fb.SendData(400,"question_num is empty",nil)
		return
	}

	questionNumInt,err:=strconv.Atoi(questionNum)
	if err!=nil{
		TbLogger.Error("get questionNum to int fail",err)
		_=fb.SendData(400,"question_num must be a int",nil)
		return
	}
	qa := struct {
		QuestionAns string `json:"question_ans"`
	}{}
	err = db.Db.QueryRow(
		`SELECT question_ans 
				from pra_record WHERE account = $1 and question_num = $2 and chapter_num = $3 and flag = $4;`,
				account,testNumInt,questionNumInt,flagInt).Scan(&qa.QuestionAns)
	if err!=nil{
	    TbLogger.Error("get question_ans from db fail",err)
	    _=fb.SendData(400,"get question_ans from db fail",nil)
	    return
	}
	
	_=fb.SendData(200,"获取数据成功",qa)
	return
}

//todo 获取用户头像
func GetUserHeadImage(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)

	if r.Method != "GET" {
		fb.SendData(400, "请求的方法不是get", "null")
		return
	}

	value := r.URL.Query()
	account := value.Get("account")
	if account==""{
		_=fb.SendData(400,"account is empty",nil)
		return
	}

 	var retUser UserInfo
	err := db.Db.QueryRow(`SELECT headimage FROM	cluser WHERE  account = $1`,account).Scan(&retUser.UserHead)
	if err!=nil{
	    TbLogger.Error("get user head image fail",err)
	    _=fb.SendData(400,"get user head image fail",nil)
	    return
	}

	_=fb.SendData(200,"获取用户头像成功",retUser)
	return

}

//todo 更新用户头像
func UpdateUesrHeadImage(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)

	if r.Method != "GET" {
		fb.SendData(400, "请求的方法不是GET", "null")
		return
	}

	detail, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error("读取数据失败", err)
		fb.SendData(400, "读取数据失败", "null")
		return
	}

	var data UserInfo
	err = json.Unmarshal(detail, &data)
	if err != nil {
		TbLogger.Error("解析数据失败", err)
		fb.SendData(503, "解析数据失败", "")
		return
	}

	_,err = db.Db.Exec(`UPDATE cluser SET headimage = $1 WHERE account = $2`,data.UserHead,data.Account)
	if err!=nil{
	    TbLogger.Error("update user head image fai;",err)
	    _=fb.SendData(400,"update user head image fai;",nil)
	    return
	}

	_=fb.SendData(200,"更新头像成功","")
	return
}