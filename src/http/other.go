package http

import (
	"fmt"
	"go-js/src/js"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/levigross/grequests"
)

type DataManage struct {
	PvId string

	sync.RWMutex
}

var Manage = &DataManage{
	PvId: "",
}

var PvId = ""

func GetCookie(param *ApiParam, proxies Proxies) ([]*http.Cookie, error) {
	//var apiUrl = "https://h5api.m.taobao.com/h5/mtop.alimama.union.xt.biz.default.api.entry/1.0/?jsv=2.6.1&appKey=12574478&t=1616580991579&sign=662fcb1ab475bcb9a54ca2aa8167d961&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&data=%7B%22floorId%22%3A%2232447%22%2C%22variableMap%22%3A%22%7B%5C%22relationId%5C%22%3A%5C%222412826531%5C%22%2C%5C%22vegasCode%5C%22%3A%5C%22BA4AUQVU%5C%22%2C%5C%22lensId%5C%22%3A%5C%2221053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm%5C%22%2C%5C%22union_lens%5C%22%3A%5C%22lensId%3A21053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm%5C%22%2C%5C%22recoveryId%5C%22%3A%5C%22201_11.87.178.209_8868291_1616580989890%5C%22%2C%5C%22pvid%5C%22%3A%5C%22201_11.87.178.209_8868291_1616580989890%5C%22%7D%22%7D"
	var options = &grequests.RequestOptions{
		RequestTimeout: time.Second * 2,
		DialTimeout:    time.Second * 2,
		UseCookieJar:   true,
		Proxies:        proxies,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57",
	}
	param.Token = "undefined"
	param.T = time.Now().Unix() * 1000
	var needHexStr = fmt.Sprintf("%s&%d&%s&", param.Token, param.T, param.AppKey)
	param.Sign = js.JSHex(needHexStr, param.Data)
	//`https://h5api.m.taobao.com/h5/${this.api}/1.0/?jsv=${this.jsv}&appKey=${this.appKey}&t=${this.t}&sign=${sign}&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&dataType=jsonp&callback=mtopjsonp1&data=${this.data}`
	param.URL = fmt.Sprintf(`https://h5api.m.taobao.com/h5/mtop.alimama.union.xt.biz.default.api.entry/1.0/?jsv=%s&appKey=%s&t=%d&sign=%s&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&dataType=jsonp&callback=mtopjsonp1&data=%s`, param.Jsv, param.AppKey, param.T, param.Sign, param.Data)

	req, err := grequests.Get(param.URL, options)

	if err != nil {
		return nil, err
	}

	log.Println("first api :", req.String())
	return req.RawResponse.Cookies(), nil
}

func GetPvid(url string, proxies Proxies) (pvid string, err error) {
	// var url = "https://uland.taobao.com/taolijin/edetail?vegasCode=FZTJ973I&type=qtz&union_lens=lensId%3A2127f8a3_08e1_1785ee4b626_da40%3Btraffic_flag%3Dlm&un=39b55ad2e92c55f5e7c4cca5375f99dd&share_crt_v=1&ut_sk=1.utdid_29215734_1616499750620.TaoPassword-Outside.taoketop&spm=a2159r.13376465.0.0&sp_tk=dW9EZ1haeFdOU04=&bxsign=tcdSzk0iWHkdYe3PzvGzb6dMIWgpLxCKBQY9jgQo9isiSE3_L82Jrxln-wMAndRx-3arKer7sMgBe_CPsGYdI3iAq-tKb8wSMM2AqCl3XLqr9Q/"

	var options = &grequests.RequestOptions{
		RequestTimeout: time.Second * 2,
		DialTimeout:    time.Second * 2,
		Proxies:        proxies,
		UseCookieJar:   true,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57",
	}

	req, err := grequests.Get(url, options)

	if err != nil {
		return
	}

	//log.Println(req.String())

	var bytes = req.Bytes()
	reg, err := regexp.Compile(`window.pvid=\"(.+?)\";`)
	if err != nil {
		return
	}
	matched := reg.Match(bytes)
	if matched {

		pvid = reg.FindStringSubmatch(req.String())[1]
	}

	log.Println("pvid Cookie: ", req.RawResponse.Cookies())
	grequests.Get("https://px.effirst.com/api/v1/jconfig?wpk-header=app%3Dalimama_lego2%26tm%3D1616817602%26ud%3D38e0f412-8833-4079-0706-560c7018124c%26sver%3D0.7.7%26sign%3Dc41e43c828c16c16a6eb1c9c1e68e8ce", options)

	return
}
