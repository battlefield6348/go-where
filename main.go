//go:build js && wasm

package main

import (
	"encoding/json"
	"go-where/geo"
	"go-where/model"
	"go-where/storage"
	"syscall/js"
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
		errMap := map[string]string{"error": err.Error()}
		errJSON, _ := json.Marshal(errMap)
		return js.ValueOf(string(errJSON))
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
