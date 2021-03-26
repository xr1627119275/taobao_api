package main

import (
	"fmt"
	"go-js/src/http"
	"os"
	"strings"
)

func main() {
	var args = os.Args
	if len(args) != 2 {
		os.Exit(-1)
	}

	Proxys := http.GetAvailProxys()

	for _, proxy := range Proxys {
		content, err := http.GetContent(strings.Trim(args[1], "'"), http.ProxyParse(proxy))
		if err != nil {
			continue
		}
		content = strings.Trim(strings.Trim(content, " callback("), ")")
		fmt.Println(content)
		os.Exit(0)
	}
	os.Exit(-1)
}
