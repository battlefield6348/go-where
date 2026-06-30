# go-where

`go-where` 是一個旅遊盲盒 / 隨機公路旅行生成器。使用者開啟網頁、提供定位、設定交通工具與可接受的移動時間後，系統會從可抵達的台灣觀光景點中隨機抽出下一站，並把整天的移動軌跡串成一條冒險路線。

專案預計以純靜態方式部署在 GitHub Pages，核心邏輯使用 Go 編譯成 WebAssembly，在瀏覽器內完成景點篩選、旅程狀態保存與地圖互動。

## 文件導覽

- [開發需求](./docs/requirements.md)
- [產品規格](./docs/product-spec.md)
- [專案結構](./docs/project-structure.md)
- [開發教訓](./PROJECT_LESSONS.md)
