package telegram

import (
	"alarm-system/config"
	"alarm-system/telegram/callback"
	"alarm-system/telegram/cmd"
	"alarm-system/telegram/lib"
	"os"
	"strings"

	"188.166.240.198/GAIUS/lib/errorCode"

	"188.166.240.198/GAIUS/lib/errorhandler"
	"188.166.240.198/GAIUS/lib/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	TelegramSys *TelegramServer
)

type command struct {
	tgcommand tgbotapi.BotCommand
	fn        func(update *tgbotapi.Update) (int, tgbotapi.MessageConfig, error)
}

type callbackCommand struct {
	tgcommand tgbotapi.BotCommand
	fn        func(message *tgbotapi.Update, req []byte) (int, tgbotapi.MessageConfig, error)
}

type TelegramServer struct {
	tgChannel    tgbotapi.UpdatesChannel
	mainBot      *tgbotapi.BotAPI
	commandList  map[string]*command
	callbackList map[string]*callbackCommand
	sendChan     chan tgbotapi.Chattable
	callbackChan chan tgbotapi.CallbackConfig
	deletemsChan chan tgbotapi.DeleteMessageConfig
}

func NewTelegramServer() *TelegramServer {
	bot, err := tgbotapi.NewBotAPI(config.ServerConfig.TelegramToken)
	if err != nil {
		logger.Error(err)
		os.Exit(2)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	channel, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Error(err)
		os.Exit(2)
	}

	TelegramSys = &TelegramServer{
		tgChannel:    channel,
		mainBot:      bot,
		commandList:  make(map[string]*command),
		callbackList: make(map[string]*callbackCommand),
		sendChan:     make(chan tgbotapi.Chattable),
		callbackChan: make(chan tgbotapi.CallbackConfig),
		deletemsChan: make(chan tgbotapi.DeleteMessageConfig),
	}

	// 新增command
	TelegramSys.AddCommandList("start", "for start command", Start)
	TelegramSys.AddCommandList("isalive", "check sysOk", IsAlive)
	TelegramSys.AddCommandList("weather", "get weather info", cmd.GetWeatherInfo)
	TelegramSys.AddCommandList("weather_list", "get weather list", cmd.GetWeatherList)
	TelegramSys.AddCommandList("location", "get location", cmd.GetLocationButton)

	// 新增callback
	TelegramSys.AddCallBackList(callback.PerWeatherInfo, "PerWeatherInfo", callback.GetPerWeatherInfo)

	// 設定command list
	_, _, err = TelegramSys.SetTgCommandList()
	if err != nil {
		logger.Error(err)
		os.Exit(2)
	}

	return TelegramSys
}

func (server *TelegramServer) RunJob() {
	for {
		select {
		case ms := <-server.sendChan:
			_, err := server.mainBot.Send(ms)
			if err != nil {
				logger.Error(err)
			}
		case callback := <-server.callbackChan:
			_, err := server.mainBot.AnswerCallbackQuery(callback)
			if err != nil {
				logger.Error(err)
			}
		case delete := <-server.deletemsChan:
			_, err := server.mainBot.DeleteMessage(delete)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// 將可用的指令 統一整理
func (server *TelegramServer) AddCommandList(act string, des string, fn func(update *tgbotapi.Update) (int, tgbotapi.MessageConfig, error)) {
	server.commandList[act] = &command{
		tgcommand: tgbotapi.BotCommand{
			Command:     act,
			Description: des,
		},
		fn: fn,
	}
}

// 將callback指令 統一整理
func (server *TelegramServer) AddCallBackList(act string, des string, fn func(message *tgbotapi.Update, req []byte) (int, tgbotapi.MessageConfig, error)) {
	server.callbackList[act] = &callbackCommand{
		tgcommand: tgbotapi.BotCommand{
			Command:     act,
			Description: des,
		},
		fn: fn,
	}
}

// 設定機器人的menu list
func (server *TelegramServer) SetTgCommandList() (code int, data interface{}, err error) {
	tgCommandList := []tgbotapi.BotCommand{}
	for _, command := range server.commandList {
		if command.tgcommand.Command != "start" {
			tgCommandList = append(tgCommandList, command.tgcommand)
		}
	}
	err = server.mainBot.SetMyCommands(tgCommandList)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// 傳送訊息
func SendMs(ms tgbotapi.Chattable) {
	TelegramSys.sendChan <- ms
}

// run telegram server
func (server *TelegramServer) RunServer() {
	for message := range server.tgChannel {
		if message.Message != nil {

			// 此區為測試用的gps
			if message.Message.Location != nil {
				code, ms, err := cmd.SendLocation(&message)
				if code != errorCode.Success {
					res := errorhandler.NewResponse(code)
					logger.Error(res.SetExtra(err))
				}
				server.sendChan <- ms
				continue
			}

			// 僅針對 是command的message
			if message.Message.IsCommand() {
				command, exist := server.commandList[strings.ToLower(message.Message.Command())]
				if exist {
					code, res, err := command.fn(&message)
					if code != errorCode.Success {
						response := errorhandler.NewResponse(code)
						response.SetExtra(err)
						logger.Error(response)
						server.sendChan <- (lib.WarnMessage(message.CallbackQuery.Message.Chat.ID, response.Message))
						continue
					}
					server.sendChan <- res
				} else {
					server.sendChan <- (lib.WarnMessage(message.Message.Chat.ID, "不在喔喔喔喔"))
				}
				continue
			}

			if message.Message.LeftChatMember != nil {
				continue
			}

			server.sendChan <- lib.AlertMessage(message.Message.Chat.ID, "Hi")
			continue
		}

		// 針對callback query的message type
		if message.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(message.CallbackQuery.ID, "")
			server.callbackChan <- callback

			code, _, _ := lib.CheckChatID(message.CallbackQuery.Message.Chat.ID)
			if code != errorCode.Success {
				logger.Error(code)
				server.sendChan <- lib.AlertMessage(message.Message.Chat.ID, "Hi")
				continue
			}

			logger.Notice("拿到callback Data" + message.CallbackQuery.Data)

			code, req, _ := lib.StringToReq(message.CallbackQuery.Data)
			if code != errorCode.Success {
				response := errorhandler.NewResponse(code)
				logger.Error(response)
				server.sendChan <- (lib.WarnMessage(message.CallbackQuery.Message.Chat.ID, response.Message))
				continue
			}

			callbackFn, exist := server.callbackList[req.Action]
			if !exist {
				server.sendChan <- (lib.WarnMessage(message.CallbackQuery.Message.Chat.ID, "不在喔喔喔喔"))
				continue
			}

			code, res, err := callbackFn.fn(&message, []byte(req.Req))
			if code != errorCode.Success {
				response := errorhandler.NewResponse(code)
				response.SetExtra(err)
				logger.Error(response)
				server.sendChan <- (lib.WarnMessage(message.CallbackQuery.Message.Chat.ID, response.Message))
				continue
			}
			server.sendChan <- res

			deleteMs := tgbotapi.NewDeleteMessage(message.CallbackQuery.Message.Chat.ID, message.CallbackQuery.Message.MessageID)
			server.deletemsChan <- deleteMs

		}
	}
}
