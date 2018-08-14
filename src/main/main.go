package main

import (
	"config"
	_ "db"
	"log"
	_ "logger"
	_ "memory"
	"net/http"
	"serve"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	config.LoadConfig()
	http.HandleFunc("/tt2b", serve.SayhelloName)
	http.HandleFunc("/login", serve.Login)
	http.HandleFunc("/count", serve.Count)
	http.HandleFunc("/gsql", serve.SqlGet)
	http.HandleFunc("/gsqls", serve.SqlGets)
	http.HandleFunc("/testpost", serve.TestPost)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
