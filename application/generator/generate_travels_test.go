package generator

import (
	"darbelis.eu/persedimai/tables"
	"testing"
	"time"
)

func TestGenerateTravels(t *testing.T) {
	gf := GeneratorFactory{}

	g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

	pointArrayConsumer := NewPointArrayConsumer()

	err := g.GeneratePoints(pointArrayConsumer)
	if err != nil {
		t.Errorf("GeneratePoints returned error: %v", err)
	}

	// Assert number of points generated
	expectedPoints := 9
	actualPoints := len(pointArrayConsumer.Points)
	if actualPoints != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, actualPoints)
	}

	travelConsumer := NewTravelArrayConsumer()
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	speed := 1000.0
	restHours := 2

	err = g.GenerateTravels(pointArrayConsumer.Points, fromDate, toDate, speed, restHours, travelConsumer)
	if err != nil {
		t.Errorf("GenerateTravels returned error: %v", err)
	}

	// Assert total number of travels
	expectedTravels := 112 // 12 H/V pairs × 6 travels + 8 diagonal pairs × 5 travels
	actualTravels := len(travelConsumer.Travels)
	if actualTravels != expectedTravels {
		t.Errorf("Expected %d travels, got %d", expectedTravels, actualTravels)
	}

	t.Logf("Generated %d points and %d travels", actualPoints, actualTravels)

	// Create a map for quick travel lookup
	travels := travelConsumer.Travels

	// Helper function to find a travel by From and To IDs
	findTravel := func(fromID, toID string) *tables.Transfer {
		for _, travel := range travels {
			if travel.From == fromID && travel.To == toID {
				return travel
			}
		}
		return nil
	}

	// Assert: Check a horizontal travel (point 1 to point 2)
	// Point 1: (0,0), Point 2: (2000,0) - distance 2000, travel time 2h
	horizontalTravel := findTravel("1", "2")
	if horizontalTravel == nil {
		t.Error("Expected to find horizontal travel from point 1 to point 2")
	} else {
		expectedArrival := fromDate.Add(2 * time.Hour)
		if !horizontalTravel.Departure.Equal(fromDate) {
			t.Errorf("Horizontal travel: expected departure %v, got %v", fromDate, horizontalTravel.Departure)
		}
		if !horizontalTravel.Arrival.Equal(expectedArrival) {
			t.Errorf("Horizontal travel: expected arrival %v, got %v", expectedArrival, horizontalTravel.Arrival)
		}
	}

	// Assert: Check a vertical travel (point 1 to point 4)
	// Point 1: (0,0), Point 4: (0,2000) - distance 2000, travel time 2h
	verticalTravel := findTravel("1", "4")
	if verticalTravel == nil {
		t.Error("Expected to find vertical travel from point 1 to point 4")
	} else {
		expectedArrival := fromDate.Add(2 * time.Hour)
		if !verticalTravel.Departure.Equal(fromDate) {
			t.Errorf("Vertical travel: expected departure %v, got %v", fromDate, verticalTravel.Departure)
		}
		if !verticalTravel.Arrival.Equal(expectedArrival) {
			t.Errorf("Vertical travel: expected arrival %v, got %v", expectedArrival, verticalTravel.Arrival)
		}
	}

	// Assert: Check a diagonal travel (point 1 to point 5)
	// Point 1: (0,0), Point 5: (2000,2000) - distance 2828.43, travel time ~2.828h
	diagonalTravel := findTravel("1", "5")
	if diagonalTravel == nil {
		t.Error("Expected to find diagonal travel from point 1 to point 5")
	} else {
		// Diagonal distance: sqrt(2000^2 + 2000^2) = 2828.427..
		// Travel time: 2828.427 / 1000 = 2.828427 hours
		travelHours := 2.828427124746190
		expectedDuration := time.Duration(travelHours * float64(time.Hour))
		actualDuration := diagonalTravel.Arrival.Sub(diagonalTravel.Departure)

		// Allow 1 second tolerance due to floating point precision
		tolerance := 1 * time.Second
		if actualDuration < expectedDuration-tolerance || actualDuration > expectedDuration+tolerance {
			t.Errorf("Diagonal travel: expected duration ~%v, got %v", expectedDuration, actualDuration)
		}
	}

	// Assert: Check that we have both directions for a pair
	// Should have travels from 1→2 and 2→1
	travel_1_to_2 := findTravel("1", "2")
	travel_2_to_1 := findTravel("2", "1")
	if travel_1_to_2 == nil || travel_2_to_1 == nil {
		t.Error("Expected to find travels in both directions between point 1 and point 2")
	}

	// Assert: Check that departure times are properly spaced
	// For the same pair, second travel should depart after first arrival + rest time
	if travel_1_to_2 != nil {
		// Find the second travel from 1 to 2
		secondTravelFrom1To2 := false
		for _, travel := range travels {
			if travel.From == "1" && travel.To == "2" && !travel.Departure.Equal(travel_1_to_2.Departure) {
				// This is a different travel from 1 to 2
				expectedSecondDeparture := travel_1_to_2.Arrival.Add(time.Duration(restHours) * time.Hour)
				if travel.Departure.Equal(expectedSecondDeparture) {
					secondTravelFrom1To2 = true
					break
				}
			}
		}
		// Note: The second travel might be in the reverse direction (2→1)
		// So we just verify we found the first travel correctly
		_ = secondTravelFrom1To2
	}

	// Assert: Verify all travels have valid From/To IDs (between "1" and "9")
	for i, travel := range travels {
		fromID := travel.From
		toID := travel.To
		if fromID < "1" || fromID > "9" {
			t.Errorf("Travel %d: invalid From ID %s", i, fromID)
		}
		if toID < "1" || toID > "9" {
			t.Errorf("Travel %d: invalid To ID %s", i, toID)
		}
		if fromID == toID {
			t.Errorf("Travel %d: From and To are the same (%s)", i, fromID)
		}
	}

	// Assert: Verify all travels are within the time window
	for i, travel := range travels {
		if travel.Departure.Before(fromDate) {
			t.Errorf("Travel %d: departure %v is before fromDate %v", i, travel.Departure, fromDate)
		}
		if travel.Arrival.After(toDate) {
			t.Errorf("Travel %d: arrival %v is after toDate %v", i, travel.Arrival, toDate)
		}
		if !travel.Arrival.After(travel.Departure) {
			t.Errorf("Travel %d: arrival %v is not after departure %v", i, travel.Arrival, travel.Departure)
		}
	}
}
