package api

// 發送給tg group的訊息
type SendMessageReq struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

// 接收 發生事故位置
type AccidentReq struct {
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Content   string  `json:"content"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
