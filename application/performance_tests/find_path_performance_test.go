package performance_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/travel_finder"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// TestPerformanceFindPath tests the performance of SimpleTravelSearchStrategy.FindPath
// with a large dataset (~1000 points, ~1 million travels)
func TestPerformanceFindPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Setup database and fill with test data
	t.Log("Setting up database and filling with test data...")
	setupStart := time.Now()

	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	// Fill test database (this takes ~45 seconds)
	err = integration_tests.FillTestDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	setupDuration := time.Since(setupStart)
	t.Logf("Setup completed in %v", setupDuration)

	// Create DAOs and strategy
	pointDao := dao.NewPointDao(db)
	travelDao := dao.NewTravelDao(db)
	strategy := travel_finder.NewSimpleTravelSearchStrategy(travelDao)

	// Get point count for reporting
	pointCount, err := pointDao.Count()
	if err != nil {
		t.Fatal(err)
	}

	travelCount, err := travelDao.Count()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Test data: %d points, %d travels", pointCount, travelCount)

	// Helper function to generate valid coordinates
	// Points are at coordinates: 0, 6000, 12000, 18000, ..., 186000
	// (n=63, squareSize=3000, skip pattern i%2==0 gives indices 0,2,4,...,62)
	getValidCoordinate := func() float64 {
		// Valid indices: 0, 2, 4, 6, ..., 62 (32 values)
		validIndex := rand.Intn(32) * 2 // 0, 2, 4, ..., 62
		return float64(validIndex * 3000)
	}

	// Helper function to get a random point by coordinates
	getRandomPoint := func() (*tables.Point, error) {
		x := getValidCoordinate()
		y := getValidCoordinate()
		return pointDao.FindByCoordinates(x, y)
	}

	// Number of test iterations
	numTests := 100
	t.Logf("Running %d FindPath tests...", numTests)

	// Time period for searches (same as test data)
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 3, 1, 0, 0, 0, 0, time.UTC)

	// Statistics tracking
	var durations []time.Duration
	successCount := 0
	noPathCount := 0
	errorCount := 0

	// Run performance tests
	testStart := time.Now()

	for i := 0; i < numTests; i++ {
		// Get random source and destination points
		source, err := getRandomPoint()
		if err != nil {
			t.Logf("Warning: Failed to get source point: %v", err)
			errorCount++
			continue
		}

		destination, err := getRandomPoint()
		if err != nil {
			t.Logf("Warning: Failed to get destination point: %v", err)
			errorCount++
			continue
		}

		// Avoid same source and destination
		if source.ID == destination.ID {
			i-- // retry this iteration
			continue
		}

		// Convert point IDs to integers for filter
		sourceID, _ := strconv.Atoi(source.ID)
		destID, _ := strconv.Atoi(destination.ID)

		filter := &data.TravelFilter{
			Source:          sourceID,
			Destination:     destID,
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     2,
		}

		// Measure FindPath duration only
		start := time.Now()
		path, err := strategy.FindPath(filter)
		duration := time.Since(start)

		durations = append(durations, duration)

		if err != nil {
			t.Logf("Test %d: Error - %v", i+1, err)
			errorCount++
			continue
		}

		if path == nil {
			noPathCount++
		} else {
			successCount++
			if i < 5 { // Log first 5 successful paths
				t.Logf("Test %d: Found path %s -> %s (via %s), duration: %v, query time: %v",
					i+1, source.ID, destination.ID, path.Travels[0].To, path.TotalDuration, duration)
			}
		}
	}

	totalTestDuration := time.Since(testStart)

	// Calculate statistics
	if len(durations) == 0 {
		t.Fatal("No successful tests to analyze")
	}

	var totalDuration time.Duration
	minDuration := durations[0]
	maxDuration := durations[0]

	for _, d := range durations {
		totalDuration += d
		if d < minDuration {
			minDuration = d
		}
		if d > maxDuration {
			maxDuration = d
		}
	}

	avgDuration := totalDuration / time.Duration(len(durations))

	// Calculate median
	sortedDurations := make([]time.Duration, len(durations))
	copy(sortedDurations, durations)
	// Simple bubble sort for median
	for i := 0; i < len(sortedDurations); i++ {
		for j := i + 1; j < len(sortedDurations); j++ {
			if sortedDurations[i] > sortedDurations[j] {
				sortedDurations[i], sortedDurations[j] = sortedDurations[j], sortedDurations[i]
			}
		}
	}
	medianDuration := sortedDurations[len(sortedDurations)/2]

	// Calculate percentiles
	p95Index := int(float64(len(sortedDurations)) * 0.95)
	p99Index := int(float64(len(sortedDurations)) * 0.99)
	p95Duration := sortedDurations[p95Index]
	p99Duration := sortedDurations[p99Index]

	// Report results
	t.Log("\n========================================")
	t.Log("PERFORMANCE TEST RESULTS")
	t.Log("========================================")
	t.Logf("Dataset: %d points, %d transfers", pointCount, travelCount)
	t.Logf("Tests run: %d", numTests)
	t.Logf("Successful paths found: %d", successCount)
	t.Logf("No path found: %d", noPathCount)
	t.Logf("Errors: %d", errorCount)
	t.Log("----------------------------------------")
	t.Logf("Min query time:    %v", minDuration)
	t.Logf("Max query time:    %v", maxDuration)
	t.Logf("Avg query time:    %v", avgDuration)
	t.Logf("Median query time: %v", medianDuration)
	t.Logf("95th percentile:   %v", p95Duration)
	t.Logf("99th percentile:   %v", p99Duration)
	t.Log("----------------------------------------")
	t.Logf("Total test time:   %v", totalTestDuration)
	t.Logf("Queries per second: %.2f", float64(numTests)/totalTestDuration.Seconds())
	t.Log("========================================")

	// Performance thresholds (adjust based on requirements)
	if avgDuration > 100*time.Millisecond {
		t.Logf("Warning: Average query time (%v) exceeds 100ms threshold", avgDuration)
	}
	if p95Duration > 200*time.Millisecond {
		t.Logf("Warning: 95th percentile (%v) exceeds 200ms threshold", p95Duration)
	}
}

