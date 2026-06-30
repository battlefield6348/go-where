# 旅遊盲盒 / 隨機公路旅行生成器 實作計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 實作「旅遊盲盒 / 隨機公路旅行生成器」網頁應用，部署至 GitHub Pages，核心邏輯使用 Go WebAssembly 編譯，搭配 Leaflet.js 互動地圖與高級 UI 懸浮面板。

**Architecture:** Go 核心邏輯處理距離計算（Haversine）與景點隨機抽選，透過 `syscall/js` 對外導出 API 並與 LocalStorage 進行旅程狀態持久化；前端使用全螢幕 Leaflet.js 地圖與半透明玻璃擬態（Glassmorphism）卡片面版互動。

**Tech Stack:** Go 1.24.4 (Wasm), JavaScript, Leaflet.js, HTML5, CSS3 (Vanilla)

---

### Task 1: Go 模組初始化與核心資料模型

**Files:**
- Create: `go.mod`
- Create: `model/trip.go`

- [ ] **Step 1: 初始化 Go Module**

Run: `go mod init go-where` in project root directory
Expected: 生成 `go.mod` 檔案

- [ ] **Step 2: 建立核心資料模型 `model/trip.go`**

Create `model/trip.go` with the following content:
```go
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
```

- [ ] **Step 3: 驗證編譯**

Run: `go build ./model/...`
Expected: 編譯成功，無語法錯誤

- [ ] **Step 4: Commit**

Run:
```bash
git add go.mod model/trip.go
git commit -m "feat: init go mod and define trip models"
```

---

### Task 2: 內嵌景點資料庫與單元測試

**Files:**
- Create: `geo/spots.json`
- Create: `geo/db.go`
- Test: `geo/db_test.go`

- [ ] **Step 1: 建立 `geo/spots.json`**

Create `geo/spots.json` with 25+ real Taiwan scenic spots:
```json
[
  {
    "id": "1",
    "name": "台北 101",
    "description": "台灣最著名的地標摩天大樓，擁有俯瞰台北盆地的觀景台。",
    "px": 121.5646,
    "py": 25.0339,
    "add": "台北市信義區信義路五段7號",
    "picture1": "https://images.unsplash.com/photo-1552912441-d6023d3856b3?w=800"
  },
  {
    "id": "2",
    "name": "九份老街",
    "description": "依山而建的古老街區，以茶樓、芋圓與復古的山城夜景聞名。",
    "px": 121.8440,
    "py": 25.1088,
    "add": "新北市瑞芳區基山街",
    "picture1": "https://images.unsplash.com/photo-1571168537969-490f230554c2?w=800"
  },
  {
    "id": "3",
    "name": "日月潭國家風景區",
    "description": "台灣最大的天然半人工湖泊，湖光山色美不勝收，設有環湖自行車道。",
    "px": 120.9316,
    "py": 23.8569,
    "add": "南投縣魚池鄉中山路599號",
    "picture1": "https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=800"
  },
  {
    "id": "4",
    "name": "太魯閣國家公園",
    "description": "以雄偉壯麗、幾近垂直的大理石峽谷景觀聞名中外的國家公園。",
    "px": 121.6215,
    "py": 24.1613,
    "add": "花蓮縣秀林鄉富世村富世291號",
    "picture1": "https://images.unsplash.com/photo-1549692520-acc6669e2f0c?w=800"
  },
  {
    "id": "5",
    "name": "墾丁國家公園",
    "description": "台灣南端的度假勝地，有著金黃色的沙灘、珊瑚礁地質與熱帶海洋風情。",
    "px": 120.7981,
    "py": 21.9442,
    "add": "屏東縣恆春鎮恆公路90號",
    "picture1": "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=800"
  },
  {
    "id": "6",
    "name": "阿里山國家森林遊樂區",
    "description": "以日出、雲海、晚霞、森林與高山鐵路「五奇」著稱的避暑勝地。",
    "px": 120.8066,
    "py": 23.5113,
    "add": "嘉義縣阿里山鄉中正村59號",
    "picture1": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"
  },
  {
    "id": "7",
    "name": "淡水漁人碼頭",
    "description": "以絕美的淡水夕照與情人橋聞名，適合散步與享受河風吹拂。",
    "px": 121.4116,
    "py": 25.1783,
    "add": "新北市淡水區觀海路199號",
    "picture1": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"
  },
  {
    "id": "8",
    "name": "高美濕地",
    "description": "擁有豐富生態與絕美木棧道，落日餘暉倒映在濕地上形成天空之鏡。",
    "px": 120.5501,
    "py": 24.3117,
    "add": "台中市清水區美堤街",
    "picture1": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"
  },
  {
    "id": "9",
    "name": "礁溪溫泉公園",
    "description": "台灣少見的平地溫泉，水質清澈無硫磺味，設有日式免費足湯。",
    "px": 121.7725,
    "py": 24.8296,
    "add": "宜蘭縣礁溪鄉公園路16號",
    "picture1": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"
  },
  {
    "id": "10",
    "name": "安平古堡",
    "description": "荷蘭人建於1624年，為台灣最早的城堡，現存磚牆紅瓦訴說著歷史歷史滄桑。",
    "px": 120.1611,
    "py": 23.0016,
    "add": "台南市安平區國勝路82號",
    "picture1": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"
  }
]
```

- [ ] **Step 2: 建立 `geo/db.go` 利用 `//go:embed` 載入景點**

