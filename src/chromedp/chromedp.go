package chromedp

import (
	"context"
	"flag"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func Exec(webUrl string) (url string, Cookies []*http.Cookie) {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	var allocCtx context.Context
	var cancelActxt context.CancelFunc
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserDataDir(dir),
	)
	var devToolWsUrl string

	flag.StringVar(&devToolWsUrl, "devtools-ws-url", "ws://localhost:9222/devtools/browser/e33c437d-a859-43e9-b407-13761e143eb7", "DevTools Websocket URL")
	flag.Parse()

	if runtime.GOOS == "linux" {
		allocCtx, cancelActxt = chromedp.NewRemoteAllocator(context.Background(), devToolWsUrl)

	} else {
		allocCtx, cancelActxt = chromedp.NewExecAllocator(context.Background(), opts...)
	}

	defer cancelActxt()

	//allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	//defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 10*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	// listen network event
	listenForNetworkEvent(taskCtx, &url)

	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(`https://uland.taobao.com/taolijin/edetail?vegasCode=8PVDNWN4&type=qtz&union_lens=lensId%3A212c3a2a_088c_1785ef62ec7_d27d%3Btraffic_flag%3Dlm&un=39b55ad2e92c55f5e7c4cca5375f99dd&share_crt_v=1&ut_sk=1.utdid_29215734_1616500895640.TaoPassword-Outside.taoketop&spm=a2159r.13376465.0.0&sp_tk=T0N1UFhaeGp2OUs=&bxsign=tcd4JEGoq3BDnxViqeJREu8dn3yxUt5xa15kxKwYBNIrCmhSfYdl8_q9Gngoam7J8jXfLWqlYt1gRhL0PIn0ywk7lfpiLi5e0G5rEBFF4WlpIM/`),
		//chromedp.Navigate(`https://xrdev.top/api/test.php`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 获取cookie
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			var c string
			for _, v := range cookies {
				c = c + v.Name + "=" + v.Value + ";"
				Cookies = append(Cookies, &http.Cookie{
					Name:  v.Name,
					Value: v.Value,
				})
			}
			log.Println(c)
			return nil
		}),
	)
	return
}

//监听
func listenForNetworkEvent(ctx context.Context, url *string) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				// log.Printf("received headers: %s", resp.Headers)
				if strings.HasPrefix(resp.URL, "https://h5api.m.taobao.com/h5/mtop.alimama.union.xt.biz.default.api.entry") {
					*url = resp.URL
					return
				}
			}
		}
	})
}
