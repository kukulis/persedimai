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
		filter := data.NewTravelFilter("1", "2", fromDate, toDate, 1)

		paths, err := strategy.FindPaths(filter)
		if err != nil {
			t.Fatalf("FindPaths returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPaths didnt find any path")
		}

		for _, path := range paths {

			if len(path.Transfers) != 1 {
				t.Errorf("Expected 1 transfer, got %d", len(path.Transfers))
			}

			if path.TransferCount != 1 {
				t.Errorf("Expected TransferCount=1 (no intermediate stops), got %d", path.TransferCount)
			}

			expectedDuration := 2 * time.Hour
			if path.TotalDuration != expectedDuration {
				t.Errorf("Expected duration %v, got %v", expectedDuration, path.TotalDuration)
			}

			t.Logf("Direct path: %s -> %s, duration: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.TotalDuration)
		}
	})

	t.Run("FindPath_TwoTransfers", func(t *testing.T) {
		// Test 2-transfer path: point 1 -> intermediate -> point 9
		filter := data.NewTravelFilter("1", "9", fromDate, toDate, 2)

		paths, err := strategy.FindPaths(filter)
		if err != nil {
			t.Fatalf("FindPaths returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPaths didnt find any path")
		}

		for _, path := range paths {
			if len(path.Transfers) != 2 {
				t.Errorf("Expected 2 transfers, got %d", len(path.Transfers))
			}

			if path.TransferCount != 2 {
				t.Errorf("Expected TransferCount=2 (one intermediate stop), got %d", path.TransferCount)
			}

			// Verify path continuity
			if path.Transfers[0].To != path.Transfers[1].From {
				t.Errorf("Path not continuous: %s != %s", path.Transfers[0].To, path.Transfers[1].From)
			}

			t.Logf("Two-transfer path: %s -> %s -> %s, duration: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.Transfers[1].To, path.TotalDuration)
		}
	})

	t.Run("FindPath_NoPath", func(t *testing.T) {
		// Test non-existent route
		filter := data.NewTravelFilter("999", "888", fromDate, toDate, 1)

		paths, err := strategy.FindPaths(filter)
		if err != nil {
			t.Fatalf("FindPaths returned error: %v", err)
		}

		if len(paths) != 0 {
			t.Errorf("Expected zero paths for non-existent route, got path with %d transfers", len(paths))
		}

		t.Log("Correctly returned nil for non-existent route")
	})

	t.Run("FindPath_InvalidTravelCount", func(t *testing.T) {
		filter := data.NewTravelFilter("1", "2", fromDate, toDate, 5)

		path, err := strategy.FindPaths(filter)
		if err == nil {
			t.Fatal("Expected error for TravelCount=5, got nil")
		}

		if path != nil {
			t.Error("Expected nil path when error occurs")
		}

		if err.Error() != "unimplemented: TravelCount > 3 not supported" {
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