// BenchmarkFindPath is a standard Go benchmark for FindPath
func BenchmarkFindPath(b *testing.B) {
	// Setup database (only once)
	db, err := di.NewDatabase("test")
	if err != nil {
		b.Fatal(err)
	}

	// Check if data exists, if not fill it
	travelDao := dao.NewTravelDao(db)
	count, err := travelDao.Count()
	if err != nil || count == 0 {
		b.Log("Filling test database...")
		err = integration_tests.FillTestDatabase(db)
		if err != nil {
			b.Fatal(err)
		}
	}

	pointDao := dao.NewPointDao(db)
	strategy := travel_finder.NewSimpleTravelSearchStrategy(travelDao)

	// Helper to get valid coordinate
	getValidCoordinate := func() float64 {
		validIndex := rand.Intn(32) * 2
		return float64(validIndex * 3000)
	}

	// Pre-generate test cases
	type testCase struct {
		sourceID int
		destID   int
	}

	var testCases []testCase
	for i := 0; i < 1000; i++ {
		x1 := getValidCoordinate()
		y1 := getValidCoordinate()
		x2 := getValidCoordinate()
		y2 := getValidCoordinate()

		p1, err := pointDao.FindByCoordinates(x1, y1)
		if err != nil {
			continue
		}
		p2, err := pointDao.FindByCoordinates(x2, y2)
		if err != nil {
			continue
		}

		if p1.ID == p2.ID {
			continue
		}

		sourceID, _ := strconv.Atoi(p1.ID)
		destID, _ := strconv.Atoi(p2.ID)

		testCases = append(testCases, testCase{sourceID, destID})

		if len(testCases) >= 100 {
			break
		}
	}

	if len(testCases) == 0 {
		b.Fatal("No test cases generated")
	}

	b.Logf("Generated %d test cases", len(testCases))

	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 3, 1, 0, 0, 0, 0, time.UTC)

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		tc := testCases[i%len(testCases)]

		filter := &data.TravelFilter{
			Source:          tc.sourceID,
			Destination:     tc.destID,
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     2,
		}

		_, err := strategy.FindPath(filter)
		if err != nil {
			// Ignore errors in benchmark
		}
	}
}
