package FileService

import (
	. "clap/staging/TBLogger"
	"clap/staging/feedback"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const UserHeadImageSavePath = "./head_image/"

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//POST takes the uploaded file(s) and saves it to disk.
	fb := feedback.NewFeedBack(w)
	if r.Method != "POST" {
		TbLogger.Error("请求方法出错",nil)
		_=fb.SendData(400,"请求方法出错",nil)
		return
	}
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		TbLogger.Error("设置文件最多字节保存内存失败",err)
		_=fb.SendData(400,"设置文件最多字节保存内存失败",nil)
		return
	}

	//获取表单中的文件
	file, handler, err := r.FormFile("file")
	if err!=nil{
		TbLogger.Error("获取文件失败",err)
		_=fb.SendData(400,"获取文件失败",nil)
		return
	}

	//文件扩展名
	fileext := filepath.Ext(handler.Filename)
	fmt.Println(fileext)
	//用时间戳做文件名防止重名
	filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext

 	dirExist,err := PathExists(UserHeadImageSavePath)
 	if err!=nil{
 	    TbLogger.Error("判断文件夹是否存在错误",err)
 	    _=fb.SendData(400,"判断文件夹是否存在错误",nil)
 	    return
 	}

 	if !dirExist{
 		err := os.Mkdir(UserHeadImageSavePath,os.ModePerm)
 		if err!=nil{
 		    TbLogger.Error("创建文件夹失败",err)
 		    _=fb.SendData(400,"创建文件夹失败",nil)
 		    return
 		}
	}

	//新建文件
	f, _ := os.OpenFile(UserHeadImageSavePath+filename, os.O_CREATE|os.O_WRONLY, 0660)
	//保存文件
	_, err = io.Copy(f, file)

}


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
