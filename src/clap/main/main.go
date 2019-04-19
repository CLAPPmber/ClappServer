package main

import (
	"clap/serve/FileService"
	"clap/serve/PracticeServe"
	"clap/serve/TestServe"
	"clap/serve/UserServe"
	_ "clap/staging/TBCache"
	_ "clap/staging/db"
	_ "clap/staging/memory"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/gsql", TestServe.SqlGetHandle)
	http.HandleFunc("/clear",PracticeServe.ClearHandle)
	http.HandleFunc("/gsqls", TestServe.SqlGetsHandle)
	http.HandleFunc("/testpost", TestServe.TestPostHandle)

	http.HandleFunc("/prarecord", PracticeServe.PrarecordHandle)
	http.HandleFunc("/getallrec", PracticeServe.GetallrecHandle)
	http.HandleFunc("/login", UserServe.LoginHandle)
	http.HandleFunc("/register", UserServe.RegisteredHandle)
	http.HandleFunc("/changepassword",UserServe.ChangePasswordHandle)

	http.HandleFunc("/api/user/get_user_progress",UserServe.GerAllRecord)
	http.HandleFunc("/api/exersys/get_chapter_msg",TestServe.GetChapterMsg)
	http.HandleFunc("/api/exersys/get_progress",UserServe.GetUserProgress)
	http.HandleFunc("/api/exersys/get_rec",UserServe.GetUserTestRecord)
	http.HandleFunc("/api/exersys/put_doing_ans",TestServe.PutDoingAns)
	http.HandleFunc("/api/exersys/get_test_msg",TestServe.GetQuestionMsg)
	http.HandleFunc("/api/get_user_head_image",UserServe.GetUserHeadImage)
	http.HandleFunc("/api/update_user_head_image",UserServe.UpdateUesrHeadImage)

	http.HandleFunc("/api/upload_file",FileService.UploadFileHandler)



	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



