package generator

import (
	"darbelis.eu/persedimai/tables"
	"math"
	"testing"
)

func TestGenerateSingleTravel(t *testing.T) {
	gf := &GeneratorFactory{}
	g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

	// Two points 10,000 units apart on X axis
	pointA := tables.Point{ID: "1", X: 0, Y: 0}
	pointB := tables.Point{ID: "2", X: 10000, Y: 0}

	// Speed: 1000 units per hour
	// Distance: 10,000 units
	// Expected travel time: 10 hours
	departure := GetTime("2027-01-01 00:00:00")
	speed := 1000.0

	g.GenerateSingleTravel(pointA, pointB, departure, speed)

	travels := g.Travels()

	if len(travels) != 1 {
		t.Errorf("Expected 1 travel, got %d", len(travels))
		return
	}

	expectedArrival := GetTime("2027-01-01 10:00:00")
	expectedTravel := &tables.Travel{
		From:      "1",
		To:        "2",
		Departure: departure,
		Arrival:   expectedArrival,
	}

	actualTravel := travels[0]

	// Compare without ID since ID might be auto-generated
	if actualTravel.From != expectedTravel.From {
		t.Errorf("Expected From=%s, got From=%s", expectedTravel.From, actualTravel.From)
	}
	if actualTravel.To != expectedTravel.To {
		t.Errorf("Expected To=%s, got To=%s", expectedTravel.To, actualTravel.To)
	}
	if !actualTravel.Departure.Equal(expectedTravel.Departure) {
		t.Errorf("Expected Departure=%v, got Departure=%v", expectedTravel.Departure, actualTravel.Departure)
	}
	if !actualTravel.Arrival.Equal(expectedTravel.Arrival) {
		t.Errorf("Expected Arrival=%v, got Arrival=%v", expectedTravel.Arrival, actualTravel.Arrival)
	}
}

func TestGenerateSingleTravelDiagonal(t *testing.T) {
	gf := &GeneratorFactory{}
	g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

	// Two points forming a 3-4-5 right triangle (distance = 5000)
	pointA := tables.Point{ID: "1", X: 0, Y: 0}
	pointB := tables.Point{ID: "2", X: 3000, Y: 4000}

	// Speed: 500 units per hour
	// Distance: sqrt(3000^2 + 4000^2) = 5000 units
	// Expected travel time: 10 hours
	departure := GetTime("2027-01-15 08:30:00")
	speed := 500.0

	g.GenerateSingleTravel(pointA, pointB, departure, speed)

	travels := g.Travels()

	if len(travels) != 1 {
		t.Errorf("Expected 1 travel, got %d", len(travels))
		return
	}

	expectedArrival := GetTime("2027-01-15 18:30:00")
	actualTravel := travels[0]

	if actualTravel.From != "1" {
		t.Errorf("Expected From=1, got From=%s", actualTravel.From)
	}
	if actualTravel.To != "2" {
		t.Errorf("Expected To=2, got To=%s", actualTravel.To)
	}
	if !actualTravel.Departure.Equal(departure) {
		t.Errorf("Expected Departure=%v, got Departure=%v", departure, actualTravel.Departure)
	}
	if !actualTravel.Arrival.Equal(expectedArrival) {
		t.Errorf("Expected Arrival=%v, got Arrival=%v", expectedArrival, actualTravel.Arrival)
	}
}

func TestGenerateSingleTravelMultipleCalls(t *testing.T) {
	gf := &GeneratorFactory{}
	g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

	pointA := tables.Point{ID: "1", X: 0, Y: 0}
	pointB := tables.Point{ID: "2", X: 2000, Y: 0}
	pointC := tables.Point{ID: "3", X: 4000, Y: 0}

	departure1 := GetTime("2027-01-01 00:00:00")
	departure2 := GetTime("2027-01-01 10:00:00")
	speed := 1000.0

	g.GenerateSingleTravel(pointA, pointB, departure1, speed)
	g.GenerateSingleTravel(pointB, pointC, departure2, speed)

	travels := g.Travels()

	if len(travels) != 2 {
		t.Errorf("Expected 2 travels, got %d", len(travels))
		return
	}

	// First travel: A to B
	if travels[0].From != "1" || travels[0].To != "2" {
		t.Errorf("First travel: expected From=1, To=2, got From=%s, To=%s", travels[0].From, travels[0].To)
	}
	expectedArrival1 := GetTime("2027-01-01 02:00:00")
	if !travels[0].Arrival.Equal(expectedArrival1) {
		t.Errorf("First travel: expected Arrival=%v, got Arrival=%v", expectedArrival1, travels[0].Arrival)
	}

	// Second travel: B to C
	if travels[1].From != "2" || travels[1].To != "3" {
		t.Errorf("Second travel: expected From=2, To=3, got From=%s, To=%s", travels[1].From, travels[1].To)
	}
	expectedArrival2 := GetTime("2027-01-01 12:00:00")
	if !travels[1].Arrival.Equal(expectedArrival2) {
		t.Errorf("Second travel: expected Arrival=%v, got Arrival=%v", expectedArrival2, travels[1].Arrival)
	}
}

// Helper function to calculate expected distance
func calculateDistance(p1, p2 tables.Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}
