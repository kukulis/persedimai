package travel_finder

import (
	"darbelis.eu/persedimai/internal/data"
	"testing"
	"time"
)

func TestSimpleTravelSearchStrategy_FindPath_InvalidTravelCount(t *testing.T) {
	strategy := NewSimpleTravelSearchStrategy(nil)

	t.Run("TravelCount=0", func(t *testing.T) {
		filter := data.NewTravelFilter("", "", time.Time{}, time.Time{}, 0)
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=0, got nil")
		}
		if err.Error() != "invalid TravelCount: must be 1, 2, or 3" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("TravelCount=4", func(t *testing.T) {
		filter := data.NewTravelFilter("", "", time.Time{}, time.Time{}, 4)
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=4, got nil")
		}
		if err.Error() != "unimplemented: TravelCount > 3 not supported" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("TravelCount=5", func(t *testing.T) {
		filter := data.NewTravelFilter("", "", time.Time{}, time.Time{}, 5)
		_, err := strategy.FindPaths(filter)
		if err == nil {
			t.Error("Expected error for TravelCount=5, got nil")
		}
		if err.Error() != "unimplemented: TravelCount > 3 not supported" {
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
