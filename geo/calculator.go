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

// DrawNextSpot 根據當前座標、交通工具及時間篩選景點，並隨機抽取下一站，且會自動排除 excludedIDs 中已造訪的景點
// 傳回值：景點、是否為超出範圍的備選景點 (isFallback)、錯誤
func DrawNextSpot(currentLng, currentLat float64, transport string, travelTimeMinutes int, excludedIDs []string) (model.TouristSpot, bool, error) {
	if len(CachedSpots) == 0 {
		return model.TouristSpot{}, false, errors.New("spots database is empty")
	}

	// 建立已去過景點的 Map 用於快速過濾
	excludedMap := make(map[string]bool)
	for _, id := range excludedIDs {
		excludedMap[id] = true
	}

	// 若使用者已造訪所有景點，則清空排除清單以防沒有景點可抽
	if len(excludedIDs) >= len(CachedSpots) {
		excludedMap = make(map[string]bool)
	}

	// 依交通工具決定平均時速 (km/h)
	var speedKmh float64
	switch transport {
	case "walking":
		speedKmh = 5.0
	case "cycling":
		speedKmh = 15.0
	case "scooter":
		speedKmh = 40.0
	case "transit":
		speedKmh = 30.0
	default: // driving
		speedKmh = 60.0
	}
	
	// 移動半徑 = 時速 * (時間分鐘 / 60)
	radiusKm := speedKmh * (float64(travelTimeMinutes) / 60.0)
	
	// 篩選半徑內的候選景點 (排除已去過的景點)
	var candidates []model.TouristSpot
	for _, spot := range CachedSpots {
		if excludedMap[spot.ID] {
			continue
		}
		dist := Haversine(currentLng, currentLat, spot.Px, spot.Py)
		if dist <= radiusKm {
			candidates = append(candidates, spot)
		}
	}
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 若半徑內沒有任何景點，退而求其次，擴大尋找最近的三個景點並隨機抽一個 (同樣排除已去過的景點)
	if len(candidates) == 0 {
		type distTuple struct {
			spot model.TouristSpot
			dist float64
		}
		
		var sorted []distTuple
		for _, spot := range CachedSpots {
			if excludedMap[spot.ID] {
				continue
			}
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
			return sorted[chosenIdx].spot, true, nil
		}
		
		return model.TouristSpot{}, false, errors.New("no spots available")
	}
	
	return candidates[r.Intn(len(candidates))], false, nil
}
