package userModel

import (
	. "clap/staging/db"
	. "clap/staging/TBLogger"
	//"database/sql"
)

func CheckUser(data UserInfo) (bool, error) {

	//account = data.account

	sqlStat := "SElECT from cluser WHERE account = $1"
	stmt, err := Db.Prepare(sqlStat)
	defer stmt.Close()
	if err != nil {
		//fmt.Println("数据库语句准备失败")
		TbLogger.Error("数据库语句准备失败", err)
		return false, err
	}

	rows, err := stmt.Query(data.Account)
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}

	return false, nil

}

func Registered(data UserInfo) (bool, error) {

	sqlStatement := "INSERT INTO cluser(account, password) VALUES ($1, $2);"
	stmt, err := Db.Prepare(sqlStatement)

	if err != nil {
		TbLogger.Error("插入数据语句准备失败", err)
		return false, err
	}

	_, err = stmt.Exec(data.Account,data.Password)
	defer stmt.Close()

	if err != nil {
		TbLogger.Error("插入数据失败", err)
		return false, err
	}

	return true, nil
}


