# 旅遊盲盒 / 隨機公路旅行生成器設計規格書

本專案 `go-where` 是一個純靜態部署於 GitHub Pages 的「旅遊盲盒 / 隨機公路旅行生成器」應用。核心邏輯使用 Go 編譯為 WebAssembly (Wasm) 並運行於瀏覽器中，前端介面使用 Leaflet.js 互動地圖與高級半透明玻璃擬態 (Glassmorphism) 控制面板。

---

## 1. 資料模型與景點資料集 (`model/trip.go` & `geo/`)

### 1.1 核心結構定義
```go
package model

type TouristSpot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Px          float64 `json:"px"`         // 經度 (Longitude)
	Py          float64 `json:"py"`         // 緯度 (Latitude)
	Address     string  `json:"add"`        // 地址
	Picture1    string  `json:"picture1"`   // 景點照片網址
}

type TripNode struct {
	Step       int         `json:"step"`
	Spot       TouristSpot `json:"spot"`
	Transport  string      `json:"transport"`
	TravelTime int         `json:"travel_time"`
}

type UserTrip struct {
	CurrentCoords [2]float64 `json:"current_coords"` // [lng, lat]
	History       []TripNode `json:"history"`
}
```

### 1.2 景點資料內嵌 (`geo/db.go` & `geo/spots.json`)
精選全台各大縣市約 100 個知名觀光景點，包含座標、名稱、描述與照片網址，打包於 `geo/spots.json` 中。
```go
package geo

import (
	_ "embed"
	"encoding/json"
	"go-where/model"
)

//go:embed spots.json
var spotsData []byte

var CachedSpots []model.TouristSpot

func init() {
	if err := json.Unmarshal(spotsData, &CachedSpots); err != nil {
		panic("Failed to parse embedded spots data: " + err.Error())
	}
}
```

---

## 2. Go WASM 邏輯與 JS 橋接 (`geo/calculator.go`, `storage/localstorage.go`, `main.go`)

### 2.1 距離與抽選計算 (`geo/calculator.go`)
- 採用 **Haversine** 公式計算地表兩點間直線距離。
- 交通工具半徑對應：
  - `walking` (步行): 5 km/h
  - `cycling` (單車): 15 km/h
  - `driving` (汽車): 60 km/h
- 若篩選半徑內無景點，則採用備份方案隨機抽選最鄰近的 3 個景點之一，避免旅程中斷。

### 2.2 LocalStorage 狀態儲存 (`storage/localstorage.go`)
- 透過 `syscall/js` 操作瀏覽器 `localStorage`，鍵值為 `go_where_trip_state`。
- 支援讀取 (`LoadTrip`)、寫入 (`SaveTrip`) 與清除 (`ClearTrip`)。

### 2.3 WASM JS 介面橋接 (`main.go`)
- `window.goWhereGetState()`: 取得當前旅程 JSON 字串。
- `window.goWhereDraw(transport, timeMinutes)`: 進行下一站抽選，並更新歷史紀錄與儲存。
- `window.goWhereReset()`: 重設狀態。
- `window.goWhereUpdateCoords(lng, lat)`: 更新起點定位座標。

---

## 3. 前端 Layout 與高級 UI 設計 (`web/index.html`)

### 3.1 視覺美感設計
- **色調**：深色科技風 (`#0b0f19` 背景、`#00f5d4` 霓虹青起點標示、`#7b2cbf` 紫色盲盒特效)。
- **版面**：全螢幕 Leaflet.js 地圖為底層，上方浮動 **Glassmorphism** 半透明控制面版與卡片。
- **微交互與動畫**：
  - 景點抽取時地圖中央產生雷達掃描擴散波動。
  - 地圖上的歷史軌跡連線採用流暢的虛線前進動畫 (flow animation)。
  - 抽取結果卡片使用平滑的 `Slide-in` 與浮動投影效果呈現。

---

## 4. CI/CD 與部署流程

### 4.1 本地端測試伺服器 (`server.go`)
- 由於 CORS 限制，本機提供簡單的 `server.go` 靜態檔案伺服器，可在 `localhost:8080` 直接預覽。

### 4.2 GitHub Actions 自動化部署 (`.github/workflows/deploy.yml`)
- 推送至 `main` 分支時，觸發自動編譯 Go 至 Wasm 並部署 `web/` 資料夾至 `gh-pages` 分支。
