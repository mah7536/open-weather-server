package checker

import (
	"alarm-system/config"
	"alarm-system/db/cache"
	"alarm-system/telegram"
	tlib "alarm-system/telegram/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"188.166.240.198/GAIUS/lib/errorCode"
	"188.166.240.198/GAIUS/lib/errorhandler"
	"188.166.240.198/GAIUS/lib/httpClient"
	"188.166.240.198/GAIUS/lib/logger"
)

const (
	UrlForGet36HourWeather = "/v1/rest/datastore/F-C0032-001"
)

type EachTimeInfo struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Parameter struct {
		Name  string `json:"parameterName"`
		Value string `json:"parameterValue"`
		Unit  string `json:"parameterUnit"`
	} `json:"parameter"`
}

type WeatherElement struct {
	ElementName string          `json:"elementName"`
	Time        []*EachTimeInfo `json:"time"`
}

type Location struct {
	LocationName    string            `json:"locationName"`
	WeatherElements []*WeatherElement `json:"weatherElement"`
}

type Record struct {
	Description string      `json:"datasetDescription"`
	Location    []*Location `json:"location"`
}

type WeatherResponse struct {
	Success string  `json:"success"`
	Records *Record `json:"records"`
}

type WeatherChecker struct {
	*config.WeatherChecker
	Client   *httpClient.Client `json:"-"`
	LastCode int
	IsHealth bool
	Duration int
	List     map[string]*cache.WeatherData
}

func NewWeatherChecker(config *config.WeatherChecker) *WeatherChecker {
	return &WeatherChecker{
		WeatherChecker: config,
		Client:         httpClient.NewClient(),
		LastCode:       errorCode.Success,
	}
}

func NewWeatherRequest(url string, token string) (code int, req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		code = errorCode.Error
		return
	}
	q := req.URL.Query()
	q.Add("Authorization", token)

	req.URL.RawQuery = q.Encode()
	return
}

func (checker *WeatherChecker) Run() {
	checkWeatherStatus := time.NewTicker(1 * time.Hour)
	name := "撈取天氣"
	for {
		code, weatherData, err := checker.Job()

		checker.LastCode = code

		if code != errorCode.Success {
			ms := tlib.NewCommonMessage(config.Arther, tlib.TypeWarn, name, "查無天氣資料(1)")
			telegram.SendMs(ms)
			time.Sleep(1 * time.Minute)
			continue
		}

		if checker.LastCode != errorCode.Success && code == errorCode.Success {
			ms := tlib.NewCommonMessage(config.Arther, tlib.TypeInfo, name, "恢復正常")
			telegram.SendMs(ms)
		}

		if weatherData.Records == nil {
			ms := tlib.NewCommonMessage(config.Arther, tlib.TypeWarn, name, "查無天氣資料(2)")
			telegram.SendMs(ms)
			time.Sleep(1 * time.Minute)
			continue
		}

		if weatherData.Records.Location == nil {
			ms := tlib.NewCommonMessage(config.Arther, tlib.TypeWarn, name, "查無天氣資料(3)")
			telegram.SendMs(ms)
			time.Sleep(1 * time.Minute)
			continue
		}

		if len(weatherData.Records.Location) == 0 {
			ms := tlib.NewCommonMessage(config.Arther, tlib.TypeWarn, name, "查無天氣資料(4)")
			telegram.SendMs(ms)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, eachData := range weatherData.Records.Location {
			code, data, err := DecodeWeatherStatus(eachData)
			if code != errorCode.Success {
				logger.Error(err)
				continue
			}

			cache.Server.SetWeatherData <- data
		}

		if code != errorCode.Success {
			res := errorhandler.NewResponse(code)
			res.SetExtra(err)
			logger.Error(res)
			if checker.LastCode == errorCode.Success {
				ms := tlib.NewCommonMessage(config.Arther, tlib.TypeDanger, name, fmt.Sprint(err.Error()))
				telegram.SendMs(ms)
			}
		}

		<-checkWeatherStatus.C
	}
}

func (checker *WeatherChecker) Job() (code int, data *WeatherResponse, err error) {

	code, req, err := NewWeatherRequest(checker.Url+UrlForGet36HourWeather, checker.Token)
	if code != errorCode.Success {
		code = errorCode.RequestCreateError
		checker.IsHealth = false
		return
	}

	start := time.Now()
	code, res, err := checker.Client.Send(req)
	if code != errorCode.Success {
		checker.IsHealth = false
		return
	}

	data = &WeatherResponse{}
	err = json.Unmarshal(res, data)
	if err != nil {
		code = errorCode.DecodeJsonError
		return
	}

	if data.Success != "true" {
		code = errorCode.Error
		return
	}
	checker.IsHealth = true
	checker.Duration = int(time.Since(start).Seconds())
	return
}

func DecodeWeatherStatus(weatherData *Location) (code int, data *cache.WeatherData, err error) {
	data = &cache.WeatherData{
		LocationName:  weatherData.LocationName,
		WeatherStatus: &cache.WeatherStatus{},
	}

	if len(weatherData.WeatherElements) == 0 {
		code = errorCode.DecodeJsonError
		return
	}
	now := time.Now()
	for _, eachInfo := range weatherData.WeatherElements {
		if eachInfo.ElementName == "Wx" {
			if len(eachInfo.Time) != 0 {
				data.StartTime = eachInfo.Time[0].StartTime
				data.EndTime = eachInfo.Time[0].EndTime
				data.Weather = eachInfo.Time[0].Parameter.Name
			}
		}

		if eachInfo.ElementName == "PoP" {
			if len(eachInfo.Time) != 0 {
				data.ChanceOfRain = eachInfo.Time[0].Parameter.Name + "%"
			}
		}

		if eachInfo.ElementName == "MinT" {
			if len(eachInfo.Time) != 0 {
				data.MinTemperature = eachInfo.Time[0].Parameter.Name + "C"
			}
		}

		if eachInfo.ElementName == "MaxT" {
			if len(eachInfo.Time) != 0 {
				data.MaxTemperature = eachInfo.Time[0].Parameter.Name + "C"
			}
		}
		data.UpdateTime = &now
	}
	return
}
