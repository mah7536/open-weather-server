package scraper

import (
	"alarm-system/config"
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
	UrlForWeather = "/v1/rest/datastore/%s"
)

type EachTimeInfo struct {
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	ElementValue []struct {
		Measures string `json:"measures"`
		Value    string `json:"value"`
	} `json:"elementValue"`
}

type WeatherElement struct {
	ElementName string          `json:"elementName"`
	Time        []*EachTimeInfo `json:"time"`
}

type Location struct {
	LocationName    string            `json:"locationName"`
	WeatherElements []*WeatherElement `json:"weatherElement"`
}

type Locations struct {
	LocationsName string      `json:"locationsName"`
	Location      []*Location `json:"location"`
}

type Record struct {
	Locations []*Locations `json:"locations"`
}

type WeatherRes struct {
	Success string  `json:"success"`
	Records *Record `json:"records"`
}

type WeatherScraper struct {
	*config.WeatherChecker
	Client   *httpClient.Client `json:"-"`
	Req      chan *WeatherReq
	Res      chan *WeatherStatus
	Duration int
}

type WeatherStatus struct {
	ChanceOfRain string     `json:"chance_of_rain"`
	Weather      string     `json:"whether"`
	Temperature  string     `json:"temperature"`
	WindSpeed    string     `json:"wind_speed"`
	RH           string     `json:"r_h"`
	UpdateTime   *time.Time `json:"last_check_time"`
	StartTime    string     `json:"start_time"`
	EndTime      string     `json:"end_time"`
}

type WeatherReq struct {
	CityName string
	TownName string
}

var WeatherC *WeatherScraper

func NewWeatherScraper(config *config.WeatherChecker) *WeatherScraper {
	return &WeatherScraper{
		WeatherChecker: config,
		Client:         httpClient.NewClient(),
		Req:            make(chan *WeatherReq),
		Res:            make(chan *WeatherStatus),
		Duration:       0,
	}
}

func NewWeatherRequest(url string, townName string, token string) (code int, req *http.Request, err error) {
	// currentHour := time.Now().Hour()

	// startHour := (currentHour / 3) * 3
	// endHour := startHour + 3

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		code = errorCode.Error
	}
	q := req.URL.Query()
	q.Add("Authorization", token)
	q.Add("locationName", townName)
	req.URL.RawQuery = q.Encode()
	return
}

func (scraper *WeatherScraper) Run() {
	for {

		select {
		case location := <-scraper.Req:
			code, data, err := scraper.Job(location)
			if code != errorCode.Success {
				res := errorhandler.NewResponse(code)
				res.SetExtra(err)
				logger.Error(res)
				scraper.Res <- nil
				continue
			}

			scraper.Res <- data
		}
	}
}

func (scraper *WeatherScraper) Job(param *WeatherReq) (code int, data *WeatherStatus, err error) {

	router, exist := config.LocationAndCodeMap[param.CityName]
	if !exist {
		code = errorCode.DBNoData
		return
	}

	code, req, err := NewWeatherRequest(fmt.Sprintf(scraper.Url+UrlForWeather, router.Code), param.TownName, scraper.Token)
	if code != errorCode.Success {
		code = errorCode.RequestCreateError
		return
	}

	start := time.Now()
	code, res, err := scraper.Client.Send(req)
	if err != nil {
		code = errorCode.RequestSendError
		return
	}
	weatherRes := &WeatherRes{}
	err = json.Unmarshal(res, weatherRes)
	if err != nil {
		code = errorCode.DecodeJsonError
		return
	}

	if weatherRes.Success != "true" {
		code = errorCode.Error
		return
	}

	if len(weatherRes.Records.Locations) == 0 {
		code = errorCode.DBNoData
		return
	}
	code, data, err = DecodeWeatherStatus(weatherRes.Records.Locations[0].Location[0])
	if code != errorCode.Success {
		return
	}
	scraper.Duration = int(time.Now().Sub(start).Seconds())
	return
}

func DecodeWeatherStatus(weatherData *Location) (code int, data *WeatherStatus, err error) {
	data = &WeatherStatus{}
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
				data.Weather = eachInfo.Time[0].ElementValue[0].Value
			}
		}

		if eachInfo.ElementName == "PoP12h" {
			if len(eachInfo.Time) != 0 {
				data.ChanceOfRain = eachInfo.Time[0].ElementValue[0].Value + "%"
			}
		}

		if eachInfo.ElementName == "T" {
			if len(eachInfo.Time) != 0 {
				data.Temperature = eachInfo.Time[0].ElementValue[0].Value + "C"
			}
		}

		if eachInfo.ElementName == "WS" {
			if len(eachInfo.Time) != 0 {
				data.WindSpeed = eachInfo.Time[0].ElementValue[0].Value + "m/s"
			}
		}

		if eachInfo.ElementName == "RH" {
			if len(eachInfo.Time) != 0 {
				data.RH = eachInfo.Time[0].ElementValue[0].Value + "%(濕度)"
			}
		}
		data.UpdateTime = &now
	}
	return
}
