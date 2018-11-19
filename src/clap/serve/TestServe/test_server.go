package TestServe

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	. "clap/staging/TBLogger"
	"clap/staging/feedback"
	. "clap/model/testModel"
)

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