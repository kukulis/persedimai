package generator

import (
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"testing"
	"time"
)

func TestGenerateTravels_WithHubs(t *testing.T) {
	gf := GeneratorFactory{}

	g := gf.CreateGenerator(7, 1000, 0, &SimpleIdGenerator{})

	pointArrayConsumer := NewPointArrayConsumer()

	err := g.GeneratePoints(pointArrayConsumer)
	if err != nil {
		t.Errorf("GeneratePoints returned error: %v", err)
	}

	// Assert number of points generated
	expectedPoints := 16
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
	expectedTravels := 234 // 12 H/V pairs × 6 travels + 8 diagonal pairs × 5 travels
	actualTravels := len(travelConsumer.Travels)
	if actualTravels != expectedTravels {
		t.Errorf("Expected %d travels, got %d", expectedTravels, actualTravels)
	}

	t.Logf("Generated %d points and %d travels", actualPoints, actualTravels)

	//g.GenerateTravelsBetweenHubPoints()

	allPoints := pointArrayConsumer.Points
	hubPoints := util.ArrayFilter(allPoints, func(point *tables.Point) bool { return point.X == 0 })

	if len(hubPoints) != 4 {
		t.Errorf("Expected %d hub points, got %d", 4, len(hubPoints))
	}

	hubTravelsConsumer := NewTravelArrayConsumer()
	hubSourceTravelsConsumer := NewTravelArrayConsumer()

	// The dates period might be too short
	err = g.GenerateTravelsBetweenHubPoints(hubPoints, fromDate, toDate, speed, restHours, hubTravelsConsumer)
	if err != nil {
		t.Errorf("GenerateTravelsBetweenHubPoints returned error: %v", err)
	}
	err = g.GenerateTravelsFromHubToNonHubPoints(hubPoints, allPoints, fromDate, toDate, speed, restHours, hubSourceTravelsConsumer)
	if err != nil {
		t.Errorf("GenerateTravels returned error: %v", err)
	}

	if len(hubTravelsConsumer.Travels) < 10 {
		t.Errorf("Expected at least %d travels between hub points, got %d", 10, len(hubTravelsConsumer.Travels))
	}

	if len(hubTravelsConsumer.Travels) < 16 {
		t.Errorf("Expected at least %d travels between hub and non hub points, got %d", 16, len(hubSourceTravelsConsumer.Travels))
	}
}
