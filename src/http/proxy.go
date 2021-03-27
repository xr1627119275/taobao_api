package http

import (
	"fmt"
	"go-js/src/conf"
	"net/url"
	"time"

	"github.com/levigross/grequests"
)

type ProxyResponse struct {
	Code    int          `json:"code"`
	Success string       `json:"success"`
	Msg     string       `json:"msg"`
	Data    []*ProxyData `json:"data"`
}

type ProxyData struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
}

//const proxy_url_1 = "http://a.ipjldl.com/getapi?packid=2&unkey=&tid=&qty=1&time=2&port=1&format=txt&ss=1&css=&pro=&city=&dt=1&usertype=17"
var proxy_url = conf.Config.ProxyURL

func GetProxy() *ProxyResponse {
	req, _ := grequests.Get(proxy_url, nil)
	var proxyRes = &ProxyResponse{}
	_ = req.JSON(proxyRes)
	return proxyRes
}

func (proxy *ProxyData) Check() bool {
	proxies := ProxyParse(proxy)
	req, err := grequests.Get("https://uland.taobao.com/", &grequests.RequestOptions{
		RequestTimeout:      time.Second * 3,
		DialTimeout:         time.Second * 3,
		TLSHandshakeTimeout: time.Second * 3,
		Proxies:             proxies,
	})
	if err != nil {
		return false
	}

	if req.StatusCode == 200 {
		return true
	}
	return false
}

func GetAvailProxys() (proxys []*ProxyData) {
	res := GetProxy()
	//wg := &sync.WaitGroup{}
	for i := 0; i < len(res.Data); i++ {
		item := res.Data[i]

		//proxys = append(proxys, item)
		//wg.Add(1)
		//go func(wg *sync.WaitGroup) {
		//	check := item.Check()
		//	if check {
		proxys = append(proxys, item)
		//	}
		//	wg.Done()
		//}(wg)
	}
	//wg.Wait()
	return
}

func ProxyParse(proxy *ProxyData) (proxies Proxies) {
	if proxy != nil {
		URL, _ := url.Parse(fmt.Sprintf("socks5://%s:%d", proxy.IP, proxy.Port))
		proxies = Proxies{
			"http":  URL,
			"https": URL,
		}
	}

	return
}
