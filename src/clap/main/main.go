package main

import (
	"clap/config"
	_ "clap/db"
	_ "clap/logger"
	"clap/login"
	_ "clap/memory"
	"clap/serve"
	"log"
	"net/http"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	config.LoadConfig()
	http.HandleFunc("/tt2b", serve.SayhelloName)
	http.HandleFunc("/count", serve.Count)
	http.HandleFunc("/gsql", serve.SqlGet)
	http.HandleFunc("/gsqls", serve.SqlGets)
	http.HandleFunc("/testpost", serve.TestPost)
	http.HandleFunc("/prarecord", serve.Prarecord)
	http.HandleFunc("/getallrec", serve.Getallrec)
	http.HandleFunc("/login", login.LoginHandle)
	http.HandleFunc("/register", login.RegisteredHandle)
	http.HandleFunc("/changepassword",serve.ChangePassword)
	http.HandleFunc("/clearrecord",serve.ClearRecord)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