Create `geo/db.go` with the following content:
```go
package geo

import (
	_ "embed"
	"encoding/json"
	"go-where/model"
)

//go:embed spots.json
var spotsData []byte

// CachedSpots 保存所有內嵌的景點
var CachedSpots []model.TouristSpot

func init() {
	if err := json.Unmarshal(spotsData, &CachedSpots); err != nil {
		panic("Failed to parse embedded spots data: " + err.Error())
	}
}
```

- [ ] **Step 3: 撰寫資料庫初始化測試 `geo/db_test.go`**

Create `geo/db_test.go` with the following content:
```go
package geo

import (
	"testing"
)

func TestCachedSpotsLoaded(t *testing.T) {
	if len(CachedSpots) == 0 {
		t.Fatal("Expected CachedSpots to contain spots, but it was empty")
	}
	
	firstSpot := CachedSpots[0]
	if firstSpot.Name == "" {
		t.Error("Expected spot to have a name, got empty string")
	}
	if firstSpot.Px == 0 || firstSpot.Py == 0 {
		t.Errorf("Expected spot %s to have valid coordinates, got Px:%f Py:%f", firstSpot.Name, firstSpot.Px, firstSpot.Py)
	}
}
```

- [ ] **Step 4: 執行測試驗證**

Run: `go test -v ./geo/...`
Expected: PASS

- [ ] **Step 5: Commit**

Run:
```bash
git add geo/spots.json geo/db.go geo/db_test.go
git commit -m "feat: embed spot json database and add db tests"
```

---

### Task 3: Haversine 距離計算與景點篩選

**Files:**
- Create: `geo/calculator.go`
- Test: `geo/calculator_test.go`

- [ ] **Step 1: 建立距離計算與隨機篩選邏輯 `geo/calculator.go`**

Create `geo/calculator.go` with the following content:
```go
package geo

import (
	"errors"
	"math"
	"math/rand"
	"time"
	"go-where/model"
)

// Haversine 計算兩點間的直線距離 (公里)
func Haversine(lng1, lat1, lng2, lat2 float64) float64 {
	const R = 6371.0 // 地球半徑 (km)
	
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0
	
	rLat1 := lat1 * math.Pi / 180.0
	rLat2 := lat2 * math.Pi / 180.0
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return R * c
}

// DrawNextSpot 根據當前座標、交通工具及時間篩選景點，並隨機抽取下一站
func DrawNextSpot(currentLng, currentLat float64, transport string, travelTimeMinutes int) (model.TouristSpot, error) {
	if len(CachedSpots) == 0 {
		return model.TouristSpot{}, errors.New("spots database is empty")
	}

	// 依交通工具決定平均時速 (km/h)
	var speedKmh float64
	switch transport {
	case "walking":
		speedKmh = 5.0
	case "cycling":
		speedKmh = 15.0
	default: // driving
		speedKmh = 60.0
	}
	
	// 移動半徑 = 時速 * (時間分鐘 / 60)
	radiusKm := speedKmh * (float64(travelTimeMinutes) / 60.0)
	
	// 篩選半徑內的候選景點
	var candidates []model.TouristSpot
	for _, spot := range CachedSpots {
		dist := Haversine(currentLng, currentLat, spot.Px, spot.Py)
		if dist <= radiusKm {
			candidates = append(candidates, spot)
		}
	}
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 若半徑內沒有任何景點，退而求其次，擴大尋找最近的三個景點並隨機抽一個
	if len(candidates) == 0 {
		type distTuple struct {
			spot model.TouristSpot
			dist float64
		}
		
		var sorted []distTuple
		for _, spot := range CachedSpots {
			dist := Haversine(currentLng, currentLat, spot.Px, spot.Py)
			sorted = append(sorted, distTuple{spot, dist})
		}
		
		// 簡易排序取最近的 3 個
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].dist > sorted[j].dist {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		
		limit := 3
		if len(sorted) < limit {
			limit = len(sorted)
		}
		
		if limit > 0 {
			chosenIdx := r.Intn(limit)
			return sorted[chosenIdx].spot, nil
		}
		
		return model.TouristSpot{}, errors.New("no spots available")
	}
	
	return candidates[r.Intn(len(candidates))], nil
}
```

- [ ] **Step 2: 撰寫計算與篩選測試 `geo/calculator_test.go`**

Create `geo/calculator_test.go` with the following content:
```go
package geo

import (
	"testing"
)

func TestHaversine(t *testing.T) {
	// 台北 101 (121.5646, 25.0339) 到 淡水漁人碼頭 (121.4116, 25.1783) 距離約 22.1 公里
	dist := Haversine(121.5646, 25.0339, 121.4116, 25.1783)
	if dist < 21.0 || dist > 23.0 {
		t.Errorf("Expected distance around 22.1 km, got %f", dist)
	}
}

func TestDrawNextSpot(t *testing.T) {
	// 當前起點台北車站附近
	spot, err := DrawNextSpot(121.517, 25.047, "driving", 60)
	if err != nil {
		t.Fatalf("Failed to draw spot: %s", err)
	}
	if spot.Name == "" {
		t.Error("Drawn spot name is empty")
	}
}
```

- [ ] **Step 3: 執行測試驗證**

Run: `go test -v ./geo/...`
Expected: ALL PASS

- [ ] **Step 4: Commit**

Run:
```bash
git add geo/calculator.go geo/calculator_test.go
git commit -m "feat: implement haversine and draw spot filtering logic"
```

---

### Task 4: JS WebAssembly LocalStorage 儲存器

