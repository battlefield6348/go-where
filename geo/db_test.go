package geo

import (
	"testing"
)

func TestCachedSpotsLoaded(t *testing.T) {
	if len(CachedSpots) == 0 {
		t.Fatal("Expected CachedSpots to contain spots, but it was empty")
	}
	
	firstSpot := CachedSpots[0]
	if firstSpot.Name == "" {
		t.Error("Expected spot to have a name, got empty string")
	}
	if firstSpot.Px == 0 || firstSpot.Py == 0 {
		t.Errorf("Expected spot %s to have valid coordinates, got Px:%f Py:%f", firstSpot.Name, firstSpot.Px, firstSpot.Py)
	}
}
