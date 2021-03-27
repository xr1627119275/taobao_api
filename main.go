package main

import (
	"encoding/json"
	"fmt"
	"go-js/src/conf"
	"go-js/src/http"
	"log"
	"os"
	"strings"
)

func toRequest(url string, proxy *http.ProxyData) (err error) {
	content, err := http.GetContent(url, http.ProxyParse(proxy))

	if err != nil {
		return
	}
	content = strings.Trim(strings.Trim(content, " callback("), ")")
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
func main() {
	var args = os.Args
	var tabaoUrl = "https://uland.taobao.com/taolijin/edetail?vegasCode=8PVDNWN4&type=qtz&union_lens=lensId%3A212c3a2a_088c_1785ef62ec7_d27d%3Btraffic_flag%3Dlm&un=39b55ad2e92c55f5e7c4cca5375f99dd&share_crt_v=1&ut_sk=1.utdid_29215734_1616500895640.TaoPassword-Outside.taoketop&spm=a2159r.13376465.0.0&sp_tk=T0N1UFhaeGp2OUs=&bxsign=tcd4JEGoq3BDnxViqeJREu8dn3yxUt5xa15kxKwYBNIrCmhSfYdl8_q9Gngoam7J8jXfLWqlYt1gRhL0PIn0ywk7lfpiLi5e0G5rEBFF4WlpIM/"

	if len(args) != 2 && os.Getenv("env") != "dev" {
		os.Exit(-1)
		return
	}
	if os.Getenv("env") != "dev" {
		tabaoUrl = strings.Trim(os.Args[1], "'")
	}
	//http.PvId , _  = http.GetPvid(tabaoUrl, nil)

	lastProxy := ReadLastProxy()

	if lastProxy != nil {
		err := toRequest(tabaoUrl, lastProxy)
		if err == nil {
			os.Exit(0)
			return
		}
		SetLastProxyIP(nil)
	}

	//Proxys := http.GetAvailProxys()
	Proxys := []*http.ProxyData{nil}

	log.Println(Proxys)

	//wg := &sync.WaitGroup{}
	for _, proxy := range Proxys {
		log.Println(proxy)
		//log.Println(http.GetPvid(tabaoUrl , http.ProxyParse(proxy)))
		//wg.Add(1)
		//go func(wg *sync.WaitGroup, proxy *http.ProxyData) {
		//	err := toRequest(tabaoUrl, proxy)
		//	if err == nil {
		//		log.Println("get:" , proxy)
		//		os.Exit(0)
		//	}
		//	wg.Done()
		//}(wg, proxy)
		err := toRequest(tabaoUrl, proxy)
		if err == nil {
			log.Println("get:", proxy)
			os.Exit(0)
			return
		}
	}
	//wg.Wait()
	SetLastProxyIP(nil)

	if toRequest(tabaoUrl, nil) == nil {
		log.Println("get_default")
		os.Exit(0)
	}
	os.Exit(-1)

}
