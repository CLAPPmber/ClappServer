package APPtable

type LoginData struct {
	Account  string `json:account`
	Password string `json:"password"`
}

type FeedBack struct {
	Msg  string    `json:msg`
	Code int       `json:code`
	Data LoginData `json:data`
}
