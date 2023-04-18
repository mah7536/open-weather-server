package scraper

import "alarm-system/config"

func init() {
	LocationC = NewLocationScraper(config.ServerConfig.LocationConfig)
	go LocationC.Run()

	WeatherC = NewWeatherScraper(config.ServerConfig.WeatherChecker)
	go WeatherC.Run()
}
