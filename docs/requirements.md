# 開發需求

這份文件集中記錄目前的實作需求、技術約束與開發順序。

## 功能需求

- 使用者可授權 GPS 定位，作為旅程起點。
- 使用者可設定交通工具與移動時間。
- 系統需從可抵達的真實景點中隨機抽選下一站，並在地圖上標示。
- 使用者抵達景點後，可從該景點繼續抽選下一站。
- 整段旅程需保留歷史軌跡，讓一天的路線能持續累積。

## 技術約束

- 部署環境為 GitHub Pages，不能依賴後端服務或資料庫。
- 核心邏輯使用 Golang，並編譯成 WebAssembly 在瀏覽器執行。
- 地圖使用 Leaflet.js 與免費圖磚。
- 旅程狀態透過瀏覽器 LocalStorage 保存。
- 景點資料來自台灣交通部觀光署開放資料，並以 `//go:embed` 內嵌進 Wasm。

## 預計實作順序

1. `geo/calculator.go`：實作 Haversine 距離計算與半徑篩選。
2. `storage/localstorage.go`：實作 `UserTrip` 的 LocalStorage 讀寫。
3. `main.go`：提供 JavaScript 可呼叫的 Wasm 入口。
4. `web/index.html`：整合地圖與 Wasm 橋接。

## 延伸文件

- [產品規格](./product-spec.md)：產品目標、資料模型與演算法需求。
- [專案結構](./project-structure.md)：目錄用途與初始化原則。
