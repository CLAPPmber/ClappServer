package testModel

import (
	. "clap/staging/TBLogger"
	. "clap/staging/db"
	"errors"
	."clap/DataStructure"
)
func SqlGets()([]Cluser,error){
	var clusers []Cluser
	rows, err := Db.Query("select account,password from cluser")
	defer rows.Close()
	if err != nil {
		TbLogger.Error(err)
		return nil,err
	}
	for rows.Next() {
		account := ""
		password := ""
		err := rows.Scan(&account, &password)
		if err != nil {
			TbLogger.Error("Sacn err:", err)
			return nil,err
		}
		clusers = append(clusers, Cluser{Account: account, Password: password})
	}
	return clusers,err
}

func SqlGet()(Cluser,error){
	var cluser Cluser
	rows, err := Db.Query("select account,password from cluser")
	defer rows.Close()
	if err != nil {
		TbLogger.Error("Query err:",err)
		return Cluser{},errors.New("Query err")
	}
	for rows.Next() {
		account := ""
		password := ""
		err := rows.Scan(&account, &password)
		if err != nil {
			TbLogger.Error("Sacn err:",err)
			return  Cluser{},errors.New("Sacn err")
		}
		cluser = Cluser{Account: account, Password: password}
	}
	return cluser,nil
}