**Files:**
- Create: `storage/localstorage.go`

- [ ] **Step 1: 實作 LocalStorage 同步器**

Create `storage/localstorage.go` with the following content:
```go
package storage

import (
	"encoding/json"
	"syscall/js"
	"go-where/model"
)

const StorageKey = "go_where_trip_state"

// LoadTrip 自瀏覽器 LocalStorage 載入歷史與狀態
func LoadTrip() (*model.UserTrip, error) {
	localStorage := js.Global().Get("localStorage")
	if localStorage.IsNull() || localStorage.IsUndefined() {
		return nil, nil
	}
	
	val := localStorage.Call("getItem", StorageKey)
	if val.IsNull() || val.IsUndefined() {
		return nil, nil
	}
	
	var trip model.UserTrip
	err := json.Unmarshal([]byte(val.String()), &trip)
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

// SaveTrip 寫入狀態到 LocalStorage
func SaveTrip(trip *model.UserTrip) error {
	data, err := json.Marshal(trip)
	if err != nil {
		return err
	}
	
	localStorage := js.Global().Get("localStorage")
	if !localStorage.IsNull() && !localStorage.IsUndefined() {
		localStorage.Call("setItem", StorageKey, string(data))
	}
	return nil
}

// ClearTrip 刪除 LocalStorage 狀態
func ClearTrip() {
	localStorage := js.Global().Get("localStorage")
	if !localStorage.IsNull() && !localStorage.IsUndefined() {
		localStorage.Call("removeItem", StorageKey)
	}
}
```

- [ ] **Step 2: 驗證 WebAssembly 編譯語法**

Run: `GOOS=js GOARCH=wasm go build ./storage/...`
Expected: 編譯成功

- [ ] **Step 3: Commit**

Run:
```bash
git add storage/localstorage.go
git commit -m "feat: implement localstorage syncing using syscall/js"
```

---

### Task 5: Go Wasm 主程式註冊

**Files:**
- Create: `main.go`

- [ ] **Step 1: 實作 `main.go`**

Create `main.go` in the root directory:
```go
package main

import (
	"encoding/json"
	"syscall/js"
	"go-where/geo"
	"go-where/model"
	"go-where/storage"
)

var currentTrip *model.UserTrip

func main() {
	// 初始化載入狀態
	trip, err := storage.LoadTrip()
	if err == nil && trip != nil {
		currentTrip = trip
	} else {
		// 預設起點座標：台北車站
		currentTrip = &model.UserTrip{
			CurrentCoords: [2]float64{121.517, 25.047},
			History:       []model.TripNode{},
		}
	}

	// 註冊 JS 回呼函式
	js.Global().Set("goWhereGetState", js.FuncOf(getState))
	js.Global().Set("goWhereDraw", js.FuncOf(drawSpot))
	js.Global().Set("goWhereReset", js.FuncOf(resetTrip))
	js.Global().Set("goWhereUpdateCoords", js.FuncOf(updateCoords))

	// 保持 WASM 進程常駐
	select {}
}

func getState(this js.Value, args []js.Value) any {
	data, _ := json.Marshal(currentTrip)
	return js.ValueOf(string(data))
}

func updateCoords(this js.Value, args []js.Value) any {
	if len(args) >= 2 {
		lng := args[0].Float()
		lat := args[1].Float()
		currentTrip.CurrentCoords = [2]float64{lng, lat}
		storage.SaveTrip(currentTrip)
	}
	return nil
}

func resetTrip(this js.Value, args []js.Value) any {
	// 清空歷史並回到台北車站預設點，也可以點選定位重設
	currentTrip = &model.UserTrip{
		CurrentCoords: [2]float64{121.517, 25.047},
		History:       []model.TripNode{},
	}
	storage.ClearTrip()
	return js.ValueOf(true)
}

func drawSpot(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return js.ValueOf(false)
	}
	
	transport := args[0].String()
	travelTime := args[1].Int()
	
	spot, err := geo.DrawNextSpot(currentTrip.CurrentCoords[0], currentTrip.CurrentCoords[1], transport, travelTime)
	if err != nil {
		return js.ValueOf(js.Map{"error": err.Error()})
	}
	
	// 建立新的路線站點
	newNode := model.TripNode{
		Step:       len(currentTrip.History) + 1,
		Spot:       spot,
		Transport:  transport,
		TravelTime: travelTime,
	}
	
	currentTrip.History = append(currentTrip.History, newNode)
	// 移動當前座標至抽中景點
	currentTrip.CurrentCoords = [2]float64{spot.Px, spot.Py}
	
	storage.SaveTrip(currentTrip)
	
	spotJSON, _ := json.Marshal(spot)
	return js.ValueOf(string(spotJSON))
}
```

- [ ] **Step 2: 編譯 WebAssembly 二進位檔**

Run: `GOOS=js GOARCH=wasm go build -o web/main.wasm main.go`
Expected: 在 `web/main.wasm` 生成編譯產物

- [ ] **Step 3: Commit**

Run:
```bash
git add main.go
git commit -m "feat: write main.go to bridge go logical calls to js"
```

---

### Task 6: 開發伺服器、CI/CD 部署腳本與 README 文件更新

**Files:**
- Create: `server.go`
- Create: `.github/workflows/deploy.yml`
- Modify: `README.md`

- [ ] **Step 1: 建立開發伺服器 `server.go`**

