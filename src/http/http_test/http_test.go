package test

import (
	"go-js/src/http"
	"log"
	"testing"
)

var webUrl = "https://uland.taobao.com/taolijin/edetail?vegasCode=8PVDNWN4&type=qtz&union_lens=lensId%3A212c3a2a_088c_1785ef62ec7_d27d%3Btraffic_flag%3Dlm&un=39b55ad2e92c55f5e7c4cca5375f99dd&share_crt_v=1&ut_sk=1.utdid_29215734_1616500895640.TaoPassword-Outside.taoketop&spm=a2159r.13376465.0.0&sp_tk=T0N1UFhaeGp2OUs=&bxsign=tcd4JEGoq3BDnxViqeJREu8dn3yxUt5xa15kxKwYBNIrCmhSfYdl8_q9Gngoam7J8jXfLWqlYt1gRhL0PIn0ywk7lfpiLi5e0G5rEBFF4WlpIM/"

func TestProxy(t *testing.T) {

	content, err := http.GetContent(webUrl, http.ProxyParse(&http.ProxyData{
		IP:   "140.250.150.58",
		Port: 44124,
	}))
	log.Println(err, content)
}
