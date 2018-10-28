package serve

import (
	"bytes"
	. "clap/db"
	"clap/feedback"
	"clap/logger"
	. "clap/login"
	"clap/session"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
	"errors"
)

//SayhelloName for test http
func SayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello astaxie!")
}

//count ???
func count(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		session.GlobalSessions.SessionDestory(w, r)
		sess = session.GlobalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

func Count(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

//SqlGet get data
func SqlGets(w http.ResponseWriter, r *http.Request) {
	var clusers []Cluser
	rows, err := Db.Query("select * from cluser")
	if err != nil {
		logger.Errorln("Query err:", err)
		return
	}
	for rows.Next() {
		account := ""
		password := ""
		err := rows.Scan(&account, &password)
		if err != nil {
			logger.Errorln("Sacn err:", err)
			return
		}
		clusers = append(clusers, Cluser{Account: account, Password: password})
	}
	fb := feedback.NewFeedBack(w)
	fb.SendData(200, "Request data !  last test 2", clusers)
}

func SqlGet(w http.ResponseWriter, r *http.Request) {
	var cluser Cluser
	rows, err := Db.Query("select * from cluser")
	if err != nil {
		logger.Errorln("Query err:", err)
		return
	}
	for rows.Next() {
		account := ""
		password := ""
		err := rows.Scan(&account, &password)
		if err != nil {
			logger.Errorln("Sacn err:", err)
			return
		}
		cluser = Cluser{Account: account, Password: password}
	}
	jsonuse, _ := json.Marshal(cluser)
	fmt.Fprintln(w, string(jsonuse))
}

//
func TestPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var at interface{}
		postdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorln(err)
			return
		}
		err = json.Unmarshal(postdata, &at)
		if err != nil {
			logger.Errorln(err)
			return
		}
		stringsdata := bytes.NewBuffer(postdata).String()
		logger.Infoln("post request ,requestdata:", stringsdata)
		if err != nil {
			logger.Errorln(err)
		}
		fb := feedback.NewFeedBack(w)
		fb.SendData(200, "Post data", at)
	}
}

//提交做题记录
func Prarecord(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)
	tx, err := Db.Begin()
	if err != nil {
		fb.SendErr(err, "提交记录失败")
		return
	}
	defer GetPanic(tx)
	postdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fb.SendErr(err, "提交记录失败")
		return
	}
	var prarec prarecord
	err = json.Unmarshal(postdata, &prarec)
	if err != nil {
		fb.SendErr(err, "提交记录失败")
		return
	}
	var clu Cluser
	err = tx.QueryRow("SELECT * from cluser where account = $1",
		prarec.Account).Scan(&clu.Account, &clu.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			fb.SendStatus(501, "账号不存在")
			return
		} else {
			fb.SendErr(err, "查询账号错误")
			return
		}
	}

	sqlStatement := `INSERT INTO pra_record(chapter_num,question_num,account) values($1,$2,$3)`
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		fb.SendErr(err, "插入失败")
	}
	for _, onerec := range prarec.Record {
		_, err = stmt.Exec(onerec.Chapter_num, onerec.Quesiont_num, prarec.Account)
		if err != nil {
			tx.Rollback()
			fb.SendErr(err, "提交记录失败,记录可能已存在")
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		fb.SendErr(err, "提交记录失败")
		tx.Rollback()
	}
	fb.SendStatus(200, "提交记录成功")
}

func Getallrec(w http.ResponseWriter, r *http.Request) {
	fb := feedback.NewFeedBack(w)
	tx, err := Db.Begin()
	var retrec []Retprorec
	reterr := []Retprorec{{Chapter_num: 0, Chapter_rec: 0}}
	if err != nil {
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}
	defer GetPanic(tx)
	postdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}
	var cluser Cluser
	err = json.Unmarshal(postdata, &cluser)
	if err != nil {
		fb.SendErr(err, "获取记录失败", reterr)
		return
	}
	var clu Cluser
	err = tx.QueryRow("SELECT * from cluser where account = $1",
		cluser.Account).Scan(&clu.Account, &clu.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			fb.SendData(501, "账号不存在", reterr)
			return
		} else {
			fb.SendErr(err, "查询账号错误", reterr)
			return
		}
	}
	rows, err := tx.Query("SELECT chapter_num ,COUNT(*) FROM pra_record WHERE account = $1 group by chapter_num", cluser.Account)
	if err != nil {
		fb.SendErr(err, "获取错误", reterr)
		return
	}
	for rows.Next() {
		var cl Retprorec
		err := rows.Scan(&cl.Chapter_num, &cl.Chapter_rec)
		if err != nil {
			fb.SendErr(err, "获取错误", reterr)
			return
		}
		retrec = append(retrec, cl)
	}
	fb.SendData(200, "成功获取记录", retrec)
}

//修改密码
func ChangePassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fb := feedback.NewFeedBack(w)

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln("ioutil失败", err)

	}
	var userInfo UserInfo
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		logger.Errorln("读取数据失败", err)
		fb.SendData(501, "读取数据失败", "null")
		return
	}

	if userInfo.Account == "" {
		fb.SendData(502, "账号不能为空", "null")
		return
	}

	if userInfo.Password == "" {
		fb.SendData(503, "密码不能为空", "null")
		return
	}

	sqlState := "UPDATE cluser SET password = $1 WHERE account = $2;"
	stmt, err := Db.Prepare(sqlState)
	if err != nil {
		logger.Errorln("获取更新stmt失败", err)
		fb.SendStatus(502, "获取更新stmt失败")
		return
	}

	_, err = stmt.Exec(userInfo.Password, userInfo.Account)
	if err != nil {
		logger.Errorln("修改密码Exec失败")
		fb.SendStatus(502, "修改密码失败")
		return
	}
	fb.SendStatus(200, "修改密码成功")
}

func ClearRecord(w http.ResponseWriter, Account string) error {
	if Account == "" {
		return  errors.New("账号为空")
	}
	fb := feedback.NewFeedBack(w)
	sqlstmt := "DELETE FROM pra_record where account = $1"
	stmt, err := Db.Prepare(sqlstmt)
	if err != nil {
		logger.Errorln("获取stmt失败", err)
		fb.SendData(501, "清楚记录失败", nil)
		return err
	}

	_, err = stmt.Exec(Account)
	if err != nil {
		logger.Errorln("清楚记录失败")
		fb.SendData(500, "清楚记录失败", nil)
		return err
	}
	return nil
}

func Clear(w http.ResponseWriter, r *http.Request){
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, err := template.ParseFiles("login.html")
		if err != nil {
			logger.Errorln(err)
			fmt.Println(err)
		}
		t.Execute(w, nil)
		w.Header().Set("Content-type", "text/html")

	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		err := r.ParseForm()
		if err != nil {
			logger.Errorln(err)
			fmt.Println(err)
		}
		userName := template.HTMLEscapeString(r.Form.Get("username"))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("username:", userName) //输出到服务器端
		if userName == ""{
			logger.Errorln("账号为空")
			template.HTMLEscape(w,[]byte("不能为空"))
		}
		err = ClearRecord(w,userName)
		if err!=nil{
			logger.Errorln(err)
			template.HTMLEscape(w,[]byte("清除失败"))
			return
		}
		template.HTMLEscape(w, []byte(r.Form.Get("username")+"记录清楚成功")) //输出到客户端
	http.Redirect(w, r, "/", 302)
}
}


//GetPanic Rollback tx
func GetPanic(tx *sql.Tx) {
	if p := recover(); p != nil {
		tx.Rollback()

	}

}
