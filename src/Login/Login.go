package login

import(
	"feedback"
	"logger"
	. "db"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

func Login(data interface{})(bool,error){
	sqlStatement := "SELECT * FROM ish2b WHERE account = $1 AND password = $2;"
	stmt, err := Db.Prepare(sqlStatement)
	if err != nil {
		logger.Errorln("Login失败",err)
		return false,err
	}
	rows, err := stmt.Query(data)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(data)
		return true, nil
	}
	logger.Errorln("无法查询",err)
    return false,err
}

func LoginHandle(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb:=feedback.NewFeedBack(w)

	result,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		logger.Errorln("ioutil失败",err)
		
	}
	var data interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		logger.Errorln("读取数据失败",err)
		fb.SendData(501,"读取数据失败","null")
		return
	}
	fmt.Println(data)
	
	ok,err := Login(data)
	if err !=nil {
		logger.Errorln("Login失败",err)
		fb.SendData(501,"解析数据失败","null")
		return
	}

	if ok {
		fb.SendData(200,"登录成功","null")
		return
	}else{
		fb.SendData(501,"账号不存在","null")
		return
	}
}