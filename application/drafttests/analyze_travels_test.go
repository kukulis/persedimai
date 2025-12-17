//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/generator"
	"math"
	"testing"
	"time"
)

func TestAnalyzeTravelsGeneration(t *testing.T) {
	t.Log("=== Analyzing GenerateTravels output ===")

	// Setup
	gf := generator.GeneratorFactory{}
	g := gf.CreateGenerator(5, 1000, 0, &generator.SimpleIdGenerator{})

	// Generate points
	pointArrayConsumer := generator.NewPointArrayConsumer()
	err := g.GeneratePoints(pointArrayConsumer)
	if err != nil {
		t.Fatalf("GeneratePoints error: %v", err)
	}

	points := pointArrayConsumer.Points
	t.Logf("Generated %d points in the grid", len(points))

	// Grid structure explanation
	t.Log("\nGrid structure:")
	t.Log("  - Generator with n=5, squareSize=1000")
	t.Log("  - Points generated where i%2==0 and j%2==0")
	t.Log("  - Results in a 3×3 grid (indices: 0, 2, 4)")
	t.Log("\nPoint coordinates:")
	for i, p := range points {
		t.Logf("  P%d: (%.0f, %.0f)", i+1, p.X, p.Y)
	}

	// Count neighbor pairs by type
	horizontalPairs := 0
	verticalPairs := 0
	diagonalPairs := 0

	squareSize := 1000.0
	pointsMap := make(map[string]bool)
	for _, p := range points {
		key := p.BuildLocationKey()
		pointsMap[key] = true
	}

	processed := make(map[string]bool)
	deltas := []float64{-squareSize * 2, 0, squareSize * 2}

	for _, p1 := range points {
		for _, dX := range deltas {
			for _, dY := range deltas {
				if dX == 0 && dY == 0 {
					continue
				}

				neighborX := p1.X + dX
				neighborY := p1.Y + dY
				neighborKey := (&struct {
					X, Y float64
				}{neighborX, neighborY})

				// Simulate BuildLocationKey
				neighborKeyStr := ""
				for _, p2 := range points {
					if p2.X == neighborX && p2.Y == neighborY {
						key1 := p1.ID + "_" + p2.ID
						key2 := p2.ID + "_" + p1.ID

						if !processed[key1] && !processed[key2] {
							processed[key1] = true

							if dX == 0 {
								verticalPairs++
							} else if dY == 0 {
								horizontalPairs++
							} else {
								diagonalPairs++
							}
						}
						break
					}
				}
				_ = neighborKeyStr
				_ = neighborKey
			}
		}
	}

	totalPairs := horizontalPairs + verticalPairs + diagonalPairs
	t.Log("\nNeighbor pairs analysis:")
	t.Logf("  Horizontal pairs: %d", horizontalPairs)
	t.Logf("  Vertical pairs: %d", verticalPairs)
	t.Logf("  Diagonal pairs: %d", diagonalPairs)
	t.Logf("  Total unique pairs: %d", totalPairs)

	// Calculate travels per pair type
	speed := 1000.0
	restHours := 2
	totalHours := 24.0

	// Horizontal/Vertical (distance = 2000 units)
	hvDistance := 2000.0
	hvTravelTime := hvDistance / speed               // 2.0 hours
	hvCycleTime := hvTravelTime + float64(restHours) // 4.0 hours
	hvTravelsPerPair := int(totalHours / hvCycleTime)

	t.Log("\nHorizontal/Vertical travel calculation:")
	t.Logf("  Distance: %.0f units", hvDistance)
	t.Logf("  Travel time: %.1f hours", hvTravelTime)
	t.Logf("  Rest time: %d hours", restHours)
	t.Logf("  Cycle time: %.1f hours", hvCycleTime)
	t.Logf("  Transfers per pair in 24h: %d", hvTravelsPerPair)

	// Diagonal (distance = sqrt(2000² + 2000²) ≈ 2828.43 units)
	diagDistance := math.Sqrt(2000*2000 + 2000*2000)
	diagTravelTime := diagDistance / speed               // ~2.828 hours
	diagCycleTime := diagTravelTime + float64(restHours) // ~4.828 hours

	// Precise diagonal count
	currentTime := 0.0
	diagCount := 0
	t.Log("\nDiagonal travel simulation:")
	for {
		arrivalTime := currentTime + diagTravelTime
		if arrivalTime > totalHours {
			t.Logf("  Travel %d would arrive at %.2fh (exceeds 24h) - STOP", diagCount+1, arrivalTime)
			break
		}
		diagCount++
		t.Logf("  Travel %d: depart %.2fh → arrive %.2fh", diagCount, currentTime, arrivalTime)

		nextDeparture := arrivalTime + float64(restHours)
		if nextDeparture >= totalHours {
			t.Logf("  Next departure would be at %.2fh (≥24h) - STOP", nextDeparture)
			break
		}
		currentTime = nextDeparture
	}

	t.Log("\nDiagonal travel calculation:")
	t.Logf("  Distance: %.2f units", diagDistance)
	t.Logf("  Travel time: %.3f hours", diagTravelTime)
	t.Logf("  Rest time: %d hours", restHours)
	t.Logf("  Cycle time: %.3f hours", diagCycleTime)
	t.Logf("  Transfers per pair in 24h: %d", diagCount)

	// Total calculation
	hvTotal := (horizontalPairs + verticalPairs) * hvTravelsPerPair
	diagTotal := diagonalPairs * diagCount
	expectedTotal := hvTotal + diagTotal

	t.Log("\n=== TOTAL TRAVELS CALCULATION ===")
	t.Logf("  Horizontal/Vertical: %d pairs × %d travels = %d", horizontalPairs+verticalPairs, hvTravelsPerPair, hvTotal)
	t.Logf("  Diagonal: %d pairs × %d travels = %d", diagonalPairs, diagCount, diagTotal)
	t.Logf("  EXPECTED TOTAL: %d travels", expectedTotal)

	// Now actually generate and verify
	travelConsumer := generator.NewTravelArrayConsumer()
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)

	err = g.GenerateTravels(points, fromDate, toDate, speed, restHours, travelConsumer)
	if err != nil {
		t.Fatalf("GenerateTravels error: %v", err)
	}

	actualTotal := len(travelConsumer.Travels)
	t.Logf("  ACTUAL TOTAL: %d travels", actualTotal)

	if actualTotal == expectedTotal {
		t.Logf("\n✓ SUCCESS: Actual matches expected!")
	} else {
		t.Errorf("\n✗ MISMATCH: Expected %d but got %d", expectedTotal, actualTotal)
	}
}
