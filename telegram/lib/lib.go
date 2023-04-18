package lib

import (
	"alarm-system/config"
	"encoding/json"
	"fmt"
	"strings"

	"188.166.240.198/GAIUS/lib/errorCode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	seperateLine = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, "=========")
	msTypeList   = map[string]string{
		TypeWarn:     "⚠️⚠️⚠️⚠️*Warn*⚠️⚠️⚠️⚠️",        //warn
		TypeDanger:   "❗️❗️❗️❗️❗️*Danger*❗️❗️❗️❗️❗️❗", //danger
		TypeInfo:     "👌👌👌👌👌*Info*👌👌👌👌👌",              //info
		TypeCommon:   "👍👍👍👍*Common*👍👍👍👍",              //common
		TypeUndefind: "❔❔❔❔*Undefined*❔❔❔❔",           //undefinded
	}
	StandardFormat = " %s \n " + seperateLine + " %s " + seperateLine + "\n`" + "%s" + "`"
)

const (
	TypeWarn     = "warn"
	TypeDanger   = "danger"
	TypeInfo     = "info"
	TypeCommon   = "common"
	TypeUndefind = "undefined"

	StatusStringYes = "是"
	StatusStringNo  = "否"
)

type CallBackReq struct {
	Action string `json:"ac"`
	Req    string `json:"req"`
}

// passive message
func NewResponseMs(chatID int64, text string) (newMS tgbotapi.MessageConfig) {
	newMS = tgbotapi.NewMessage(chatID, text)
	return
}

func AlertMessage(chatID int64, text string) (newMS tgbotapi.MessageConfig) {
	text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
	newMS = tgbotapi.NewMessage(chatID, fmt.Sprintf("`%s`", text))
	newMS.ParseMode = tgbotapi.ModeMarkdownV2
	return
}

func WarnMessage(chatID int64, text string) (newMS tgbotapi.MessageConfig) {
	text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
	newMS = tgbotapi.NewMessage(chatID, fmt.Sprintf("*%s*", text))
	newMS.ParseMode = tgbotapi.ModeMarkdownV2
	return
}

func DangerMessage(chatID int64, text string) (newMS tgbotapi.MessageConfig) {
	text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
	newMS = tgbotapi.NewMessage(chatID, fmt.Sprintf("__%s__", text))
	newMS.ParseMode = tgbotapi.ModeMarkdownV2
	return
}

func CheckChatID(id int64) (code int, data interface{}, err error) {
	for _, userId := range config.AllowChatId {
		if userId == id {
			code = errorCode.Success
			return
		}
	}
	code = errorCode.TgNotFoundUser
	return
}

// active message
// 發送一般訊息
func NewCommonMessage(chatId int64, msType string, title string, text string) (newMS tgbotapi.MessageConfig) {
	title = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, title)
	text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
	header, exist := msTypeList[strings.ToLower(msType)]
	if !exist {
		header = msTypeList[TypeUndefind]
	}
	newMS = tgbotapi.NewMessage(chatId, fmt.Sprintf(StandardFormat, header, title, text))
	newMS.ParseMode = tgbotapi.ModeMarkdownV2
	return
}

// 傳送事件 內容 及發生位置
func VenueMessage(chatId int64, msType string, title string, text string, latitude float64, longitude float64) (newVenue tgbotapi.VenueConfig) {
	text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
	header, exist := msTypeList[strings.ToLower(msType)]
	if !exist {
		header = msTypeList[TypeUndefind]
	}
	newVenue = tgbotapi.NewVenue(chatId, header+"\n"+title, text, latitude, longitude)
	return
}

func JsonToString(jsonData interface{}) (code int, data string, err error) {
	byteData, err := json.Marshal(jsonData)
	if err != nil {
		code = errorCode.EncodeJsonError
		return
	}
	data = string(byteData)
	return
}

func StringToReq(reqStr string) (code int, req *CallBackReq, err error) {
	req = &CallBackReq{}
	err = json.Unmarshal([]byte(reqStr), req)
	if err != nil {
		code = errorCode.DecodeJsonError
		return
	}
	return
}

func SetCallBackReq(action string, req interface{}) (code int, data string, err error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		code = errorCode.EncodeJsonError
		return
	}
	byteData, jErr := json.Marshal(&CallBackReq{
		Action: action,
		Req:    string(reqData),
	})
	if jErr != nil {
		code = errorCode.EncodeJsonError
		return
	}
	data = string(byteData)
	return
}

func ConvertStatusToString(status bool) (statusString string) {
	statusString = StatusStringNo
	if status {
		statusString = StatusStringYes
	}
	return
}
