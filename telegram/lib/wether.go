package lib

import (
	"alarm-system/db/cache"
	"alarm-system/scraper"
)

// 格式化 從中央氣象局抓回來的天氣資料(一般天氣預報-今明36小時天氣預報)
func FormatWeatherData(weatherData *cache.WeatherData) (code int, content string, err error) {
	content = ""
	content += "地點:" + weatherData.LocationName + "\n"
	content += "預測時間區間:" + weatherData.StartTime + " - " + weatherData.EndTime + "\n"
	content += "天氣狀況:" + weatherData.Weather + "\n"
	content += "降雨機率:" + weatherData.ChanceOfRain + "\n"
	content += "最高溫:" + weatherData.MaxTemperature + "\n"
	content += "最低溫:" + weatherData.MinTemperature + "\n"
	content += "\n"
	return
}

// 格式化 從中央氣象局抓回來的天氣資料(鄉鎮天氣預報-XXX未來兩天天氣預報)
func FormatWeatherStatus(weatherStatus *scraper.WeatherStatus) (code int, content string, err error) {
	content = ""
	content += "預測時間區間:" + weatherStatus.StartTime + " - " + weatherStatus.EndTime + "\n"
	content += "天氣狀況:" + weatherStatus.Weather + "\n"
	content += "降雨機率:" + weatherStatus.ChanceOfRain + "\n"
	content += "溫度:" + weatherStatus.Temperature + "\n"
	content += "風速:" + weatherStatus.WindSpeed + "\n"
	content += "相對溼度:" + weatherStatus.RH + "\n"
	content += "\n"
	return
}
