package cache

import (
	"time"
)

type WeatherData struct {
	LocationName string `json:"location_name"`
	*WeatherStatus
}

type WeatherStatus struct {
	ChanceOfRain   string     `json:"chance_of_rain"`
	Weather        string     `json:"whether"`
	MaxTemperature string     `json:"max_temperature"`
	MinTemperature string     `json:"min_temperature"`
	UpdateTime     *time.Time `json:"last_check_time"`
	StartTime      string     `json:"start_time"`
	EndTime        string     `json:"end_time"`
}

func NewWeatherStatus(locationName string) *WeatherData {
	now := time.Now()
	return &WeatherData{
		LocationName: locationName,
		WeatherStatus: &WeatherStatus{
			ChanceOfRain:   "",
			Weather:        "",
			MaxTemperature: "",
			MinTemperature: "",
			StartTime:      "",
			EndTime:        "",
			UpdateTime:     &now,
		},
	}
}
