package practiceModel

type Record struct {
	Chapter_num  int `json:"chapter_num"`
	Quesiont_num int `json:"question_num"`
}

type Prarecord struct {
	Account string   `json:"account"`
	Record  []Record `json:"record"`
}

type Cluser struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type Retprorec struct {
	Chapter_num int `json:"chapter_num"`
	Chapter_rec int `json:"chapter_rec"`
}
