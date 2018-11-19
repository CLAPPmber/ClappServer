package userModel

import (
	. "clap/staging/TBLogger"
	. "clap/staging/db"
	"errors"
)
func ChangePasswor(userInfo UserInfo) (error) {

	if userInfo.Account == "" {
		return errors.New("账号不能为空")
	}

	if userInfo.Password == "" {
		return errors.New("密码不能为空")
	}

	sqlState := "UPDATE cluser SET password = $1 WHERE account = $2;"
	stmt, err := Db.Prepare(sqlState)
	defer stmt.Close()
	if err != nil {
		TbLogger.Error("获取更新stmt失败", err)
		return err
	}

	_, err = stmt.Exec(userInfo.Password, userInfo.Account)
	if err != nil {
		TbLogger.Error("修改密码Exec失败")
		return err
	}
	return nil
}