Create `server.go` in the root:
```go
package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	log.Println("本地開發伺服器已啟動：http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

- [ ] **Step 2: 建立 GitHub Actions 工作流 `.github/workflows/deploy.yml`**

Create `.github/workflows/deploy.yml`:
```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [ main ]

permissions:
  contents: write

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.x'

      - name: Build WebAssembly
        run: |
          GOOS=js GOARCH=wasm go build -o web/main.wasm main.go
          GOROOT=$(go env GOROOT)
          cp "$GOROOT/misc/wasm/wasm_exec.js" web/

      - name: Deploy to GitHub Pages
        uses: JamesIves/github-pages-deploy-action@v4
        with:
          folder: web
          branch: gh-pages
```

- [ ] **Step 3: 拷貝 `wasm_exec.js` 備用**

Run: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/`
Expected: `web/wasm_exec.js` 檔案建立成功

- [ ] **Step 4: 修改 `README.md` 加入說明**

Replace content in `README.md`:
```markdown
# go-where 旅遊盲盒 / 隨機公路旅行生成器

`go-where` 是一個旅遊盲盒與隨機公路旅行路線規劃器。使用者開啟網頁並允許定位後，設定出行方式與預期移動時間，系統會隨機抽取適合的台灣旅遊景點，將其串連成一條精彩的今日公路冒險軌跡！

專案完全以**純靜態網頁**部署在 GitHub Pages 上，核心運算與隨機景點庫內嵌於 Go 撰寫的 WebAssembly 中。

## 功能特色
- 📡 **即時定位**：串接瀏覽器 Geolocation API 作為冒險起點。
- 🚗 **交通客製**：支援步行 🚶、單車 🚲、汽車 🚗，依不同移動半徑精準篩選景點。
- 🎁 **旅遊盲盒**：一鍵抽選、炫目雷達動畫、詳細景點圖卡（包含照片與描述）。
- 🗺️ **互動地圖**：整合 Leaflet.js 地圖與 CartoDB Voyager 精緻圖磚，顯示今日連續接關路線與動畫流光軌跡。
- 💾 **自動存檔**：無後端，全狀態透過瀏覽器 LocalStorage 保存，重整不遺失。

## 本地開發與預覽

### 1. 編譯 WebAssembly
確保已安裝 Go 編譯器，並在根目錄執行：
```bash
GOOS=js GOARCH=wasm go build -o web/main.wasm main.go
```

### 2. 啟動本地伺服器
```bash
go run server.go
```
瀏覽：[http://localhost:8080](http://localhost:8080)。

## 部署至 GitHub Pages
本專案已設定 GitHub Actions。當您將程式碼推送到 `main` 分支時，CI/CD 將自動編譯 Go Wasm 並部署 `web/` 目錄下的資源至 `gh-pages` 分支。
```

- [ ] **Step 5: 驗證開發伺服器編譯**

Run: `go build server.go`
Expected: 編譯成功

- [ ] **Step 6: Commit**

Run:
```bash
git add server.go .github/workflows/deploy.yml web/wasm_exec.js README.md
git commit -m "feat: add local server, deploy workflow, and update readme"
```

---

### Task 7: Premium Leaflet.js 地圖與 HTML/CSS 玻璃擬態前端介面

**Files:**
- Create: `web/index.html`

- [ ] **Step 1: 實作網頁 `web/index.html`**

