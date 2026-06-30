package model

// TouristSpot 代表台灣觀光景點
type TouristSpot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Px          float64 `json:"px"`         // 經度
	Py          float64 `json:"py"`         // 緯度
	Address     string  `json:"add"`        // 地址
	Picture1    string  `json:"picture1"`   // 景點照片網址
}

// TripNode 代表旅程中的其中一站
type TripNode struct {
	Step       int         `json:"step"`
	Spot       TouristSpot `json:"spot"`
	Transport  string      `json:"transport"`
	TravelTime int         `json:"travel_time"`
}

// UserTrip 記錄整趟旅程的當前狀態
type UserTrip struct {
	CurrentCoords [2]float64 `json:"current_coords"` // [lng, lat]
	History       []TripNode `json:"history"`
}
