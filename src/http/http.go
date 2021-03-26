package http

import (
	"encoding/json"
	"fmt"
	"go-js/src/js"
	"net/http"
	"net/url"
	"time"

	"github.com/levigross/grequests"
)


type Proxies map[string] *url.URL


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

func (param *ApiParam) GetContent(proxies Proxies) (string, error){
	var options = &grequests.RequestOptions {
		RequestTimeout: time.Second * 2,
		DialTimeout: time.Second * 2,
		Cookies: param.Cookies,
		Proxies: proxies,
		UseCookieJar: true,
	}


	req , err := grequests.Get(param.URL, options)

	if err != nil {
		return "", err
	}
	return req.String(), nil
}

func GetContent(_url string, proxies Proxies) (string, error) {
	//var _url = "https://uland.taobao.com/taolijin/edetail?vegasCode=BA4AUQVU&type=qtz&union_lens=lensId%3A21053dc7_0907_17861e80b13_10dd%3Btraffic_flag%3Dlm&relationId=2412826531&un=849771b72d5a07119ba585bab8a15b30&share_crt_v=1&ut_sk=1.utdid_28187926_1616550300551.TaoPassword-Outsidthis.taoketop&spm=a2159r.13376465.0.0&sp_tk=VHdpTlgwYUNtb2s=&bxsign=tcdvxdviWT681nhJ0adq_0OY68zqlWfOZODDSmHsMwwX19nvh1Fn1ocOyTEaLUxj7H7ujarNOKOMb_SNRlyKUdvKCI599dnv46ilEvGoOgLrZw/"

	parsedURL, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)

	if err != nil {
		return "", err
	}

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

	pvid, err := GetPvid(_url, proxies)
	if err != nil {
		return "", err
	}
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

	cookies ,err := GetCookie(proxies)
	if err != nil {
		return "", err
	}

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
	return param.GetContent(proxies)
}

