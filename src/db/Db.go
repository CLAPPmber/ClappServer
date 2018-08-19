package db

import (
	"database/sql"
	"fmt"
	"logger"

	_ "github.com/lib/pq"
)

var (
	DbDriverName  = "postgres"
	StartDatabase = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	Host          = "123.207.25.239"
	Port          = 5432
	User          = "ish2b"
	Dbname        = "ish2b"
	Password      = "123456"
	err           error
)

var Db *sql.DB

func init() {
	psqlInit := fmt.Sprintf(StartDatabase, Host, Port, User, Password, Dbname)
	Db, err = sql.Open(User, psqlInit)
	if err != nil {
		logger.Errorln("数据库启动失败：", err)

	} else {
		fmt.Println("数据库连接成功")
	}
}
