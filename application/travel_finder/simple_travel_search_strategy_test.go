package travel_finder

import (
	"darbelis.eu/persedimai/data"
	"testing"
)

func TestSimpleTravelSearchStrategy_FindPath_InvalidTravelCount(t *testing.T) {
	strategy := NewSimpleTravelSearchStrategy(nil)

	t.Run("TravelCount=0", func(t *testing.T) {
		filter := &data.TravelFilter{TravelCount: 0}
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=0, got nil")
		}
		if err.Error() != "invalid TravelCount: must be 1 or 2" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("TravelCount=3", func(t *testing.T) {
		filter := &data.TravelFilter{TravelCount: 3}
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=3, got nil")
		}
		if err.Error() != "unimplemented: TravelCount > 2 not supported" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("TravelCount=5", func(t *testing.T) {
		filter := &data.TravelFilter{TravelCount: 5}
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=5, got nil")
		}
		if err.Error() != "unimplemented: TravelCount > 2 not supported" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
}

func TestSimpleTravelSearchStrategy_GetName(t *testing.T) {
	strategy := NewSimpleTravelSearchStrategy(nil)
	if strategy.GetName() != "Simple" {
		t.Errorf("Expected name 'Simple', got '%s'", strategy.GetName())
	}
}
