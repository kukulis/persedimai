//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/internal/generator"
	"darbelis.eu/persedimai/internal/tables"
	"fmt"
	"time"
)

func DemoRandomFactor() {
	fmt.Println("=== Demonstrating Random Factor in Travel Generation ===")

	// Test with randFactor = 0 (no randomness)
	fmt.Println("1. With randFactor = 0 (deterministic):")
	gf := generator.GeneratorFactory{}
	g1 := gf.CreateGenerator(3, 1000, 0, &generator.SimpleIdGenerator{})

	pointA := tables.Point{ID: "A", X: 0, Y: 0}
	pointB := tables.Point{ID: "B", X: 10000, Y: 0}
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)

	travelConsumer1 := generator.NewTravelArrayConsumer()
	g1.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, travelConsumer1)

	fmt.Printf("  Generated %d travels\n", len(travelConsumer1.Travels))
	for i, t := range travelConsumer1.Travels {
		if i < 4 { // Show first 4 travels
			duration := t.Arrival.Sub(t.Departure)
			fmt.Printf("  Travel %d: %s→%s, Duration: %v\n", i+1, t.From, t.To, duration)
		}
	}

	// Test with randFactor = 0.1 (10% variation)
	fmt.Println("\n2. With randFactor = 0.1 (±10% variation):")
	g2 := gf.CreateGenerator(3, 1000, 0.1, &generator.SimpleIdGenerator{})

	travelConsumer2 := generator.NewTravelArrayConsumer()
	g2.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, travelConsumer2)

	fmt.Printf("  Generated %d travels\n", len(travelConsumer2.Travels))
	for i, t := range travelConsumer2.Travels {
		if i < 4 { // Show first 4 travels
			duration := t.Arrival.Sub(t.Departure)
			fmt.Printf("  Travel %d: %s→%s, Duration: %v\n", i+1, t.From, t.To, duration)
		}
	}

	// Test with randFactor = 0.2 (20% variation)
	fmt.Println("\n3. With randFactor = 0.2 (±20% variation):")
	g3 := gf.CreateGenerator(3, 1000, 0.2, &generator.SimpleIdGenerator{})

	travelConsumer3 := generator.NewTravelArrayConsumer()
	g3.GenerateTravelsForTwoPoints(pointA, pointB, fromDate, toDate, 1000, 2, travelConsumer3)

	fmt.Printf("  Generated %d travels\n", len(travelConsumer3.Travels))
	for i, t := range travelConsumer3.Travels {
		if i < 4 { // Show first 4 travels
			duration := t.Arrival.Sub(t.Departure)
			fmt.Printf("  Travel %d: %s→%s, Duration: %v\n", i+1, t.From, t.To, duration)
		}
	}

	fmt.Println("\nNote: With higher randFactor, you'll see:")
	fmt.Println("  - Varied travel durations (speed variations)")
	fmt.Println("  - Different rest periods between consecutive travels")
	fmt.Println("  - Possibly different total number of travels (due to timing changes)")
}
