package chromedp

import (
	"context"
	"flag"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go-js/src/conf"
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

	flag.StringVar(&devToolWsUrl, "devtools-ws-url", conf.Config.WsUrl, "DevTools Websocket URL")
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
		chromedp.Navigate(webUrl),
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
				if strings.HasPrefix(resp.URL, "https://api.cmsv5.iyunzk.com/apis") {
					*url = resp.URL
					return
				}
			}
		}
	})
}
