package test

import (
	"encoding/json"
	"fmt"
	http2 "go-js/src/http"
	"go-js/src/js"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/levigross/grequests"
)

func GetCookie() []*http.Cookie{
	var apiUrl = "https://h5api.m.taobao.com/h5/mtop.alimama.union.xt.biz.default.api.entry/1.0/?jsv=2.6.1&appKey=12574478&t=1616580991579&sign=662fcb1ab475bcb9a54ca2aa8167d961&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&data=%7B%22floorId%22%3A%2232447%22%2C%22variableMap%22%3A%22%7B%5C%22relationId%5C%22%3A%5C%222412826531%5C%22%2C%5C%22vegasCode%5C%22%3A%5C%22BA4AUQVU%5C%22%2C%5C%22lensId%5C%22%3A%5C%2221053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm%5C%22%2C%5C%22union_lens%5C%22%3A%5C%22lensId%3A21053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm%5C%22%2C%5C%22recoveryId%5C%22%3A%5C%22201_11.87.178.209_8868291_1616580989890%5C%22%2C%5C%22pvid%5C%22%3A%5C%22201_11.87.178.209_8868291_1616580989890%5C%22%7D%22%7D"
	var options  = &grequests.RequestOptions {
		RequestTimeout: time.Second * 2,
		DialTimeout: time.Second * 2,
		UseCookieJar: true,
	}
	req , err :=grequests.Get(apiUrl, options)

	if err != nil {
		panic(err)
		return nil
	}
	return req.RawResponse.Cookies()
}

func GetPvid() (pvid string) {
	var url = "https://uland.taobao.com/taolijin/edetail?vegasCode=FZTJ973I&type=qtz&union_lens=lensId%3A2127f8a3_08e1_1785ee4b626_da40%3Btraffic_flag%3Dlm&un=39b55ad2e92c55f5e7c4cca5375f99dd&share_crt_v=1&ut_sk=1.utdid_29215734_1616499750620.TaoPassword-Outside.taoketop&spm=a2159r.13376465.0.0&sp_tk=dW9EZ1haeFdOU04=&bxsign=tcdSzk0iWHkdYe3PzvGzb6dMIWgpLxCKBQY9jgQo9isiSE3_L82Jrxln-wMAndRx-3arKer7sMgBe_CPsGYdI3iAq-tKb8wSMM2AqCl3XLqr9Q/"
	var options = &grequests.RequestOptions {
		RequestTimeout: time.Second * 2,
		DialTimeout: time.Second * 2,
		UseCookieJar: true,
	}

	req , _:= grequests.Get(url, options)
	var bytes = req.Bytes()
	reg, _ := regexp.Compile(`window.pvid=\"(.+?)\";`)
	matched  := reg.Match(bytes)
	if matched {
		pvid = reg.FindStringSubmatch(req.String())[1]
	}

	return

}


type ApiParam struct {
	RelationId string
	VegasCode string
	UnionLens string
	Api string
	Jsv string
	AppKey string
	T int64
	Sign string
	Cookies []*http.Cookie
	Token string
	Data string
	URL string
}

type ApiData struct {
	FloorId string `json:"floorId"`
	VariableMap string `json:"variableMap"`
}

func (param *ApiParam) GetContent() {
	var options = &grequests.RequestOptions {
		RequestTimeout: time.Second * 2,
		DialTimeout: time.Second * 2,
		Cookies: param.Cookies,
		UseCookieJar: true,
	}

	//log.println(param.URL, param.Cookies)
	req , _ := grequests.Get(param.URL, options)
	//log.println(req.String())
}

func TestHttp(t *testing.T)  {
	var _url = "https://uland.taobao.com/taolijin/edetail?vegasCode=BA4AUQVU&type=qtz&union_lens=lensId%3A21053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm&relationId=2412826531&un=849771b72d5a07119ba585bab8a15b30&share_crt_v=1&ut_sk=1.utdid_28187926_1616550300551.TaoPassword-Outsidthis.taoketop&spm=a2159r.13376465.0.0&sp_tk=VHdpTlgwYUNtb2s=&bxsign=tcdvxdviWT681nhJ0adq_0OY68zqlWfOZODDSmHsMwwX19nvh1Fn1ocOyTEaLUxj7H7ujarNOKOMb_SNRlyKUdvKCI599dnv46ilEvGoOgLrZw/"

	parsedURL, err := url.Parse(_url)
	if err != nil {
		return
	}
	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)

	var param = &ApiParam{
		RelationId: parsedQuery.Get("relationId"),
		VegasCode: parsedQuery.Get("vegasCode"),
		UnionLens: parsedQuery.Get("union_lens"),
		Api:     "mtop.alimama.union.xt.biz.default.api.entry",
		Jsv:     "2.6.1",
		AppKey:  "12574478",
		T:       0,
		Sign:    "",
		Data:    "",
	}

	var pvid = GetPvid()
	var variableMap, _ = json.Marshal(map[string]string {
		"relationId": param.RelationId,
		"vegasCode": param.VegasCode,
		"lensId": js.GetLensId(param.UnionLens),
		"union_lens": param.UnionLens,
		"recoveryId": pvid,
		"pvid": pvid,
	})
	//{\"relationId\":\"2412826531\",\"vegasCode\":\"BA4AUQVU\",\"lensId\":\"21053dc7_0907_17861e80b13_10dd;traffic_flag=lm\",\"union_lens\":\"lensId:21053dc7_0907_17861e80b13_10dd;traffic_flag=lm\",\"recoveryId\":\"201_11.8.55.35_9554994_1616731785230\",\"pvid\":\"201_11.8.55.35_9554994_1616731785230\"}
	var data, _ = json.Marshal(ApiData {
		FloorId: "32477",
		VariableMap: string(variableMap),
	})

	param.Data = string(data)

	cookies := GetCookie()
	param.Cookies = cookies
	for _, cookie := range cookies {
		if cookie.Name == "_m_h5_tk" {
			//log.printf(cookie.Value)
			param.Token = js.GetToken(cookie.Value)

			//log.println("token: ", param.Token)
		}
	}

	//this.token + "&" + this.t + "&" + this.appKey + "&" + this.data
	param.T = time.Now().Unix() * 1000
	var needHexStr = fmt.Sprintf("%s&%d&%s&", param.Token, param.T, param.AppKey)
	param.Sign = js.JSHex(needHexStr, param.Data)
	//`https://h5api.m.taobao.com/h5/${this.api}/1.0/?jsv=${this.jsv}&appKey=${this.appKey}&t=${this.t}&sign=${sign}&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&dataType=jsonp&callback=mtopjsonp1&data=${this.data}`
	param.URL = fmt.Sprintf(`https://h5api.m.taobao.com/h5/%s/1.0/?jsv=%s&appKey=%s&t=%d&sign=%s&api=mtop.alimama.union.xt.biz.default.api.entry&v=1.0&AntiCreep=true&AntiFlood=true&type=jsonp&ecode=0&timeout=20000&dataType=jsonp&data=%s`, param.Api, param.Jsv, param.AppKey, param.T, param.Sign, url.PathEscape(param.Data))
	param.GetContent()
}

func TestProxy(t *testing.T) {
	res := http2.GetProxy()

	//log.println(res)
}

func TestJs(t *testing.T) {
	data := js.GetLensId("lensId:2127f8a3_08e1_1785ee4b626_da40;traffic_flag=lm")
	//log.println(data)
}