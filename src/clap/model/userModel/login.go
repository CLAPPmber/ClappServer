package userModel

import (
	. "clap/staging/db"
	. "clap/staging/TBLogger"
	"database/sql"
	."clap/DataStructure"
	"errors"
)



func Login(userInfo UserInfo) (UserInfo, error) {
	var userRet UserInfo
	sqlStatement := "SELECT account,headimage FROM cluser WHERE account = $1 AND password = $2;"
	stmt, err := Db.Prepare(sqlStatement)
	if err != nil {
		TbLogger.Error("查询出错", err)
		return userRet, errors.New("查询出错")
	}
	err = stmt.QueryRow(userInfo.Account, userInfo.Password).Scan(&userRet.Account,&userRet.UserHead)
	if err != nil {
		if err == sql.ErrNoRows {
			TbLogger.Error("账号不存在", err)
			return userRet,errors.New("账号不存在")
		} else {
			TbLogger.Error("查询出错", err)
			return userRet, errors.New("查询出错")
		}
	}
	return userRet, nil
}
