package scraper

import (
	"alarm-system/config"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"188.166.240.198/GAIUS/lib/errorCode"
	"188.166.240.198/GAIUS/lib/httpClient"
	"188.166.240.198/GAIUS/lib/logger"
)

const (
	UrlForGetLocation = "other/TownVillagePointQuery"
)

var (
	LocationC *LocationChecker
)

type GNSS struct {
	Lat float32
	Lng float32
}

type LocationErrorRes struct {
	Error   string `xml:"error"`
	Message string `xml:"message"`
}

type LocationRes struct {
	XMLName xml.Name `xml:"townVillageItem"`
	// 正常回傳時
	CityName    string `xml:"ctyName"`
	TownName    string `xml:"townName"`
	VillageName string `xml:"villageName"`

	// 查無資料時
	Error *LocationErrorRes `xml:"error"`
}

type LocationChecker struct {
	*config.LocationConfig
	Client  *httpClient.Client `json:"-"`
	ReqChan chan *GNSS
	ResChan chan *LocationRes
}

func NewLocationScraper(config *config.LocationConfig) *LocationChecker {
	return &LocationChecker{
		LocationConfig: config,
		Client:         httpClient.NewClient(),
		ReqChan:        make(chan *GNSS),
		ResChan:        make(chan *LocationRes),
	}
}

func NewLocationRequest(url string, gnss *GNSS) (code int, req *http.Request, err error) {

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		code = errorCode.Error
	}
	req.URL.Path = fmt.Sprintf("%s/%f/%f", req.URL.Path, gnss.Lng, gnss.Lat)
	return
}

func (checker *LocationChecker) Run() {
	for {
		select {
		case gnss := <-checker.ReqChan:
			code, locationData, err := checker.Job(&GNSS{
				Lat: gnss.Lat,
				Lng: gnss.Lng,
			})
			if code != errorCode.Success {
				logger.Error(err)
				checker.ResChan <- nil
				continue
			}

			if locationData.Error != nil {
				logger.Error(fmt.Sprintf("座標(經度,緯度):%f,%f 錯誤訊息:%s", gnss.Lng, gnss.Lat, locationData.Error.Message))
				checker.ResChan <- nil
				continue
			}
			checker.ResChan <- locationData
		}
	}
}

func (checker *LocationChecker) Job(gnss *GNSS) (code int, data *LocationRes, err error) {

	code, req, err := NewLocationRequest(checker.Url+UrlForGetLocation, gnss)
	if code != errorCode.Success {
		code = errorCode.RequestCreateError
		return
	}

	start := time.Now()
	code, res, err := checker.Client.Send(req)
	if err != nil {
		code = errorCode.RequestSendError
		return
	}
	data = &LocationRes{}
	err = xml.Unmarshal(res, data)
	if err != nil {
		code = errorCode.Error
		return
	}
	logger.Info(fmt.Sprintf("此次Location Reqeust共花費%f秒", time.Now().Sub(start).Seconds()))
	return
}
