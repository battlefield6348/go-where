# 專案結構

```text
go-where/
├── main.go
├── model/
│   └── trip.go
├── geo/
│   └── calculator.go
├── storage/
│   └── localstorage.go
└── web/
    ├── index.html
    └── main.wasm
```

## 目錄用途

- `main.go`：Wasm 入口，負責綁定 JavaScript 函式。
- `model/`：核心資料模型，例如景點與旅程狀態。
- `geo/`：地理計算與景點篩選。
- `storage/`：瀏覽器 LocalStorage 存取。
- `web/`：GitHub Pages 需要的靜態資產。

## 初始化原則

- 先建立最小骨架，不預先放入未使用的工具層與抽象。
- 邏輯檔案等真正開始實作時再新增，避免空殼程式碼。
- `web/main.wasm` 屬於編譯產物，初始化階段不預先提交。
