package telegram

import (
	"alarm-system/telegram/lib"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	SomeOneTalkToBot = "有人與機器人開通, 請注意 該人的名稱為%s %s 語系為%s Id為%d"
)

func Start(message *tgbotapi.Update) (code int, res tgbotapi.MessageConfig, err error) {
	res = lib.NewResponseMs(message.Message.Chat.ID, fmt.Sprintf("Hi! %s, welcome ", message.Message.From.UserName))

	return
}
