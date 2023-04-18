package config

type LocationAndLocationCode struct {
	Location string
	Code     string
}

var LocationAndCodeMap = make(map[string]*LocationAndLocationCode)

func AddToMap(location string, code string) {
	LocationAndCodeMap[location] = &LocationAndLocationCode{
		Code:     code,
		Location: location,
	}
}

func init() {
	AddToMap("連江縣", "F-D0047-081")
	AddToMap("新北市", "F-D0047-069")
	AddToMap("苗栗縣", "F-D0047-013")
	AddToMap("臺中市", "F-D0047-073")
	AddToMap("桃園市", "F-D0047-005")
	AddToMap("南投縣", "F-D0047-021")
	AddToMap("澎湖縣", "F-D0045-045")
	AddToMap("彰化縣", "F-D0047-017")
	AddToMap("嘉義縣", "F-D0047-029")
	AddToMap("新竹縣", "F-D0047-009")
	AddToMap("宜蘭縣", "F-D0047-001")
	AddToMap("臺東縣", "F-D0047-037")
	AddToMap("屏東縣", "F-D0047-033")
	AddToMap("基隆市", "F-D0047-049")
	AddToMap("嘉義市", "F-D0047-057")
	AddToMap("臺南市", "F-D0047-077")
	AddToMap("新竹市", "F-D0047-053")
	AddToMap("臺北市", "F-D0047-061")
	AddToMap("雲林縣", "F-D0047-025")
	AddToMap("花蓮縣", "F-D0047-041")
	AddToMap("高雄市", "F-D0047-065")
	AddToMap("金門縣", "F-D0047-085")
}
