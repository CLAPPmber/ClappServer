package TreadPit

import (
	"testing"
	"fmt"
	"time"
	"net/http"
	"log"
	"github.com/valyala/fasthttp"
	"net/url"
	"os"
)

type TestDbData func(*testing.T,*TestCase)  // 声明了一个函数类型

type TestCase struct {
	Protocol 	string //Url协议
	Host 		string //请求域名+端口
	Api 		string //调用的api
	Method 		string //请求方法
	Describe 	string
	ExpectResp 	string //期望返回值
	ReqParsMsp 	map[string]string //请求参数
	ReqParsJson string //请求参数
	InsertMockData TestDbData //往数据库里添加测试数据
	ClearMockData TestDbData //清楚数据库中的测试数据
	UserType 	string //
}

const srv = "127.0.0.1:9090"
const reTryTimes = 10

func TestMain(m *testing.M) {
	serverOk := make(chan int)
	url := fmt.Sprintf("http://%s/gsql", srv)
	go func() {
		reTry := 0
		for {
			if reTry >= reTryTimes {
				panic(m)
			}
			time.Sleep(2 * time.Second)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(fmt.Sprintf("retry...%d", reTry))
				reTry++
				continue
			}
			if resp.StatusCode != 200 {
				log.Println(fmt.Sprintf("retry...%d", reTry))
				reTry++
				continue
			}
			serverOk <- 1
		}
	}()
	<-serverOk
	log.Println("Server Ready")
	InitTestDb()

	ExitCode := m.Run()

	ClearTestData()

	fmt.Println("test complete")
	os.Exit(ExitCode)

}

func ReqTest(t *testing.T,testCase *TestCase,client *fasthttp.Client){

	if testCase.Protocol == ""{
		testCase.Protocol = "HTTP"
	}

	if testCase.Host == "" {
		testCase.Host = srv
	}

	if testCase.InsertMockData !=nil {
		testCase.InsertMockData(t,testCase)
	}

	if testCase.ClearMockData !=nil {
		defer testCase.ClearMockData(t,testCase)
	}


	var err error

	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}
	ParseUrl,err := url.Parse(testCase.Protocol+"://"+testCase.Host+testCase.Api)
	if err!=nil{
		t.Errorf("%s:err:%s",testCase.Describe,err.Error())
		return
	}

	switch testCase.Method {
	case "GET":

		url := ParseUrl.String()
		req.SetRequestURI(url)
		req.Header.SetMethod(http.MethodGet)
		err = client.Do(req,resp)
		break
	case "POST":
		//PostJson数据
		req.Header.Set("Content-Type","application/json")
		body := testCase.ReqParsJson

		if body == "" {
			t.Errorf("%s:err,%s",testCase.Describe,"ReqParsJson empty")
			return
		}

		req.SetRequestURI(ParseUrl.String())
		req.Header.SetMethod(http.MethodPost)
		req.SetBodyString(body)
		err = client.Do(req,resp)

		break
	default:
		t.Error(testCase.Describe + "unknown request method: " + testCase.Method)
		return
		break
	}
	if err!=nil {
		t.Fatal()
	}

	if testCase.ExpectResp !="" && testCase.ExpectResp != string(resp.Body()){
		t.Fatalf("%v fail,expect %v,but get %v",testCase.Describe,testCase.ExpectResp,string(resp.Body()))
	}
}

func TestPracticeModel(t *testing.T){
	client := &fasthttp.Client{}
	var testCase *TestCase

	//测试test
	testCase = &TestCase{
		Api:"/gsql",
		Method:http.MethodGet,
		Describe:"test get gsql",
		ExpectResp:`{"code":200,"msg":"SqlGet!","data":{"account":"usernew3","password":"123456"}}`,
	}
	ReqTest(t,testCase,client)

	//测试获取用户做题记录
	testCase = &TestCase{
		Api:"/getallrec",
		Method:http.MethodPost,
		Describe:"test get practiceRecord",
		InsertMockData:InsertMockPracticeData,
		ClearMockData:nil,
		ExpectResp: GetPracticeExpectResp,
		ReqParsJson:GetPracticeRecReqPars,
	}
	ReqTest(t,testCase,client)

	//测试用户插入新纪录
	testCase = &TestCase{
		Api:"/prarecord",
		Method:http.MethodPost,
		Describe:"test commit new practiceRecord",
		InsertMockData:nil,
		ClearMockData:nil,
		ReqParsJson:CommitRecReq,
		ExpectResp:SuccessCommitRecExpectResp,
	}
	ReqTest(t,testCase,client)


	//测试提交重复记录
	testCase = &TestCase{
		Api:"/prarecord",
		Method:http.MethodPost,
		Describe:"test commit a exist practiceRecord ",
		ReqParsJson:CommitRecReq,
		ExpectResp:FailCommitRecExpectResp,
		InsertMockData:nil,
		ClearMockData:ClearMockPracticeData,
	}
	ReqTest(t,testCase,client)

}