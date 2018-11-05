package main

import (
	_ "clap/db"
	_ "clap/memory"
	."clap/TBLogger"
	"net/http"
	"clap/serve"
	"clap/login"
	"log"
	"strconv"
	"time"
)

func main() {
	go func() {
		for i := 0;i<100;i++{
			TbLogger.Info("插入日志测试:"+strconv.Itoa(i))
			time.Sleep(1*time.Hour)
		}
	}()
	http.HandleFunc("/tt2b", serve.SayhelloName)
	http.HandleFunc("/count", serve.Count)
	http.HandleFunc("/gsql", serve.SqlGet)
	http.HandleFunc("/clear",serve.Clear)
	http.HandleFunc("/gsqls", serve.SqlGets)
	http.HandleFunc("/testpost", serve.TestPost)
	http.HandleFunc("/prarecord", serve.Prarecord)
	http.HandleFunc("/getallrec", serve.Getallrec)
	http.HandleFunc("/login", login.LoginHandle)
	http.HandleFunc("/register", login.RegisteredHandle)
	http.HandleFunc("/changepassword",serve.ChangePassword)
	http.HandleFunc("/clearrecord",serve.Clear)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


