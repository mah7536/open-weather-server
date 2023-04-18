package telegram

import (
	"alarm-system/telegram/lib"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TestJson struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
}

func IsAlive(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {
	res = lib.NewResponseMs(message.Message.Chat.ID, "👍")
	return
}
