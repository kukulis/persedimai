package integration_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/migrations"
	"darbelis.eu/persedimai/travel_finder"
	"testing"
	"time"
)

func TestSimpleTravelSearchStrategy_Integration(t *testing.T) {
	// Setup database
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	// Create tables
	err = migrations.CreatePointsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	err = migrations.CreateTravelsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	// Clear existing data
	ClearTestDatabase(db, "points")
	ClearTestDatabase(db, "travels")

	// Generate test data
	gf := generator.GeneratorFactory{}
	idGenerator := &generator.SimpleIdGenerator{}
	g := gf.CreateGenerator(5, 1000, 0, idGenerator)

	// Generate points
	pointDao := dao.NewPointDao(db)
	pointDbConsumer := generator.NewPointConsumer(pointDao, 100)
	err = g.GeneratePoints(pointDbConsumer)
	if err != nil {
		t.Fatal(err)
	}
	err = pointDbConsumer.Flush()
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve points
	points, err := pointDao.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// Generate travels for 24 hours
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	speed := 1000.0
	restHours := 2

	travelDao := dao.NewTravelDao(db)
	travelDbConsumer := generator.NewTravelConsumer(travelDao, 100)

	err = g.GenerateTravels(points, fromDate, toDate, speed, restHours, travelDbConsumer)
	if err != nil {
		t.Fatal(err)
	}
	err = travelDbConsumer.Flush()
	if err != nil {
		t.Fatal(err)
	}

	// Create strategy
	strategy := travel_finder.NewSimpleTravelSearchStrategy(travelDao)

	t.Run("FindPath_DirectConnection", func(t *testing.T) {
		// Test direct connection: point 1 -> point 2
		filter := &data.TravelFilter{
			Source:          "1",
			Destination:     "2",
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     1,
		}

		path, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if path == nil {
			t.Fatal("FindPath returned nil path")
		}

		if len(path.Travels) != 1 {
			t.Errorf("Expected 1 transfer, got %d", len(path.Travels))
		}

		if path.TransferCount != 0 {
			t.Errorf("Expected TransferCount=0 (no intermediate stops), got %d", path.TransferCount)
		}

		expectedDuration := 2 * time.Hour
		if path.TotalDuration != expectedDuration {
			t.Errorf("Expected duration %v, got %v", expectedDuration, path.TotalDuration)
		}

		t.Logf("Direct path: %s -> %s, duration: %v",
			path.Travels[0].From, path.Travels[0].To, path.TotalDuration)
	})

	t.Run("FindPath_TwoTransfers", func(t *testing.T) {
		// Test 2-transfer path: point 1 -> intermediate -> point 9
		filter := &data.TravelFilter{
			Source:          "1",
			Destination:     "9",
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     2,
		}

		path, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if path == nil {
			t.Fatal("FindPath returned nil path")
		}

		if len(path.Travels) != 2 {
			t.Errorf("Expected 2 transfers, got %d", len(path.Travels))
		}

		if path.TransferCount != 1 {
			t.Errorf("Expected TransferCount=1 (one intermediate stop), got %d", path.TransferCount)
		}

		// Verify path continuity
		if path.Travels[0].To != path.Travels[1].From {
			t.Errorf("Path not continuous: %s != %s", path.Travels[0].To, path.Travels[1].From)
		}

		t.Logf("Two-transfer path: %s -> %s -> %s, duration: %v",
			path.Travels[0].From, path.Travels[0].To, path.Travels[1].To, path.TotalDuration)
	})

	t.Run("FindPath_NoPath", func(t *testing.T) {
		// Test non-existent route
		filter := &data.TravelFilter{
			Source:          "999",
			Destination:     "888",
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     1,
		}

		path, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if path != nil {
			t.Errorf("Expected nil path for non-existent route, got path with %d transfers", len(path.Travels))
		}

		t.Log("Correctly returned nil for non-existent route")
	})

	t.Run("FindPath_InvalidTravelCount", func(t *testing.T) {
		filter := &data.TravelFilter{
			Source:          "1",
			Destination:     "2",
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     5,
		}

		path, err := strategy.FindPath(filter)
		if err == nil {
			t.Fatal("Expected error for TravelCount=5, got nil")
		}

		if path != nil {
			t.Error("Expected nil path when error occurs")
		}

		if err.Error() != "unimplemented: TravelCount > 2 not supported" {
			t.Errorf("Unexpected error message: %v", err)
		}

		t.Logf("Correctly returned error: %v", err)
	})

	t.Run("Strategy_GetName", func(t *testing.T) {
		if strategy.GetName() != "Simple" {
			t.Errorf("Expected name 'Simple', got '%s'", strategy.GetName())
		}
	})
}
