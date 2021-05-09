package main

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"go-js/src/chromedp"
	"go-js/src/conf"
	"go-js/src/http"
	"log"
	http2 "net/http"
	"os"
	"strings"
)

const url2 = "http://tkapi.apptimes.cn/tao-password/parse?appkey=7lfrrdqq&appsecret=5r2uetrr6qi6odpx&password="

func toRequest(url string, proxy *http.ProxyData) (err error) {
	req, err := grequests.Get(url, &grequests.RequestOptions{Proxies: http.ProxyParse(proxy)})

	if err != nil {
		return
	}

	var res = &ResponseDefault{}

	err = req.JSON(res)
	if err != nil {
		return
	}
	if res.Code == 200 {
		SetLastProxyIP(proxy)
	}

	content, err := grequests.Get(url2+res.Data, nil)
	if err != nil {
		return err
	}

	fmt.Println(content.String())

	return
}

func toRequestDefault(url string, Cookies []*http2.Cookie, proxy *http.ProxyData) (err error) {
	content, err := http.UrlCookie2Content(url, Cookies, http.ProxyParse(proxy))

	if err != nil {
		return
	}
	//content = js.JSCbContent(content)
	//content = strings.Trim(content, " mtopjsonp1(")
	fmt.Println(content)
	if len(content) > 200 {
		SetLastProxyIP(proxy)
	}
	return
}

func SetLastProxyIP(proxy *http.ProxyData) {
	if proxy != nil {
		bytes, _ := json.Marshal(proxy)
		conf.Config.AvailIP = string(bytes)
	} else {
		conf.Config.AvailIP = ""
	}
	conf.Config.SaveConf()
}

func ReadLastProxy() (proxy *http.ProxyData) {
	proxy = &http.ProxyData{}
	AvailIP := conf.Config.AvailIP

	if len(AvailIP) == 0 {
		proxy = nil
		return
	}

	err := json.Unmarshal([]byte(AvailIP), proxy)
	if err != nil {
		proxy = nil
		return
	}
	return
}

type ResponseDefault struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func main() {
	var args = os.Args
	var tabaoUrl = "https://yun043.kuaizhan.com/t/BLXgz9"

	if len(args) != 2 && os.Getenv("env") != "dev" {
		os.Exit(-1)
		return
	}
	if os.Getenv("env") != "dev" {
		tabaoUrl = strings.Trim(os.Args[1], "'")
	}
	//http.PvId , _  = http.GetPvid(tabaoUrl, nil)

	apiUrl, _ := chromedp.Exec(tabaoUrl)

	if len(strings.Trim(apiUrl, " ")) == 0 {
		os.Exit(-1)
		return
	}
	log.Println(apiUrl)
	log.Println("OpenProxy: ", conf.Config.OpenProxy)

	if conf.Config.OpenProxy != "1" {
		err := toRequest(apiUrl, nil)
		if err == nil {
			os.Exit(0)
		} else {
			os.Exit(-2)
		}
		return
	}
	//log.Println(cookies)
	lastProxy := ReadLastProxy()

	if lastProxy != nil {
		//err := toRequest(tabaoUrl, lastProxy)
		err := toRequest(apiUrl, lastProxy)
		if err == nil {
			os.Exit(0)
			return
		}
		SetLastProxyIP(nil)
	}

	Proxys := http.GetAvailProxys()
	log.Println(Proxys)

	for _, proxy := range Proxys {
		log.Println(proxy)

		err := toRequest(apiUrl, &proxy)
		if err == nil {
			//log.Println("get:", proxy)

			SetLastProxyIP(&proxy)
			os.Exit(0)
			return
		}
	}
	//wg.Wait()
	SetLastProxyIP(nil)

	if toRequest(apiUrl, nil) == nil {
		log.Println("get_default")
		os.Exit(0)
	}
	os.Exit(-1)

}
