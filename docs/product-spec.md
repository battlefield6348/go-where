# 產品規格

## 專案願景與核心玩法

這是一個面向一般大眾的「旅遊盲盒 / 隨機公路旅行生成器」網頁應用。

核心流程：使用者打開網頁並允許 GPS 定位，接著設定交通工具與移動時間（例如：開車 1 小時）。按下按鈕後，系統會隨機抽選一個在該時間內可抵達的真實台灣觀光景點，並在地圖上插旗。

連續接關功能：使用者抵達景點後，可以以該景點為新起點，重新設定條件並抽選下一站。整天的移動軌跡會被記錄下來，連成一條冒險路線。

## 架構與技術選型約束

- 部署環境：GitHub Pages（完全靜態託管，無後端伺服器與資料庫成本）。
- 核心邏輯：使用 Golang 撰寫，遵循 Clean Architecture 原則，並編譯成 WebAssembly 在瀏覽器運行。
- 前端地圖：使用 Leaflet.js 搭配 CartoDB Voyager 圖磚，不使用 Google Maps API。
- 狀態保存：利用瀏覽器 LocalStorage 持久化行程狀態。
- 景點資料庫：內嵌台灣交通部觀光署開放資料 JSON，以 `//go:embed` 打包進 Wasm。

## 核心資料模型

```go
package model

// TouristSpot 代表台灣官方登記的真實觀光景點
type TouristSpot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Px          float64 `json:"px"`
	Py          float64 `json:"py"`
	Address     string  `json:"add"`
	Picture1    string  `json:"picture1"`
}

// TripNode 代表冒險旅程中的其中一站
type TripNode struct {
	Step       int         `json:"step"`
	Spot       TouristSpot `json:"spot"`
	Transport  string      `json:"transport"`
	TravelTime int         `json:"travel_time"`
}

// UserTrip 記錄整趟旅程的當前狀態
type UserTrip struct {
	CurrentCoords [2]float64 `json:"current_coords"`
	History       []TripNode `json:"history"`
}
```

## 核心算法需求

- 使用 Haversine 公式計算地表兩點之間的直線距離（公里）。
- 依輸入座標與半徑（公里）過濾景點清單，再從候選結果中隨機抽選。

## 下一步實作順序

1. `geo/calculator.go`：距離計算與半徑篩選。
2. `storage/localstorage.go`：透過 `syscall/js` 與 LocalStorage 同步 `UserTrip`。
3. `main.go`：註冊可由 JavaScript 呼叫的 Wasm 函式。
4. `web/index.html`：整合 Leaflet.js 與 Wasm 橋接。
