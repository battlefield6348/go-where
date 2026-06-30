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
