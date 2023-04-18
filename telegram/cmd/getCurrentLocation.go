package cmd

import (
	"alarm-system/scraper"
	"alarm-system/telegram/lib"
	"fmt"

	"188.166.240.198/GAIUS/lib/errorCode"
	"188.166.240.198/GAIUS/lib/selfTime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetLocationButton(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {

	res = lib.NewResponseMs(message.Message.Chat.ID, "選取地區")

	// 這邊是將所有地區整理成選單
	rowList := tgbotapi.NewKeyboardButtonRow()
	rowList = append(rowList, tgbotapi.NewKeyboardButtonLocation("Send Location"))

	res.ReplyMarkup = tgbotapi.NewReplyKeyboard(rowList)
	return
}

func SendLocation(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {
	scraper.LocationC.ReqChan <- &scraper.GNSS{
		Lat: float32(message.Message.Location.Latitude),
		Lng: float32(message.Message.Location.Longitude),
	}
	locationInfo := <-scraper.LocationC.ResChan

	if locationInfo == nil {
		code = errorCode.Error
		res = lib.NewCommonMessage(message.Message.Chat.ID, lib.TypeDanger, "不在台灣", "你沒在台灣喔")
		return
	}

	code = errorCode.Success

	content := fmt.Sprintf("你在 %s %s %s\n", locationInfo.CityName, locationInfo.TownName, locationInfo.VillageName)

	scraper.WeatherC.Req <- &scraper.WeatherReq{
		CityName: locationInfo.CityName,
		TownName: locationInfo.TownName,
	}

	weatherStatus := <-scraper.WeatherC.Res
	code, tmpContent, err := lib.FormatWeatherStatus(weatherStatus)
	if code != errorCode.Success {
		res = lib.NewResponseMs(message.Message.Chat.ID, content)
		return
	}
	content += tmpContent
	content += "最後更新時間" + weatherStatus.UpdateTime.Format(selfTime.TimeLayout) + "\n"
	res = lib.NewResponseMs(message.Message.Chat.ID, content)
	return
}
