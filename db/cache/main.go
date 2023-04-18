package cache

type CacheServer struct {

	// weather
	SetWeatherData       chan *WeatherData
	GetWeatherDataReq    chan string
	GetWeatherDataRes    chan *WeatherData
	GetAllWeatherDataReq chan bool
	GetAllWeatherDataRes chan map[string]*WeatherData
	WeatherDataList      map[string]*WeatherData `json:"weather_status_list"` // mysql checker list

}

var (
	Server *CacheServer
)

func NewCacheServer() *CacheServer {
	Server = &CacheServer{

		// weather
		SetWeatherData:       make(chan *WeatherData),
		GetWeatherDataReq:    make(chan string),
		GetWeatherDataRes:    make(chan *WeatherData),
		GetAllWeatherDataReq: make(chan bool),
		GetAllWeatherDataRes: make(chan map[string]*WeatherData),
		WeatherDataList:      map[string]*WeatherData{},
	}
	return Server
}

func (s *CacheServer) Run() {

	go func() {
		for {
			select {
			case data := <-s.SetWeatherData:
				s.WeatherDataList[data.LocationName] = data
			case locationName := <-s.GetWeatherDataReq:
				if data, exist := s.WeatherDataList[locationName]; exist {
					s.GetWeatherDataRes <- data
					continue
				}
				s.GetWeatherDataRes <- nil
			case <-s.GetAllWeatherDataReq:
				s.GetAllWeatherDataRes <- s.WeatherDataList
			}
		}
	}()

}
