package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	Arther            int64 = 327835980
	RDMember                = []int64{Arther}
	AllowChatId             = []int64{}
	TaiwanLocation, _       = time.LoadLocation("Asia/Taipei")
	Power                   = true
)

type LocationConfig struct {
	Url string `json:"url"`
}

type WeatherChecker struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type Config struct {
	LocationConfig *LocationConfig `json:"location_config"`
	WeatherChecker *WeatherChecker `json:"weather_checker"`
	TelegramToken  string          `json:"telegram_token"`
}

// 設定檔
var ServerConfig = &Config{

	LocationConfig: &LocationConfig{
		Url: "https://api.nlsc.gov.tw/",
	},
	WeatherChecker: &WeatherChecker{
		Token: "",
		Url:   "https://opendata.cwb.gov.tw/api",
	},
	TelegramToken: "",
}

func init() {
	configPath := flag.String("conf", "", "set setting config path service need")
	flag.Parse()
	if *configPath != "" {
		fmt.Println("讀取設定檔")
		conf, err := ReadConfig(*configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
			return
		}
		ServerConfig = conf
	}

	AllowChatId = append(AllowChatId, RDMember...)
}

func ReadConfig(path string) (config *Config, err error) {
	ConfigContent, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		err = readErr
		return
	}
	config = &Config{}
	if err = json.Unmarshal(ConfigContent, config); err != nil {
		return
	}
	ServerConfig = config
	return
}
