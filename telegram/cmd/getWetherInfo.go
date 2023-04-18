package cmd

import (
	"alarm-system/db/cache"
	"alarm-system/telegram/callback"
	"alarm-system/telegram/lib"

	"188.166.240.198/GAIUS/lib/errorCode"
	"188.166.240.198/GAIUS/lib/logger"
	"188.166.240.198/GAIUS/lib/selfTime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetWeatherInfo(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {
	cache.Server.GetAllWeatherDataReq <- true
	weatherInfoList := <-cache.Server.GetAllWeatherDataRes

	content := ""
	i := 0
	for _, each := range weatherInfoList {

		i++
		tmpCode, tmpContent, tmpErr := lib.FormatWeatherData(each)
		if tmpCode != errorCode.Success {
			logger.Error(tmpErr)
			return
		}
		content += tmpContent
		if i == len(weatherInfoList) {
			content += "最後更新時間" + each.UpdateTime.Format(selfTime.TimeLayout) + "\n"
		}
	}

	res = lib.NewResponseMs(message.Message.Chat.ID, content)
	return
}

func GetWeatherList(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {

	cache.Server.GetAllWeatherDataReq <- true
	weatherInfoList := <-cache.Server.GetAllWeatherDataRes
	res = lib.NewResponseMs(message.Message.Chat.ID, "選取地區")

	list := [][]tgbotapi.InlineKeyboardButton{}

	// 這邊是將所有地區整理成選單
	rowList := tgbotapi.NewInlineKeyboardRow()

	for _, each := range weatherInfoList {
		qCode, data, _ := lib.SetCallBackReq(callback.PerWeatherInfo, &callback.GetPerWeatherInfoReq{
			LocationName: each.LocationName,
		})
		if qCode != errorCode.Success {
			return
		}

		rowList = append(rowList, tgbotapi.NewInlineKeyboardButtonData(each.LocationName, data))
		if len(rowList) == 4 {
			list = append(list, rowList)
			rowList = tgbotapi.NewInlineKeyboardRow()
		}
	}
	if len(rowList) != 0 {
		list = append(list, rowList)
	}

	res.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(list...)
	return
}
