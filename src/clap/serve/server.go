package serve

import (
	"bytes"
	. "clap/db"
	"clap/feedback"
	"clap/session"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"clap/logger"
	"net/http"
	"strconv"
	"time"
)

//SayhelloName for test http
func SayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello astaxie!")
}

//login test
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	sess := session.GlobalSessions.SessionStart(w, r)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, err := template.ParseFiles("login.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, token)
		w.Header().Set("Content-type", "text/html")

	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
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

//GetPanic Rollback tx
func GetPanic(tx *sql.Tx) {
	if p := recover(); p != nil {
		tx.Rollback()

	}

}