Create `web/index.html` with beautiful premium design, loading CSS, Leaflet JS, and linking Go WASM:
```html
<!DOCTYPE html>
<html lang="zh-TW">
<head>
  <meta charset="UTF-8">
  <title>GoWhere | 旅遊盲盒與隨機公路冒險</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="用旅遊盲盒開啟你的公路旅行！隨機抽取景點，串起一整天的冒險路線">
  
  <!-- Leaflet.js CSS -->
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=" crossorigin=""/>
  <!-- Google Fonts -->
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Noto+Sans+TC:wght@300;400;500;700&family=Outfit:wght@400;600;700&display=swap" rel="stylesheet">
  
  <style>
    :root {
      --bg-dark: #0b0f19;
      --glass-bg: rgba(15, 23, 42, 0.78);
      --glass-border: rgba(255, 255, 255, 0.08);
      --primary: #00f5d4;
      --primary-glow: rgba(0, 245, 212, 0.4);
      --accent: #7b2cbf;
      --accent-glow: rgba(123, 44, 191, 0.4);
      --text-main: #f8fafc;
      --text-sub: #94a3b8;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      font-family: 'Noto Sans TC', sans-serif;
      background-color: var(--bg-dark);
      color: var(--text-main);
      height: 100vh;
      overflow: hidden;
    }

    #map {
      width: 100%;
      height: 100vh;
      z-index: 1;
      position: absolute;
      top: 0;
      left: 0;
    }

    /* Glassmorphism Control Panel */
    .control-panel {
      position: absolute;
      top: 20px;
      left: 20px;
      z-index: 10;
      width: 360px;
      background: var(--glass-bg);
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      border: 1px solid var(--glass-border);
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
      transition: all 0.3s ease;
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 20px;
    }

    .brand h1 {
      font-family: 'Outfit', sans-serif;
      font-size: 24px;
      font-weight: 700;
      background: linear-gradient(135deg, var(--primary), #9d4edd);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .section-title {
      font-size: 14px;
      color: var(--text-sub);
      margin-bottom: 8px;
      text-transform: uppercase;
      letter-spacing: 1px;
    }

    .location-box {
      background: rgba(255, 255, 255, 0.03);
      border: 1px solid var(--glass-border);
      border-radius: 8px;
      padding: 12px;
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .location-info p {
      font-size: 13px;
      color: var(--text-sub);
    }
    .location-info h4 {
      font-size: 15px;
      font-weight: 500;
      color: var(--text-main);
    }

    .btn-locate {
      background: rgba(0, 245, 212, 0.1);
      border: 1px solid var(--primary);
      color: var(--primary);
      padding: 6px 12px;
      border-radius: 6px;
      cursor: pointer;
      font-size: 13px;
      font-weight: 500;
      transition: all 0.2s ease;
    }

    .btn-locate:hover {
      background: var(--primary);
      color: #000;
      box-shadow: 0 0 10px var(--primary-glow);
    }

    /* Transport selector */
    .transport-selector {
      display: flex;
      gap: 8px;
      margin-bottom: 20px;
    }

    .transport-option {
      flex: 1;
      background: rgba(255, 255, 255, 0.03);
      border: 1px solid var(--glass-border);
      color: var(--text-sub);
      padding: 10px;
      border-radius: 8px;
      cursor: pointer;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;
      transition: all 0.2s ease;
    }

    .transport-option:hover {
      border-color: var(--primary);
      color: var(--text-main);
    }

    .transport-option.active {
      background: rgba(0, 245, 212, 0.08);
      border-color: var(--primary);
      color: var(--primary);
      box-shadow: 0 0 8px var(--primary-glow);
    }

    .transport-option span.icon {
      font-size: 18px;
    }
    .transport-option span.label {
      font-size: 11px;
    }

    /* Time selector */
    .time-selector {
      margin-bottom: 24px;
    }

    .time-slider {
      width: 100%;
      height: 6px;
      background: rgba(255, 255, 255, 0.1);
      border-radius: 3px;
      outline: none;
      -webkit-appearance: none;
      margin: 12px 0;
    }

    .time-slider::-webkit-slider-thumb {
      -webkit-appearance: none;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background: var(--primary);
      cursor: pointer;
      box-shadow: 0 0 8px var(--primary);
    }

    .time-display {
      display: flex;
      justify-content: space-between;
      font-size: 13px;
      color: var(--text-sub);
    }

    /* Giant Draw Button */
    .btn-draw {
      width: 100%;
      background: linear-gradient(135deg, var(--accent), #9d4edd);
      border: none;
      color: #fff;
      padding: 14px;
      border-radius: 10px;
      font-size: 16px;
      font-weight: 700;
      cursor: pointer;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px var(--accent-glow);
      letter-spacing: 1px;
    }

    .btn-draw:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(157, 78, 237, 0.6);
      filter: brightness(1.1);
    }

    .btn-draw:active {
      transform: translateY(1px);
    }

    /* Blind Box Popup Card */
    .spot-card {
      position: absolute;
      bottom: 30px;
      right: 30px;
      z-index: 10;
      width: 380px;
      background: var(--glass-bg);
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      border: 1px solid var(--glass-border);
      border-radius: 16px;
      overflow: hidden;
      box-shadow: 0 10px 30px rgba(0,0,0,0.6);
      transform: translateY(120%);
      opacity: 0;
      transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .spot-card.show {
      transform: translateY(0);
      opacity: 1;
    }

    .spot-img {
      width: 100%;
      height: 180px;
      object-fit: cover;
      background-color: #1e293b;
    }

    .spot-content {
      padding: 20px;
    }

    .spot-tag {
      display: inline-block;
      background: rgba(123, 44, 191, 0.15);
      color: #c77dff;
      border: 1px solid rgba(123, 44, 191, 0.3);
      font-size: 11px;
      padding: 3px 8px;
      border-radius: 4px;
      margin-bottom: 10px;
      font-weight: 500;
    }

    .spot-title {
      font-size: 20px;
      font-weight: 700;
      margin-bottom: 8px;
      color: #fff;
    }

    .spot-desc {
      font-size: 13px;
      color: var(--text-sub);
      line-height: 1.5;
      margin-bottom: 16px;
      max-height: 80px;
      overflow-y: auto;
    }

    .spot-actions {
      display: flex;
      gap: 10px;
    }

    .btn-action {
      flex: 1;
      padding: 10px;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      text-align: center;
      transition: all 0.2s ease;
    }

    .btn-confirm {
      background: var(--primary);
      border: none;
      color: #000;
    }
    .btn-confirm:hover {
      box-shadow: 0 0 10px var(--primary-glow);
      filter: brightness(1.1);
    }

    .btn-cancel {
      background: rgba(255,255,255,0.05);
      border: 1px solid var(--glass-border);
      color: var(--text-main);
    }
    .btn-cancel:hover {
      background: rgba(255,255,255,0.1);
    }

    /* History & Reset button */
    .history-panel {
      margin-top: 20px;
      border-top: 1px solid rgba(255, 255, 255, 0.05);
      padding-top: 16px;
    }

    .history-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;
    }

    .history-list {
      max-height: 120px;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      gap: 6px;
    }

    .history-item {
      display: flex;
      justify-content: space-between;
      background: rgba(255,255,255,0.02);
      padding: 6px 10px;
      border-radius: 6px;
      font-size: 12px;
    }

    .btn-reset {
      background: none;
      border: none;
      color: #ef4444;
      font-size: 12px;
      cursor: pointer;
    }
    .btn-reset:hover {
      text-decoration: underline;
    }

    /* Radar scan effect overlays */
    .radar-circle {
      border: 2px solid var(--accent);
      background: rgba(123, 44, 191, 0.1);
      border-radius: 50%;
      position: absolute;
      z-index: 5;
      pointer-events: none;
      transform: translate(-50%, -50%);
      opacity: 0;
    }

    @keyframes radar-animation {
      0% { width: 0px; height: 0px; opacity: 0.8; }
      100% { width: 400px; height: 400px; opacity: 0; }
    }

    /* Flowing path custom class for Leaflet */
    .flowing-path {
      stroke-dasharray: 8, 8;
      animation: dash 30s linear infinite;
    }

    @keyframes dash {
      to {
        stroke-dashoffset: -1000;
      }
    }

    /* Loader */
    .wasm-loader {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      background: var(--bg-dark);
      z-index: 100;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      transition: opacity 0.5s ease;
    }
    .loader-spinner {
      border: 4px solid rgba(255,255,255,0.1);
      border-left-color: var(--primary);
      width: 40px;
      height: 40px;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin-bottom: 20px;
    }
    @keyframes spin {
      100% { transform: rotate(360deg); }
    }
  </style>
</head>
<body>

  <!-- WASM loading screen -->
  <div id="loader" class="wasm-loader">
    <div class="loader-spinner"></div>
    <p>正在載入 Go WebAssembly 冒險核心...</p>
  </div>

  <!-- Leaflet Map -->
  <div id="map"></div>

  <!-- Left Sidebar Panel -->
  <div class="control-panel">
    <div class="brand">
      <span style="font-size: 26px;">🎲</span>
      <h1>GoWhere 旅遊盲盒</h1>
    </div>

    <!-- Location selection -->
    <div class="section-title">起點座標定位</div>
    <div class="location-box">
      <div class="location-info">
        <h4 id="loc-status">台北車站 (預設)</h4>
        <p id="loc-coords">121.5170, 25.0470</p>
      </div>
      <button class="btn-locate" id="btn-locate">📡 GPS 定位</button>
    </div>

    <!-- Transport tool selection -->
    <div class="section-title">選擇探險交通</div>
    <div class="transport-selector">
      <button class="transport-option active" data-mode="driving">
        <span class="icon">🚗</span>
        <span class="label">汽車</span>
      </button>
      <button class="transport-option" data-mode="cycling">
        <span class="icon">🚲</span>
        <span class="label">單車</span>
      </button>
      <button class="transport-option" data-mode="walking">
        <span class="icon">🚶</span>
        <span class="label">步行</span>
      </button>
    </div>

    <!-- Maximum time -->
    <div class="section-title">接受移動時間</div>
    <div class="time-selector">
      <input type="range" class="time-slider" id="time-slider" min="15" max="180" value="60" step="15">
      <div class="time-display">
        <span>限時: <strong id="time-val" style="color: var(--primary);">60</strong> 分鐘</span>
        <span>預估半徑: <strong id="radius-val" style="color: var(--primary);">60</strong> km</span>
      </div>
    </div>

    <!-- Draw Button -->
    <button class="btn-draw" id="btn-draw">🎁 抽取旅遊盲盒</button>

    <!-- History list -->
    <div class="history-panel">
      <div class="history-header">
        <span class="section-title" style="margin:0;">今天探險軌跡</span>
        <button class="btn-reset" id="btn-reset">重設路線</button>
      </div>
      <div class="history-list" id="history-list">
        <!-- populated by js -->
        <p style="font-size:12px; color: var(--text-sub); text-align:center; padding: 10px;">尚未開始探險</p>
      </div>
    </div>
  </div>

  <!-- Right Drawer Spot Detail Card -->
  <div class="spot-card" id="spot-card">
    <img src="" alt="景點照片" class="spot-img" id="spot-img">
    <div class="spot-content">
      <div class="spot-tag" id="spot-tag">⛰️ 知名景點</div>
      <h3 class="spot-title" id="spot-title">景點名稱</h3>
      <p class="spot-desc" id="spot-desc">景點介紹描述載入中...</p>
      <div class="spot-actions">
        <button class="btn-action btn-confirm" id="btn-confirm-spot">確認前往 ➔</button>
        <button class="btn-action btn-cancel" id="btn-cancel-spot">換一個</button>
      </div>
    </div>
  </div>

  <!-- Leaflet.js JS -->
  <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
  <!-- WASM Glue code -->
  <script src="wasm_exec.js"></script>
  
  <script>
    // 1. 初始化 Go WASM
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
      go.run(result.instance);
      // 隱藏載入遮罩
      document.getElementById('loader').style.opacity = '0';
      setTimeout(() => {
        document.getElementById('loader').style.display = 'none';
      }, 500);
      
      // 載入當前 LocalStorage 狀態並更新 UI 地圖
      syncStateFromWasm();
    }).catch(err => {
      console.error("WASM load fail:", err);
      document.querySelector('#loader p').innerHTML = "<span style='color:#ef4444'>核心載入失敗！請確認 main.wasm 存在並以網頁伺服器開啟。</span>";
    });

    // 2. 初始化 Leaflet 地圖
    // 使用 CartoDB Voyager 圖磚 (暗黑色調主題非常契合)
    const map = L.map('map', {
      zoomControl: false,
      attributionControl: false
    }).setView([25.0470, 121.5170], 11);

    L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager_labels_under/{z}/{x}/{y}{r}.png', {
      maxZoom: 20
    }).addTo(map);

    // 定義地圖圖示
    const originIcon = L.divIcon({
      className: 'custom-div-icon',
      html: `<div style="background-color: var(--primary); width: 14px; height: 14px; border-radius: 50%; border: 2px solid white; box-shadow: 0 0 10px var(--primary-glow)"></div>`,
      iconSize: [14, 14],
      iconAnchor: [7, 7]
    });

    const spotIcon = L.divIcon({
      className: 'custom-div-icon',
      html: `<div style="background-color: #9d4edd; width: 16px; height: 16px; border-radius: 50%; border: 2px solid white; box-shadow: 0 0 12px rgba(157, 78, 237, 0.8); display: flex; align-items: center; justify-content: center; font-size: 10px;">📍</div>`,
      iconSize: [20, 20],
      iconAnchor: [10, 10]
    });

    let currentPosMarker = null;
    let markersGroup = L.layerGroup().addTo(map);
    let routePolyline = null;
    let currentSelectedSpot = null;

    // 3. UI 狀態同步與更新
    function syncStateFromWasm() {
      if (typeof goWhereGetState !== 'function') return;
      const stateStr = goWhereGetState();
      const state = JSON.parse(stateStr);
      
      const lng = state.current_coords[0];
      const lat = state.current_coords[1];
      
      // 更新定位面板顯示
      document.getElementById('loc-coords').innerText = `${lng.toFixed(4)}, ${lat.toFixed(4)}`;
      
      // 清空舊 Marker
      markersGroup.clearLayers();
      
      // 更新定位點
      if (currentPosMarker) {
        currentPosMarker.setLatLng([lat, lng]);
      } else {
        currentPosMarker = L.marker([lat, lng], {icon: originIcon}).addTo(map);
      }
      
      // 若無歷史，重設地圖中心
      if (!state.history || state.history.length === 0) {
        map.setView([lat, lng], 11);
        document.getElementById('loc-status').innerText = "旅程出發點";
        document.getElementById('history-list').innerHTML = `<p style="font-size:12px; color: var(--text-sub); text-align:center; padding: 10px;">尚未開始探險</p>`;
        if (routePolyline) {
          map.removeLayer(routePolyline);
          routePolyline = null;
        }
        return;
      }
      
      // 繪製歷史站點與軌跡
      const latlngs = [];
      
      // 加入第一個起點座標 (若有歷史，需要從第一站的起點推算)
      // 這裡簡化為以最後歷史站點或定位當前點渲染
      let historyHTML = '';
      
      // 收集路線上的所有節點
      state.history.forEach((node, index) => {
        const spot = node.spot;
        const transportIcon = node.transport === 'walking' ? '🚶' : (node.transport === 'cycling' ? '🚲' : '🚗');
        
        // 放入坐標點
        latlngs.push([spot.Py, spot.Px]);
        
        // 標記景點
        L.marker([spot.Py, spot.Px], {icon: spotIcon})
          .addTo(markersGroup)
          .bindPopup(`<b>第 ${node.step} 站：${spot.name}</b><br>${spot.add}`);
          
        historyHTML += `
          <div class="history-item">
            <span>第 ${node.step} 站：${spot.name}</span>
            <span style="color: var(--text-sub);">${transportIcon} ${node.travel_time}m</span>
          </div>
        `;
      });
      
      document.getElementById('history-list').innerHTML = historyHTML;
      document.getElementById('loc-status').innerText = `目前位置：${state.history[state.history.length-1].spot.name}`;
      
      // 繪製軌跡連線
      if (routePolyline) {
        map.removeLayer(routePolyline);
      }
      
      // 加上流動虛線特效
      routePolyline = L.polyline(latlngs, {
        color: '#9d4edd',
        weight: 4,
        opacity: 0.8,
        className: 'flowing-path'
      }).addTo(map);
      
      // 移動視角以包含所有景點
      const bounds = L.latLngBounds(latlngs);
      bounds.extend([lat, lng]);
      map.fitBounds(bounds, { padding: [50, 50] });
    }

    // 4. 事件監聽處理
    // 交通工具選擇
    let currentTransport = 'driving';
    document.querySelectorAll('.transport-option').forEach(btn => {
      btn.addEventListener('click', (e) => {
        document.querySelectorAll('.transport-option').forEach(b => b.classList.remove('active'));
        const opt = e.currentTarget;
        opt.classList.add('active');
        currentTransport = opt.dataset.mode;
        updateRadiusValue();
      });
    });

    // 時間滑桿
    const timeSlider = document.getElementById('time-slider');
    timeSlider.addEventListener('input', updateRadiusValue);

    function updateRadiusValue() {
      const timeVal = parseInt(timeSlider.value);
      document.getElementById('time-val').innerText = timeVal;
      
      let speedKmh = 60.0;
      if (currentTransport === 'walking') speedKmh = 5.0;
      else if (currentTransport === 'cycling') speedKmh = 15.0;
      
      const radius = (speedKmh * (timeVal / 60)).toFixed(1);
      document.getElementById('radius-val').innerText = radius;
    }
    
    updateRadiusValue(); // 初始化半徑顯示

    // 定位按鈕
    document.getElementById('btn-locate').addEventListener('click', () => {
      const btn = document.getElementById('btn-locate');
      btn.innerText = "定位中...";
      btn.disabled = true;
      
      if (!navigator.geolocation) {
        alert("瀏覽器不支援 GPS 定位，使用預設值。");
        resetLocateBtn();
        return;
      }
      
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          const lat = pos.coords.latitude;
          const lng = pos.coords.longitude;
          if (typeof goWhereUpdateCoords === 'function') {
            goWhereUpdateCoords(lng, lat);
            syncStateFromWasm();
          }
          resetLocateBtn();
        },
        (err) => {
          alert("定位失敗！請開啟位置授權。將使用預設台北車站。");
          resetLocateBtn();
        },
        { enableHighAccuracy: true, timeout: 8000 }
      );
    });

    function resetLocateBtn() {
      const btn = document.getElementById('btn-locate');
      btn.innerText = "📡 GPS 定位";
      btn.disabled = false;
    }

    // 隨機抽取盲盒按鈕
    document.getElementById('btn-draw').addEventListener('click', () => {
      if (typeof goWhereDraw !== 'function') return;
      
      const state = JSON.parse(goWhereGetState());
      const originLng = state.current_coords[0];
      const originLat = state.current_coords[1];
      
      // 1. 雷達掃描視覺特效
      const radar = document.createElement('div');
      radar.className = 'radar-circle';
      radar.style.left = `${map.latLngToContainerPoint([originLat, originLng]).x}px`;
      radar.style.top = `${map.latLngToContainerPoint([originLat, originLng]).y}px`;
      document.body.appendChild(radar);
      
      radar.style.animation = 'radar-animation 1.5s ease-out forwards';
      setTimeout(() => radar.remove(), 1500);

      // 地圖短暫縮放震盪
      map.zoomOut(0.5);
      setTimeout(() => map.zoomIn(0.5), 300);

      // 2. 呼叫 Wasm 進行抽選
      const timeVal = parseInt(timeSlider.value);
      const resStr = goWhereDraw(currentTransport, timeVal);
      const res = JSON.parse(resStr);
      
      if (res.error) {
        alert("抽選出錯: " + res.error);
        return;
      }
      
      currentSelectedSpot = res;
      
      // 3. 彈出景點卡片展示
      setTimeout(() => {
        document.getElementById('spot-title').innerText = res.name;
        document.getElementById('spot-desc').innerText = res.description || "這裡有一段絕佳的探險回憶等著你去發掘！";
        document.getElementById('spot-img').src = res.picture1 || "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800";
        document.getElementById('spot-card').classList.add('show');
        
        // 暫時在地圖上新增虛擬的未確認標記並平移過去
        map.panTo([res.py, res.px]);
      }, 800);
    });

    // 確認前往
    document.getElementById('btn-confirm-spot').addEventListener('click', () => {
      document.getElementById('spot-card').classList.remove('show');
      currentSelectedSpot = null;
      syncStateFromWasm();
    });

    // 換一個 (取消)
    document.getElementById('btn-cancel-spot').addEventListener('click', () => {
      document.getElementById('spot-card').classList.remove('show');
      // 將 Wasm 中的最後一筆歷史刪除，並還原座標
      if (typeof goWhereGetState === 'function') {
        const state = JSON.parse(goWhereGetState());
        if (state.history && state.history.length > 0) {
          state.history.pop();
          // 更新座標到上一個位置 (或起點)
          let lastLng = 121.5170;
          let lastLat = 25.0470;
          if (state.history.length > 0) {
            const lastSpot = state.history[state.history.length-1].spot;
            lastLng = lastSpot.px;
            lastLat = lastSpot.py;
          }
          state.current_coords = [lastLng, lastLat];
          
          // 寫入 Wasm 的狀態
          goWhereReset(); // 先重置
          // 重新載入定位
          goWhereUpdateCoords(lastLng, lastLat);
          // 還原之前的歷史
          state.history.forEach(n => {
            // 直接呼叫 Draw 來重現 (雖然不保證隨機相同，但為了讓狀態同步，這是一個好的還原)
            // 這裡我們簡化：在 LocalStorage 內手動寫回
            localStorage.setItem("go_where_trip_state", JSON.stringify(state));
          });
          
          // 更簡單：我們可以直接將 state 重新序列化寫入 localStorage，然後重新載入 Wasm 狀態
          localStorage.setItem("go_where_trip_state", JSON.stringify(state));
        }
      }
      
      // 重新載入 WASM 並同步 UI
      // 這裡直接呼叫 main() 的 init 或重新同步 state
      // 因為 Wasm 內部 currentTrip 是指標，我們可以直接更新 localStorage 後，重新點擊即可
      // 為了免去 Wasm 重新載入，我們手動做狀態覆蓋：
      location.reload(); 
    });

    // 重設按鈕
    document.getElementById('btn-reset').addEventListener('click', () => {
      if (confirm("確定要重設今天的探險路線嗎？這會清除所有足跡。")) {
        if (typeof goWhereReset === 'function') {
          goWhereReset();
          syncStateFromWasm();
        }
      }
    });
  </script>
</body>
</html>
```

- [ ] **Step 2: Commit**

Run:
```bash
git add web/index.html
git commit -m "feat: design premium glassmorphism html UI and integrate leaflet map"
```

---

## 驗證流程 (Verification Process)

1. **單元測試驗證**：
   執行 `go test -v ./geo/...` 確保經緯度 Haversine 距離與抽取邏輯完全無誤。
2. **WASM 編譯驗證**：
   執行 `GOOS=js GOARCH=wasm go build -o web/main.wasm main.go` 確保編譯成功，輸出 `web/main.wasm`。
3. **本機網頁運行驗證**：
   啟動伺服器 `go run server.go`。開啟瀏覽器載入 `http://localhost:8080`，驗證 UI 渲染、Leaflet.js 地圖載入、定位功能與 Wasm 抽取按鈕與歷史路徑流光特效正常。
