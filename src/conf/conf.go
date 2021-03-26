package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 配置参数管理
type mConfig struct {
	ProxyURL string // 代理提取地址
}

// Init 读取注册中心配置，订阅配置服务配置更新。
func (s *mConfig) Init() {
	_, filename := filepath.Split(os.Args[0])
	index := strings.IndexByte(filename, '.')
	if index >= 0 {
		filename = filename[0:index]
	}
	configFile := fmt.Sprintf("./configs/%v.conf", filename)
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./configs/", os.ModePerm)
			if data, err := json.MarshalIndent(s, "", "  "); err == nil {
				if err = ioutil.WriteFile(configFile, data, os.ModePerm); err != nil {
					//log.Println("create config file fail.", configFile, err)
				}
			}
		}
	} else {
		err = json.Unmarshal(configData, &s)
		if err != nil {
			//log.Println("error config!!!!", configFile, err)
		}
	}
}
func (s *mConfig) SaveConf() {
	_, filename := filepath.Split(os.Args[0])
	index := strings.IndexByte(filename, '.')
	if index >= 0 {
		filename = filename[0:index]
	}
	configFile := fmt.Sprintf("./configs/%v.conf", filename)
	_, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./configs/", os.ModePerm)
		}
	}
	if data, err := json.MarshalIndent(s, "", "  "); err == nil {
		if err = ioutil.WriteFile(configFile, data, os.ModePerm); err != nil {

		}
	}
}

func newConfig() *mConfig {
	var config = &mConfig{
		ProxyURL: "http://a.ipjldl.com/getapi?packid=2&unkey=&tid=&qty=10&time=2&port=1&format=json&ss=5&css=&pro=&city=&dt=1&usertype=17",
	}
	config.Init()
	return config
}

// Config 全局配置对象
var Config = newConfig()
