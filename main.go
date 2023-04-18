package main

import (
	_ "alarm-system/config"
	"alarm-system/telegram"
	"os"
	"os/signal"
	"syscall"

	"alarm-system/db/cache"
	_ "alarm-system/telegram"

	"alarm-system/checker"

	_ "alarm-system/scraper"
	_ "net/http/pprof"

	"188.166.240.198/GAIUS/lib/logger"
)

func main() {

	cacheServer := cache.NewCacheServer()
	go cacheServer.Run()

	telegramServer := telegram.NewTelegramServer()

	go telegramServer.RunServer()
	go telegramServer.RunJob()

	checker.StartChecker()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	<-shutdown

	logger.Info("stop server .....")
}
