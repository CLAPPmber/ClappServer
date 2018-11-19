package practiceModel

import (
	"database/sql"
	"errors"
	. "clap/staging/TBLogger"
	. "clap/staging/db"
	"net/http"
	"clap/staging/feedback"
	."clap/staging/TBCache"
)

func SubmitRecord(prarec Prarecord)(error){

	tx, err := Db.Begin()
	if err != nil {
		TbLogger.Error(err,"tx.begin Fail")
		return errors.New("tx.begin Fail")
	}

	var clu Cluser
	err = tx.QueryRow("SELECT * from cluser where account = $1",
		prarec.Account).Scan(&clu.Account, &clu.Password)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			TbLogger.Error(err,"账号不存在")
			return errors.New("账号不存在")
		} else {
			TbLogger.Error(err,"查询账号错误")
			return errors.New("查询账号错误")
		}
	}

	sqlStatement := `INSERT INTO pra_record(chapter_num,question_num,account) values($1,$2,$3)`
	stmt, err := tx.Prepare(sqlStatement)
	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		TbLogger.Error(err,"插入失败")
	}
	for _, onerec := range prarec.Record {
		_, err = stmt.Exec(onerec.Chapter_num, onerec.Quesiont_num, prarec.Account)
		if err != nil {
			tx.Rollback()
			TbLogger.Error(err,"提交记录失败,记录可能已存在")
			return errors.New("提交记录失败,记录可能已存在")
		}
	}
	err = tx.Commit()
	if err != nil {
		TbLogger.Error(err,"提交记录失败")
		tx.Rollback()
	}

	TbCache.DeleteCache(prarec.Account)

	return nil
}

func Getallrec(cluser Cluser)(error,[]Retprorec) {
	var retrec []Retprorec

	getCache := TbCache.GetValue(cluser.Account)	//从缓存中获取
	if getCache != nil {
		if retProrec,ok := getCache.([]Retprorec);ok{
			return nil,retProrec
		}
	}

	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()
		TbLogger.Error(err,"获取记录失败")
		return errors.New("获取记录失败"),nil
	}

	var clu Cluser
	err = tx.QueryRow("SELECT * from cluser where account = $1",
		cluser.Account).Scan(&clu.Account, &clu.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			TbLogger.Error(err,"账号不存在")
			return errors.New("账号不存在"),nil
		} else {
			tx.Rollback()
			TbLogger.Error(err,"查询账号错误")
			return errors.New("查询账号错误"),nil
		}
	}

	rows, err := tx.Query("SELECT chapter_num ,COUNT(*) FROM pra_record WHERE account = $1 group by chapter_num", cluser.Account)
	defer rows.Close()
	if err != nil {
		tx.Rollback()
		TbLogger.Error(err,"获取错误")
		return errors.New("获取错误"),nil
	}
	for rows.Next() {
		var cl Retprorec
		err := rows.Scan(&cl.Chapter_num, &cl.Chapter_rec)
		if err != nil {
			tx.Rollback()
			TbLogger.Error(err,"获取错误")
			return errors.New("获取错误"),nil
		}
		retrec = append(retrec, cl)
	}
	tx.Commit()
	TbCache.InsertCache(cluser.Account,retrec) //最新的数据插入到缓存中
	return nil,retrec
}

//清除指定账号做题记录
func ClearRecord(w http.ResponseWriter, Account string) error {
	if Account == "" {
		return  errors.New("账号为空")
	}
	fb := feedback.NewFeedBack(w)
	sqlstmt := "DELETE FROM pra_record where account = $1"
	stmt, err := Db.Prepare(sqlstmt)
	defer stmt.Close()
	if err != nil {
		TbLogger.Error("获取stmt失败", err)
		fb.SendData(501, "清除记录失败", nil)
		return err
	}

	_, err = stmt.Exec(Account)
	if err != nil {
		TbLogger.Error("清除记录失败")
		fb.SendData(500, "清除记录失败", nil)
		return err
	}
	TbLogger.Info("清除了记录")
	return nil
}