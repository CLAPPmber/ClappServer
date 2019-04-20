package main

import (
	."clap/serve"
	"clap/serve/FileService"
	"clap/serve/PracticeServe"
	"clap/serve/TestServe"
	"clap/serve/UserServe"
	_ "clap/staging/TBCache"
	"clap/staging/TBLogger"
	_ "clap/staging/db"
	_ "clap/staging/memory"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/gsql", LogApiAccess(TestServe.SqlGetHandle))
	http.HandleFunc("/clear",LogApiAccess(PracticeServe.ClearHandle))
	http.HandleFunc("/gsqls", LogApiAccess(TestServe.SqlGetsHandle))
	http.HandleFunc("/testpost", LogApiAccess(TestServe.TestPostHandle))

	http.HandleFunc("/prarecord", LogApiAccess(PracticeServe.PrarecordHandle))
	http.HandleFunc("/getallrec", LogApiAccess(PracticeServe.GetallrecHandle))
	http.HandleFunc("/login", LogApiAccess(UserServe.LoginHandle))
	http.HandleFunc("/register", LogApiAccess(UserServe.RegisteredHandle))
	http.HandleFunc("/changepassword",LogApiAccess(UserServe.ChangePasswordHandle))

	http.HandleFunc("/api/user/get_user_progress",LogApiAccess(UserServe.GerAllRecord))
	http.HandleFunc("/api/exersys/get_chapter_msg",LogApiAccess(TestServe.GetChapterMsg))
	http.HandleFunc("/api/exersys/get_progress",LogApiAccess(UserServe.GetUserProgress))
	http.HandleFunc("/api/exersys/get_rec",LogApiAccess(UserServe.GetUserTestRecord))
	http.HandleFunc("/api/exersys/put_doing_ans",LogApiAccess(TestServe.PutDoingAns))
	http.HandleFunc("/api/exersys/get_test_msg",LogApiAccess(TestServe.GetQuestionMsg))
	http.HandleFunc("/api/get_user_head_image",LogApiAccess(UserServe.GetUserHeadImage))
	http.HandleFunc("/api/update_user_head_image",LogApiAccess(UserServe.UpdateUesrHeadImage))

	http.HandleFunc("/api/upload_file",LogApiAccess(FileService.UploadFileHandler))

	dirPath,err := TBLogger.GetProDir()
	if err!=nil{
		TBLogger.TbLogger.Error("get dir path fail",err)
	}
	http.Handle("/api/get_image/", http.StripPrefix("/api/get_image/", LogDownFileAccess(http.FileServer(http.Dir(dirPath+FileService.UserHeadImageSavePath)))))


	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



