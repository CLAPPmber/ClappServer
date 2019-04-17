package TreadPit

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var TestDb *sql.DB

const dbDriverName = "postgres"
const dbStartDataBase = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
const dbHost = "123.207.25.239"
const dbPort = 5432
const dbUser = "ish2b"
const dbName = "ish2b"
const dbPassword = "123456"

func InitTestDb(){
	DbInfo := fmt.Sprintf(dbStartDataBase, dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	TestDb, err = sql.Open(dbDriverName, DbInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TestDb Open success!\n begin pind db serve %v:%v \n",dbHost,dbPort)
	err = TestDb.Ping()
	if err!=nil{
		fmt.Println("ping fail")
		TestDb.Close()
		log.Fatal(err)
	}
	fmt.Println("ping success!")
}

func GetTestDb() *sql.DB{
	if TestDb == nil{
		InitTestDb()
	}
	return TestDb
}

func ClearTestData(){
	//DoSomeClear
	_, err := TestDb.Exec("delete from cluser  where account like 'T%';")
	if err!=nil {
		fmt.Println("clear mock_data fail:",err)
		return
	}
	_, err = TestDb.Exec("delete from pra_record where account like 'T%';")
	if err!=nil {
		fmt.Println("clear mock_data fail:",err)
		return
	}
}