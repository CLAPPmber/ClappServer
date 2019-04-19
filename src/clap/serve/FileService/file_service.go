package FileService

import (
	. "clap/staging/TBLogger"
	"clap/staging/feedback"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const UserHeadImageSavePath = "/head_image/"

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//POST takes the uploaded file(s) and saves it to disk.
	fb := feedback.NewFeedBack(w)
	if r.Method != "POST" {
		TbLogger.Error("请求方法出错",nil)
		_=fb.SendData(400,"请求方法出错",nil)
		return
	}
	//设置接收数据存放在内存所占用的最大空间
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
	defer file.Close()

	//文件扩展名
	fileext := filepath.Ext(handler.Filename)
	if fileext!=".jpg" && fileext!=".png"{
		TbLogger.Error("文件类型错误",fileext)
		_=fb.SendData(400,"文件类型错误",nil)
		return
	}

	//用时间戳做文件名防止重名
	filename := strconv.FormatInt(time.Now().Unix(), 10) +handler.Filename

	//创建文件夹
 	dirPath,err := Mkdir(UserHeadImageSavePath)
 	if err!=nil{
 	    TbLogger.Error("判断文件夹是否存在错误",err)
 	    _=fb.SendData(400,"判断文件夹是否存在错误",nil)
 	    return
 	}

	//新建文件
	f, _ := os.OpenFile(dirPath+UserHeadImageSavePath+filename, os.O_CREATE|os.O_RDWR, 0660)
	//保存文件
	_, err = io.Copy(f, file)
	fmt.Println(dirPath+UserHeadImageSavePath+filename)
	f.Close()
	//压缩图片
	if handler.Size>=1500000{
		err = ResizeImage(dirPath+UserHeadImageSavePath+filename,fileext)
		if err!=nil{
			TbLogger.Error("压缩图片失败",err)
		}
	}

	_=fb.SendData(200,"上传成功",dirPath+UserHeadImageSavePath+filename)
	return
}


func Mkdir(path string) (string,error) {

	dirPath ,err := GetProDir()
	if err!=nil{
		return dirPath,err
	}
	err = MkdirProDir(dirPath+path)
	if err != nil {
		return dirPath,err
	}
	return dirPath,err
}

//压缩图片
func ResizeImage(path string,fileext string)(err error){
	// decode jpeg into image.Image
	var img image.Image
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	if fileext == ".jpg"{
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	}else{
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	}

	defer  file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(800, 0, img, resize.NearestNeighbor)

	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)
	return nil
}