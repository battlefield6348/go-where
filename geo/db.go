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
