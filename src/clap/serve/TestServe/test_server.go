package TestServe

import (
	"clap/staging/db"
	"database/sql"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	. "clap/staging/TBLogger"
	"clap/staging/feedback"
	. "clap/model/testModel"
	"strconv"
)

type ChapterMsg struct {
	TestNum int `json:"test_num"`
	TestName string `json:"test_name"`
	TestTotal 	int `json:"test_total"`
}

//SqlGet get data-获取当前所有用户
func SqlGetsHandle(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)
	clusers,err := SqlGets()
	if err!=nil {
		fb.SendErr(err,"调用SqlGets请求数据失败")
		return
	}
	TbLogger.Info("调用SqlGets")
	fb.SendData(200, "Request data !", clusers)
}

//获取用户
func SqlGetHandle(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)
	cluser,err := SqlGet()
	if err!=nil {
		TbLogger.Error("获取数据失败",err)
		fb.SendErr(err,"获取s数据失败")
		return
	}
	TbLogger.Info("调用SqlGets",err)
	fb.SendData(200,"SqlGet!",cluser)
}

//将收到的数据在重新发送回去
func TestPostHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var at interface{}
		postdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			TbLogger.Error(err)
			return
		}
		err = json.Unmarshal(postdata, &at)
		if err != nil {
			TbLogger.Error(err)
			return
		}
		stringsdata := bytes.NewBuffer(postdata).String()
		TbLogger.Info("post request ,requestdata:", stringsdata)
		if err != nil {
			TbLogger.Error(err)
		}
		fb := feedback.NewFeedBack(w)
		fb.SendData(200, "Post data", at)
	}
}

//获取习题目录信息
func GetChapterMsg(w http.ResponseWriter,r *http.Request){
	var chapterMsg []ChapterMsg
	chapterMsg = make([]ChapterMsg,0)
	fb := feedback.NewFeedBack(w)
	if r.Method != "GET" {
		_=fb.SendData(400,"Request Method no Get",chapterMsg)
		return
	}

	value := r.URL.Query()

	flag := value.Get("flag")
	if flag==""{
		_=fb.SendData(400,"flag is empty",chapterMsg)
		return
	}
	flagInt,err :=  strconv.Atoi(flag)
	if err!=nil{
		TbLogger.Error("flag to int fail",err)
		_=fb.SendData(400,"get flag fail,it must be a int",chapterMsg)
		return
	}

	row,err := db.Db.Query(`SELECT chapter_num,chapter_name,COUNT(chapter_num) from clapp_test WHERE flag = $1 group by chapter_num,chapter_name;`,flagInt)
	if err!=nil{
		if err==sql.ErrNoRows{
			TbLogger.Error("no record",err)
			_=fb.SendData(400,"no record,if the params is err?",nil)
			return
		}
	    TbLogger.Error("get chapter_msg from db fail",err)
	    _=fb.SendData(400,"get chapter_msg from db fail",chapterMsg)
	    return
	}
	
	defer row.Close()
	for row.Next(){
		var cm ChapterMsg
		err = row.Scan(
			&cm.TestNum,
			&cm.TestName,
			&cm.TestTotal)
		if err!=nil{
		    TbLogger.Error("scan data fail",err)
		    _=fb.SendData(400,"get chapter msg from db fail",chapterMsg)
		    return
		}
		chapterMsg = append(chapterMsg,cm)
	}
	_=fb.SendData(200,"获做数据成功",chapterMsg)
	return
}

//获取题目信息
func GetQuestionMsg(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)
	if r.Method != "GET" {
		fb.SendData(400,"Request Method no Get",nil)
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


	var qn,qa,qb,qc,qd,qans string
	err = db.Db.QueryRow(`
		SELECT question_name,question_a,question_b,question_c,question_d,question_ans 
		FROM clapp_test WHERE chapter_num = $1 and question_num = $2 and flag = $3;`,testNumInt,questionNumInt,flagInt).
		Scan(&qn,&qa,&qb,&qc,&qd,&qans)
	if err!=nil{
		if err==sql.ErrNoRows{
			TbLogger.Error("no record",err)
			_=fb.SendData(400,"no record,if the params is err?",nil)
			return
		}
		TbLogger.Error("get test msg fail",err)
		_=fb.SendData(400,"get test msg fail because get from db fail",nil)
		return
	}
	QuestionMsg:= struct {
		QuestionName string `json:"question_name"`
		QuestionA	string 	`json:"question_a"`
		QuestionB 	string 	`json:"question_b"`
		QuestionC   string 	`json:"question_c"`
		QuestionD 	string 	`json:"question_d"`
		QuestionAns string 	`json:"question_ans"`
	}{
		QuestionName:qn,
		QuestionA:qa,
		QuestionB:qb,
		QuestionC:qc,
		QuestionD:qd,
		QuestionAns:qans,
	}
	
	_=fb.SendData(200,"获取数据成功",QuestionMsg)
	return
}

//提交做题记录
func PutDoingAns(w http.ResponseWriter,r *http.Request){
	fb := feedback.NewFeedBack(w)
	if r.Method != "POST" {
		TbLogger.Error("Request Method no Post")
		_=fb.SendData(400,"Request Method no Post",nil)
		return
	}

	BodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		TbLogger.Error("get http request body data fail",err)
		_=fb.SendData(400,"get http request body data fail",nil)
		return
	}

	PostData := struct {
		Account string  		`json:"account"`
		Record []struct{
			Flag 	int 		`json:"flag"`
			TestNum	int			`json:"test_num"`
			QuestionNum	int		`json:"question_num"`
			QuestionAns string  `json:"question_ans"`
		} 						`json:"record"`
	}{}

	err = json.Unmarshal([]byte(BodyData), &PostData)
	if err != nil {
		TbLogger.Error("get postdata from json format fail",err)
		_=fb.SendData(400,"get post_data fail,is json format err?",nil)
		return
	}

	if PostData.Account==""{
		TbLogger.Error("account is empty")
		_=fb.SendData(400,"account is empty",nil)
		return
	}

	if len(PostData.Record)<=0{
		TbLogger.Error("recode is nil")
		_=fb.SendData(400,"recode is nil",nil)
		return
	}

	sqlInsert := `INSERT INTO pra_record(chapter_num, question_num, account, flag,question_ans) VALUES ($1,$2,$3,$4,$5)`
	stmt,err := db.Db.Prepare(sqlInsert)
	if err!=nil{
		if err==sql.ErrNoRows{
			TbLogger.Error("no record",err)
			_=fb.SendData(400,"no record,if the params is err?",nil)
			return
		}
	    TbLogger.Error("sql prepare fail",err)
	    _=fb.SendData(400,"db fail",err)
	    return
	}

	for i:=0;i<len(PostData.Record);i++{
		_,err = stmt.Exec(
			PostData.Record[i].TestNum,
			PostData.Record[i].QuestionNum,
			PostData.Account,
			PostData.Record[i].Flag,
			PostData.Record[i].QuestionAns,
			)
		if err!=nil{
		    TbLogger.Error("insert pra record fail",err)
		    _=fb.SendData(400,"insert pra record fail",nil)
		    return
		}
	}
	
	_=fb.SendData(200,"提交记录成功",nil)
	return
}