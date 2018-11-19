package main

import (
	_ "clap/staging/db"
	_ "clap/staging/memory"
	"net/http"
	"log"
	"clap/serve/UserServe"
	"clap/serve/PracticeServe"
	"clap/serve/TestServe"
	_"clap/staging/TBCache"
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
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


