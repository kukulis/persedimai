package generator

import (
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

	travelConsumer := NewTravelArrayConsumer()
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	speed := 1000.0
	restHours := 2

	err = g.GenerateTravels(pointArrayConsumer.Points, fromDate, toDate, speed, restHours, travelConsumer)
	if err != nil {
		t.Errorf("GenerateTravels returned error: %v", err)
	}

	// Verify that some travels were generated
	if len(travelConsumer.Travels) == 0 {
		t.Errorf("Expected some travels to be generated, got 0")
	}

	// Log the number of travels for visibility
	t.Logf("Generated %d travels", len(travelConsumer.Travels))
}
