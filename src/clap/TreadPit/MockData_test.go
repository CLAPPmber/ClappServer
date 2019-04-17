package TreadPit

import (
	"testing"
	."clap/DataStructure"

)

var GetPracticeRecReqPars  = `{"account": "T_User"}`
var GetPracticeExpectResp  = `{"code":200,"msg":"成功获取记录","data":[{"chapter_num":3,"chapter_rec":3},{"chapter_num":1,"chapter_rec":1}]}`

var PracticeRec = []Record{
	{Chapter_num:3,Quesiont_num:2},
	{Chapter_num:3,Quesiont_num:1},
	{Chapter_num:3,Quesiont_num:3},
	{Chapter_num:1,Quesiont_num:1},
}

func InsertMockPracticeData(t *testing.T,tc *TestCase){
	tx,err := TestDb.Begin()
	if err!=nil {
		tx.Rollback()
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
		return
	}
	
	stmt,err := tx.Prepare("insert into cluser(account, password) VALUES ($1,$2);")
	defer stmt.Close()
	if err!=nil {
		tx.Rollback()
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
		return
	}

	_, err = stmt.Exec("T_User","123456")
	if err!=nil {
		tx.Rollback()
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
	}

	stmt ,err = tx.Prepare("INSERT INTO pra_record(chapter_num,question_num,account) values($1,$2,$3)")
	if err!=nil {
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
		return
	}
	for _,insertRec := range PracticeRec{
		_,err = stmt.Exec(insertRec.Chapter_num,insertRec.Quesiont_num,"T_User")
		if err!=nil {
			tx.Rollback()
			t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
			return
		}
	}
	tx.Commit()
}

func ClearMockPracticeData(t *testing.T,tc *TestCase){
	_, err := TestDb.Exec("delete from cluser  where account like 'T%';")
	if err!=nil {
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
		return
	}
	_, err = TestDb.Exec("delete from pra_record where account like 'T%';")
	if err!=nil {
		t.Fatalf("%v Prepare mock data fail err:%v",tc.Describe,err)
		return
	}
}


var CommitRecReq = `{"account": "T_User","record":[{"chapter_num": 1,"question_num": 2}]}`
var SuccessCommitRecExpectResp = `{"code":200,"msg":"提交记录成功","data":""}`
var FailCommitRecExpectResp = `{"code":505,"msg":"提交记录失败","data":""}`
