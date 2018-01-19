package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
    "net/http/cookiejar"
    "io/ioutil"
    "regexp"
    "os"
)

type BaseClient struct {
	Url    string
	Data   url.Values
	Client http.Client
    Cookie  *cookiejar.Jar
}

func (this *BaseClient) Login() {
	request, err := http.NewRequest("POST", this.Url, strings.NewReader(this.Data.Encode()))
	if err != nil {
		fmt.Println("new request error: ", err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Set("Referer", "http://poj.org/problemlist")
	_, Clienterr := this.Client.Do(request)
	if Clienterr != nil {
		fmt.Println("client.Do error: ", err.Error())
		return
	}
    /*
    fmt.Println("resp", resp)
    fmt.Println("================================")
    for _, cookie := range this.Cookie.Cookies(request.URL) {
        fmt.Printf("== %v", cookie)
    }
    */
}


func (this *BaseClient)GetRunId(problem_id string) string{
    query_str := "http://acm.nyist.net/JudgeOnline/status.php?do=search&pid=" + problem_id + "&userid=Lost_in_wine&language=0&result=Accepted"
    req, _ := http.NewRequest("GET", query_str, nil)
    resp, _ := this.Client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    reg := regexp.MustCompile(`(?U)<form name="acoin(?P<run_id>\d+)" method="post"`)
    mp := reg.FindAllStringSubmatch(string(body), -1)
    return mp[0][1]
}

func (this *BaseClient)GetProblem(run_id string) string {
    req, _ := http.NewRequest("GET", "http://acm.nyist.net/JudgeOnline/code.php?runid=" + run_id, nil)
    resp, _ := this.Client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    reg := regexp.MustCompile(`(?U)<pre class="brush: cpp; ">(?P<code>[\s\S]*) *</pre`)
    mp := reg.FindAllStringSubmatch(string(body), -1)
    ret := strings.Replace(mp[0][1], "&#34;", "\"", -1)  
    ret = strings.Replace(ret, "&#39;", "'", -1)  
    ret = strings.Replace(ret, "&#38;", "&", -1)  
    ret = strings.Replace(ret, "&#60;", "<", -1)  
    ret = strings.Replace(ret, "&#62;", ">", -1)  

    ret = strings.Replace(ret, "&quot;", "\"", -1)  
    ret = strings.Replace(ret, "&apos;", "'", -1)  
    ret = strings.Replace(ret, "&amp;", "&", -1)  
    ret = strings.Replace(ret, "&lt;", "<", -1)  
    ret = strings.Replace(ret, "&gt;", ">", -1)  
    return ret
}

func main() {
	poj := BaseClient{
        Url:    "http://acm.nyist.net/JudgeOnline/dologin.php?url=http%3A%2F%2Facm.nyist.net%2FJudgeOnline%2Fproblemset.php",
        Cookie: nil,
		Client: http.Client{
            CheckRedirect: nil,
        },
		Data:   url.Values{},
	}

    poj.Cookie, _ = cookiejar.New(nil)
    poj.Client.Jar = poj.Cookie

	poj.Data.Add("userid", "your_id")
	poj.Data.Add("password", "your_password")
	poj.Data.Add("btn_submit", "登录")
	poj.Data.Add("url", "/")

	poj.Login()

    problems := []string{
        "1", "2", "3"
    }
    for _, v := range problems {
        run_id := poj.GetRunId(v)
        fmt.Println("problem_id = ", v, "run_id = ", run_id)
        body := poj.GetProblem(run_id)
        file, err := os.OpenFile("nyoj/" + v + ".cc", os.O_RDWR|os.O_CREATE, 0666);
        if (err != nil) {
            fmt.Println(err);
            continue
        }
        file.WriteString(body)
        file.Close()
    }
}
