package robot

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
    req, _ := http.NewRequest("GET", "http://poj.org/status?problem_id=" + problem_id + "&user_id=Lost_in_wine&result=0&language=", nil)
    resp, _ := this.Client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    reg := regexp.MustCompile(`(?U)<tr align=center><td>(?P<run_id>\d+)</td>[\s\S]*(?P<problem_id>\d+)>`)
    mp := reg.FindAllStringSubmatch(string(body), -1)
    return mp[0][1]
}

func (this *BaseClient)GetProblem(run_id string) string {
    req, _ := http.NewRequest("GET", "http://poj.org/showsource?solution_id=" + run_id, nil)
    resp, _ := this.Client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    reg := regexp.MustCompile(`(?U)monospace">(?P<code>[\s\S]*)</pre`)
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
		Url:    "http://poj.org/login",
        Cookie: nil,
		Client: http.Client{
            CheckRedirect: nil,
        },
		Data:   url.Values{},
	}

    poj.Cookie, _ = cookiejar.New(nil)
    poj.Client.Jar = poj.Cookie

	poj.Data.Add("user_id1", "your_user_id")
	poj.Data.Add("password1", "your_passwd")
	poj.Data.Add("B1", "login")
	poj.Data.Add("url", "/")

	poj.Login()
    //problems := []string{ "1000", "1006", "1061", "1149", "1182", "1201", "1273", "1308", "1321", "1330", "1364", "1509", "1611", "1631", "1743", "1816", "1887", "1988", "2015", "2084", "2155", "2236", "2251", "2255", "2312", "2406", "2420", "2478", "2488", "2492", "2524", "2533", "2549", "2559", "2676", "2752", "2773", "2774", "2778", "2785", "2796", "2891", "3009", "3061", "3069", "3070", "3083", "3159", "3253", "3278", "3437", "3468", "3617", "3660", "3974", "3984", }
    problems := []string{"3253", "3278", "3437", "3468", "3617", "3660", "3974", "3984", }

    for _, v := range problems {
        run_id := poj.GetRunId(v)
        fmt.Println("problem_id = ", v, "run_id = ", run_id)
        body := poj.GetProblem(run_id)
        file, err := os.OpenFile("poj/" + v + ".cc", os.O_RDWR|os.O_CREATE, 0666);
        if (err != nil) {
            fmt.Println(err);
            continue
        }
        file.WriteString(body)
        file.Close()
    }
}
