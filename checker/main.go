package checker

import (
	"alarm-system/config"
)

func StartChecker() {

	newWeatherChcker := NewWeatherChecker(config.ServerConfig.WeatherChecker)
	go newWeatherChcker.Run()

}
