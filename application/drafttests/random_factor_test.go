//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/tables"
	"testing"
	"time"
)

func TestRandomFactorZero(t *testing.T) {
	// With randFactor = 0, results should be deterministic
	gf := generator.GeneratorFactory{}
	g := gf.CreateGenerator(3, 1000, 0, &generator.SimpleIdGenerator{})

	pointA := tables.Point{ID: "A", X: 0, Y: 0}
	pointB := tables.Point{ID: "B", X: 10000, Y: 0}
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)

	// Generate twice and compare
	consumer1 := generator.NewTravelArrayConsumer()
	g.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, consumer1)

	consumer2 := generator.NewTravelArrayConsumer()
	g.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, consumer2)

	if len(consumer1.Travels) != len(consumer2.Travels) {
		t.Errorf("With randFactor=0, expected same number of travels, got %d and %d",
			len(consumer1.Travels), len(consumer2.Travels))
	}

	// Check that travel durations are identical
	for i := 0; i < len(consumer1.Travels); i++ {
		duration1 := consumer1.Travels[i].Arrival.Sub(consumer1.Travels[i].Departure)
		duration2 := consumer2.Travels[i].Arrival.Sub(consumer2.Travels[i].Departure)

		if duration1 != duration2 {
			t.Errorf("Travel %d: with randFactor=0, expected same duration, got %v and %v",
				i, duration1, duration2)
		}
	}

	t.Logf("✓ With randFactor=0: Generated %d identical travels both times", len(consumer1.Travels))
}

func TestRandomFactorNonZero(t *testing.T) {
	// With randFactor > 0, results should vary
	gf := generator.GeneratorFactory{}
	g := gf.CreateGenerator(3, 1000, 0.15, &generator.SimpleIdGenerator{})

	pointA := tables.Point{ID: "A", X: 0, Y: 0}
	pointB := tables.Point{ID: "B", X: 10000, Y: 0}
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)

	// Generate twice
	consumer1 := generator.NewTravelArrayConsumer()
	g.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, consumer1)

	consumer2 := generator.NewTravelArrayConsumer()
	g.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, consumer2)

	// Count how many travels have different durations
	differentCount := 0
	minLen := len(consumer1.Travels)
	if len(consumer2.Travels) < minLen {
		minLen = len(consumer2.Travels)
	}

	for i := 0; i < minLen; i++ {
		duration1 := consumer1.Travels[i].Arrival.Sub(consumer1.Travels[i].Departure)
		duration2 := consumer2.Travels[i].Arrival.Sub(consumer2.Travels[i].Departure)

		if duration1 != duration2 {
			differentCount++
		}
	}

	// With randFactor=0.15, we expect most travels to have different durations
	if differentCount == 0 {
		t.Error("With randFactor=0.15, expected some variation in travel durations, but all were identical")
	}

	t.Logf("✓ With randFactor=0.15: %d/%d travels had different durations", differentCount, minLen)
	t.Logf("  First run: %d travels", len(consumer1.Travels))
	t.Logf("  Second run: %d travels", len(consumer2.Travels))
}

func TestRandomFactorRange(t *testing.T) {
	// Verify that random factor produces values within expected range
	gf := generator.GeneratorFactory{}
	randFactor := 0.2 // ±20%
	g := gf.CreateGenerator(3, 1000, randFactor, &generator.SimpleIdGenerator{})

	pointA := tables.Point{ID: "A", X: 0, Y: 0}
	pointB := tables.Point{ID: "B", X: 10000, Y: 0}
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 3, 0, 0, 0, 0, time.UTC) // 2 days for more samples

	consumer := generator.NewTravelArrayConsumer()
	g.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, consumer)

	// Expected duration without randomness: 10000/1000 = 10 hours
	expectedDuration := 10 * time.Hour
	minExpected := time.Duration(float64(expectedDuration) * (1 - randFactor))
	maxExpected := time.Duration(float64(expectedDuration) * (1 + randFactor))

	outOfRangeCount := 0
	for _, travel := range consumer.Travels {
		duration := travel.Arrival.Sub(travel.Departure)

		if duration < minExpected || duration > maxExpected {
			outOfRangeCount++
			t.Logf("  Travel %s→%s has duration %v (outside expected range %v-%v)",
				travel.From, travel.To, duration, minExpected, maxExpected)
		}
	}

	if outOfRangeCount > 0 {
		t.Errorf("Found %d/%d travels outside expected duration range", outOfRangeCount, len(consumer.Travels))
	} else {
		t.Logf("✓ All %d travels within expected range: %v to %v", len(consumer.Travels), minExpected, maxExpected)
	}
}
