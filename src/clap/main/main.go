package main

import (
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
	http.HandleFunc("/logn", UserServe.LoginHandle)
	http.HandleFunc("/register", UserServe.RegisteredHandle)
	http.HandleFunc("/changepassword",UserServe.ChangePasswordHandle)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}



