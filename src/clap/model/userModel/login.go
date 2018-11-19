package userModel

import (
	. "clap/staging/db"
	. "clap/staging/TBLogger"
	"database/sql"
)



func Login(userInfo UserInfo) (bool, string) {
	sqlStatement := "SELECT account FROM cluser WHERE account = $1 AND password = $2;"
	stmt, err := Db.Prepare(sqlStatement)
	if err != nil {
		TbLogger.Error("查询出错", err)
		return false, "查询出错"
	}
	var account string
	err = stmt.QueryRow(userInfo.Account, userInfo.Password).Scan(&account)
	if err != nil {
		if err == sql.ErrNoRows {
			TbLogger.Error("账号不存在", err)
			return false, "账号不存在"
		} else {
			TbLogger.Error("查询出错", err)
			return false, "查询出错"
		}
	}
	return true, ""
}
