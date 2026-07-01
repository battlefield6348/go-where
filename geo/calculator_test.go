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

func TestDrawNextSpotWithDuplicateExclusions(t *testing.T) {
	// 假設排除名單為 [ "1", "1", ..., "1" ] 長度大於等於所有的景點數 (56)
	// 但實際去重後只有 1 個景點。
	// 這時不應該過早清空去重機制，台北 101 (ID: 1) 依然要被排除。
	excluded := make([]string, 60)
	for i := range excluded {
		excluded[i] = "1"
	}
	
	for i := 0; i < 20; i++ {
		spot, _, err := DrawNextSpot(121.517, 25.047, "driving", 60, excluded)
		if err != nil {
			t.Fatalf("Failed to draw spot: %s", err)
		}
		if spot.ID == "1" {
			t.Error("Should not draw excluded spot台北 101 (ID 1) even with large duplicate exclusions list")
		}
	}
}

func TestDrawNextSpotWrapAroundNotDuplicate(t *testing.T) {
	// 當排除名單包含了所有景點，此時會觸發 wrap-around (清空排除清單)。
	// 為了避免下一抽直接抽回當前所在的最後一站景點，我們將當前位置設在最後一站 (假設是幾米主題廣場, ID: 56, 座標 121.7578, 24.7525)。
	// 此時不論怎麼抽，都不應該抽中 ID: 56。
	excluded := make([]string, len(CachedSpots))
	for i, spot := range CachedSpots {
		excluded[i] = spot.ID
	}
	
	// 最後一個確認的景點是 ID: 56
	lastSpot := CachedSpots[len(CachedSpots)-1] // 假設 ID 是 "56"
	
	// 在幾米主題廣場 (121.7578, 24.7525) 抽下一站，限速 60 分鐘開車 (radiusKm 很大，所有景點都能到)
	// 抽 100 次，驗證沒有一次會是 lastSpot.ID
	for i := 0; i < 100; i++ {
		spot, _, err := DrawNextSpot(lastSpot.Px, lastSpot.Py, "driving", 60, excluded)
		if err != nil {
			t.Fatalf("Failed to draw spot: %s", err)
		}
		if spot.ID == lastSpot.ID {
			t.Fatalf("Should not draw the same current spot %s (ID %s) immediately after wrap-around", lastSpot.Name, lastSpot.ID)
		}
	}
}

