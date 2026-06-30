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
go run cmd/server/main.go
```
瀏覽：[http://localhost:8080](http://localhost:8080)。

## 部署至 GitHub Pages
本專案已設定 GitHub Actions。當您將程式碼推送到 `main` 分支時，CI/CD 將自動編譯 Go Wasm 並部署 `web/` 目錄下的資源至 `gh-pages` 分支。
