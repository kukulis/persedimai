package generator

import (
	"darbelis.eu/persedimai/tables"
	"testing"
)

func TestGenerateSingleTravel(t *testing.T) {
	t.Run("BasicHorizontalTravel", func(t *testing.T) {
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

		actualTravel := g.GenerateSingleTravel(pointA, pointB, departure, speed)

		expectedArrival := GetTime("2027-01-01 10:00:00")
		expectedTravel := &tables.Travel{
			From:      "1",
			To:        "2",
			Departure: departure,
			Arrival:   expectedArrival,
		}

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
	})

	t.Run("DiagonalTravel", func(t *testing.T) {
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

		actualTravel := g.GenerateSingleTravel(pointA, pointB, departure, speed)

		expectedArrival := GetTime("2027-01-15 18:30:00")

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
	})

	t.Run("MultipleCalls", func(t *testing.T) {
		gf := &GeneratorFactory{}
		g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

		pointA := tables.Point{ID: "1", X: 0, Y: 0}
		pointB := tables.Point{ID: "2", X: 2000, Y: 0}
		pointC := tables.Point{ID: "3", X: 4000, Y: 0}

		departure1 := GetTime("2027-01-01 00:00:00")
		departure2 := GetTime("2027-01-01 10:00:00")
		speed := 1000.0

		travel1 := g.GenerateSingleTravel(pointA, pointB, departure1, speed)
		travel2 := g.GenerateSingleTravel(pointB, pointC, departure2, speed)

		// First travel: A to B
		if travel1.From != "1" || travel1.To != "2" {
			t.Errorf("First travel: expected From=1, To=2, got From=%s, To=%s", travel1.From, travel1.To)
		}
		expectedArrival1 := GetTime("2027-01-01 02:00:00")
		if !travel1.Arrival.Equal(expectedArrival1) {
			t.Errorf("First travel: expected Arrival=%v, got Arrival=%v", expectedArrival1, travel1.Arrival)
		}

		// Second travel: B to C
		if travel2.From != "2" || travel2.To != "3" {
			t.Errorf("Second travel: expected From=2, To=3, got From=%s, To=%s", travel2.From, travel2.To)
		}
		expectedArrival2 := GetTime("2027-01-01 12:00:00")
		if !travel2.Arrival.Equal(expectedArrival2) {
			t.Errorf("Second travel: expected Arrival=%v, got Arrival=%v", expectedArrival2, travel2.Arrival)
		}
	})
}
