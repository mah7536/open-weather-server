package callback

import (
	"alarm-system/db/cache"
	"alarm-system/telegram/lib"
	"encoding/json"

	"188.166.240.198/GAIUS/lib/errorCode"
	"188.166.240.198/GAIUS/lib/selfTime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	PerWeatherInfo = "p"
)

type GetPerWeatherInfoReq struct {
	LocationName string `json:"l"`
}

func GetPerWeatherInfo(message *tgbotapi.Update, req []byte) (code int, res tgbotapi.MessageConfig, err error) {
	request := &GetPerWeatherInfoReq{}
	err = json.Unmarshal(req, request)
	if err != nil {
		code = errorCode.DecodeJsonError
		return
	}

	cache.Server.GetWeatherDataReq <- request.LocationName
	weatherInfo := <-cache.Server.GetWeatherDataRes

	if weatherInfo == nil {
		code = errorCode.DBNoData
		return
	}

	code, content, err := lib.FormatWeatherData(weatherInfo)
	if err != nil {
		return
	}

	content += "最後更新時間" + weatherInfo.UpdateTime.Format(selfTime.TimeLayout) + "\n"

	res = lib.NewResponseMs(message.CallbackQuery.Message.Chat.ID, content)
	return
}
