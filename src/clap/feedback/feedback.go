package feedback

import (
	"encoding/json"
	"net/http"

)
//FeedBack 返回数据体
type FeedBack struct {
	W    http.ResponseWriter `json:"-"`
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data interface{}         `json:"data"`
}

//NewFeedBack new fb
func NewFeedBack(w http.ResponseWriter) FeedBack {
	return FeedBack{W: w, Code: 0, Msg: "", Data: ""}
}

//SendStatus only set code and msg
func (fb *FeedBack) SendStatus(code int, msg string) error {
	fb.Code = code
	fb.Msg = msg
	return fb.SendTo()
}

//SendData send data
func (fb *FeedBack) SendData(code int, msg string, data interface{}) error {
	fb.Code = code
	fb.Msg = msg
	fb.Data = data
	return fb.SendTo()
}

//SendErr send err
func (fb *FeedBack) SendErr(err error, msg string, data ...interface{}) {
	fb.Code = 505
	fb.Msg = msg
	if data != nil {
		fb.Data = data
	}
	fb.SendTo()
}

//SendTo send fb to request
func (fb *FeedBack) SendTo() error {
	jsonbyte, err := json.Marshal(fb)
	fb.Code = 0
	fb.Msg = ""
	fb.Data = ""
	if err != nil {
		return err
	}
	fb.W.Write(jsonbyte)
	return err
}
