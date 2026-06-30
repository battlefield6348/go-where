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
	// 當前起點台北車站附近，限速 60 分鐘開車應能抽中附近景點
	spot, _, err := DrawNextSpot(121.517, 25.047, "driving", 60, nil)
	if err != nil {
		t.Fatalf("Failed to draw spot: %s", err)
	}
	if spot.Name == "" {
		t.Error("Drawn spot name is empty")
	}
}

func TestDrawNextSpotWithExclusions(t *testing.T) {
	// 設想我們在台北車站，若排除台北 101 (ID: 1)
	// 應該不能抽中台北 101
	excluded := []string{"1"}
	for i := 0; i < 20; i++ {
		spot, _, err := DrawNextSpot(121.517, 25.047, "driving", 60, excluded)
		if err != nil {
			t.Fatalf("Failed to draw spot: %s", err)
		}
		if spot.ID == "1" {
			t.Error("Should not draw excluded spot台北 101 (ID 1)")
		}
	}
}
