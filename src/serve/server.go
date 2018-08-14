package serve

import (
	"bytes"
	"crypto/md5"
	. "db"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"session"
	"strconv"
	"time"
)

type Cluser struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type Td struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

//SayhelloName for test http
func SayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello astaxie!")
}

//Login test
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

	rd := Td{Msg: "Request data", Code: 200, Data: clusers}
	jsonuse, _ := json.Marshal(rd)
	logger.Info("Get:" + string(jsonuse))
	fmt.Fprintln(w, string(jsonuse))
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
		rd := Td{Msg: "Post data", Code: 200, Data: at}
		resp, err := json.Marshal(rd)
		if err != nil {
			logger.Errorln(err)
		}
		w.Write(resp)
	}
}
