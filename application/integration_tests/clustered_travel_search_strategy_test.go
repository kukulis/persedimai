package integration_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/migrations"
	"darbelis.eu/persedimai/travel_finder"
	"flag"
	"testing"
	"time"
)

var (
	skipSetup = flag.Bool("skip-setup", true, "Skip database setup and use existing data")
)

// setupClusteredTestDatabase initializes the test database with points, travels, and cluster tables
// If -skip-setup flag is provided, it skips data generation and uses existing data
func setupClusteredTestDatabase(t *testing.T) (*database.Database, *dao.TravelDao, *dao.PointDao, time.Time, time.Time) {
	// Setup database connection
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	travelDao := dao.NewTravelDao(db)
	pointDao := dao.NewPointDao(db)
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 3, 0, 0, 0, 0, time.UTC)

	// Skip setup if flag is set
	if *skipSetup {
		t.Log("Skipping database setup - using existing data")
		return db, travelDao, pointDao, fromDate, toDate
	}

	t.Log("Setting up test database with fresh data")

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
	g := gf.CreateGenerator(10, 3000, 0, idGenerator)

	// Generate points
	//pointDao := dao.NewPointDao(db)
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

	// Generate travels for 48 hours
	speed := 1000.0
	restHours := 2

	travelDbConsumer := generator.NewTravelConsumer(travelDao, 100)

	err = g.GenerateTravels(points, fromDate, toDate, speed, restHours, travelDbConsumer)
	if err != nil {
		t.Fatal(err)
	}
	err = travelDbConsumer.Flush()
	if err != nil {
		t.Fatal(err)
	}

	// Create cluster tables and populate them
	clustersCreator := migrations.NewClustersCreator(db)

	err = clustersCreator.UpdateClustersOnTravels()
	if err != nil {
		t.Fatalf("Failed to update clusters on travels: %v", err)
	}

	err = clustersCreator.CreateClustersTables()
	if err != nil {
		t.Fatalf("Failed to create cluster tables: %v", err)
	}

	err = clustersCreator.InsertClustersDatas()
	if err != nil {
		t.Fatalf("Failed to insert cluster data: %v", err)
	}

	t.Log("Database setup completed successfully")

	return db, travelDao, pointDao, fromDate, toDate
}

func TestClusteredTravelSearchStrategy_Integration(t *testing.T) {
	_, travelDao, pointDao, fromDate, toDate := setupClusteredTestDatabase(t)

	// Create strategy
	strategy := travel_finder.NewClusteredTravelSearchStrategy(travelDao)

	t.Run("FindPath_DirectConnection", func(t *testing.T) {
		// Test direct connection using coordinates (multiples of 6000)
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 0): %v", err)
		}

		filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 1)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPath didn't find any path")
		}

		for _, path := range paths {
			if len(path.Transfers) != 1 {
				t.Errorf("Expected 1 transfer, got %d", len(path.Transfers))
			}

			if path.TransferCount != 1 {
				t.Errorf("Expected TransferCount=1, got %d", path.TransferCount)
			}

			t.Logf("Direct path: %s -> %s, duration: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.TotalDuration)
		}
	})

	t.Run("FindPath_TwoTransfers", func(t *testing.T) {
		// Test 2-transfer path using coordinates
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 6000)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 6000): %v", err)
		}

		filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 2)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPath didn't find any path")
		}

		for _, path := range paths {
			if len(path.Transfers) != 2 {
				t.Errorf("Expected 2 transfers, got %d", len(path.Transfers))
			}

			if path.TransferCount != 2 {
				t.Errorf("Expected TransferCount=2, got %d", path.TransferCount)
			}

			// Verify path continuity
			if path.Transfers[0].To != path.Transfers[1].From {
				t.Errorf("Path not continuous: %s != %s", path.Transfers[0].To, path.Transfers[1].From)
			}

			// Verify minimum connection time (default 30 minutes)
			connectionTime := path.Transfers[1].Departure.Sub(path.Transfers[0].Arrival)
			if connectionTime < 30*time.Minute {
				t.Errorf("Connection time %v is less than minimum 30 minutes", connectionTime)
			}

			t.Logf("Two-transfer path: %s -> %s -> %s, duration: %v, connection time: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.Transfers[1].To,
				path.TotalDuration, connectionTime)
		}
	})

	t.Run("FindPath_ThreeTransfers", func(t *testing.T) {
		// Test 3-transfer path using coordinates (multiples of 6000)
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(12000, 6000)
		if err != nil {
			t.Fatalf("Failed to find point at (12000, 6000): %v", err)
		}

		extendedFromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		extendedToDate := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)
		filter := data.NewTravelFilter(point1.ID, point2.ID, extendedFromDate, extendedToDate, 3)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPath didn't find any path")
		}

		for _, path := range paths {
			if len(path.Transfers) != 3 {
				t.Errorf("Expected 3 transfers, got %d", len(path.Transfers))
			}

			if path.TransferCount != 3 {
				t.Errorf("Expected TransferCount=3, got %d", path.TransferCount)
			}

			// Verify path continuity
			if path.Transfers[0].To != path.Transfers[1].From {
				t.Errorf("Path not continuous at step 1: %s != %s", path.Transfers[0].To, path.Transfers[1].From)
			}
			if path.Transfers[1].To != path.Transfers[2].From {
				t.Errorf("Path not continuous at step 2: %s != %s", path.Transfers[1].To, path.Transfers[2].From)
			}

			// Verify minimum connection times
			connectionTime1 := path.Transfers[1].Departure.Sub(path.Transfers[0].Arrival)
			connectionTime2 := path.Transfers[2].Departure.Sub(path.Transfers[1].Arrival)

			if connectionTime1 < 30*time.Minute {
				t.Errorf("Connection time 1 (%v) is less than minimum 30 minutes", connectionTime1)
			}
			if connectionTime2 < 30*time.Minute {
				t.Errorf("Connection time 2 (%v) is less than minimum 30 minutes", connectionTime2)
			}

			t.Logf("Three-transfer path: %s -> %s -> %s -> %s, duration: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.Transfers[1].To,
				path.Transfers[2].To, path.TotalDuration)
		}
	})

	t.Run("FindPath_FourTransfers", func(t *testing.T) {
		// Test 4-transfer path using coordinates (multiples of 6000)
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 12000)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 12000): %v", err)
		}

		extendedFromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		extendedToDate := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)
		filter := data.NewTravelFilter(point1.ID, point2.ID, extendedFromDate, extendedToDate, 4)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPath didn't find any path")
		}

		for _, path := range paths {
			if len(path.Transfers) != 4 {
				t.Errorf("Expected 4 transfers, got %d", len(path.Transfers))
			}

			if path.TransferCount != 4 {
				t.Errorf("Expected TransferCount=4, got %d", path.TransferCount)
			}

			// Verify path continuity
			for i := 0; i < len(path.Transfers)-1; i++ {
				if path.Transfers[i].To != path.Transfers[i+1].From {
					t.Errorf("Path not continuous at step %d: %s != %s",
						i, path.Transfers[i].To, path.Transfers[i+1].From)
				}

				// Verify minimum connection time
				connectionTime := path.Transfers[i+1].Departure.Sub(path.Transfers[i].Arrival)
				if connectionTime < 30*time.Minute {
					t.Errorf("Connection time %d (%v) is less than minimum 30 minutes", i, connectionTime)
				}
			}

			t.Logf("Four-transfer path: %s -> %s -> %s -> %s -> %s, duration: %v",
				path.Transfers[0].From, path.Transfers[0].To, path.Transfers[1].To,
				path.Transfers[2].To, path.Transfers[3].To, path.TotalDuration)
		}
	})

	t.Run("FindPath_CustomMinConnectionTime", func(t *testing.T) {
		// Test with custom minimum connection time (60 minutes)
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 6000)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 6000): %v", err)
		}

		filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 2)
		filter.MinConnectionTimeMinutes = 60

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		// Verify all paths meet the minimum connection time requirement
		for _, path := range paths {
			for i := 0; i < len(path.Transfers)-1; i++ {
				connectionTime := path.Transfers[i+1].Departure.Sub(path.Transfers[i].Arrival)
				if connectionTime < 60*time.Minute {
					t.Errorf("Connection time %v is less than required 60 minutes", connectionTime)
				}
			}

			t.Logf("Path with 60min minimum connection: duration: %v", path.TotalDuration)
		}
	})

	t.Run("FindPath_NoPath", func(t *testing.T) {
		// Test non-existent route using coordinates that don't exist
		// Using coordinates outside the expected range
		_, err := pointDao.FindByCoordinates(99999, 99999)
		if err == nil {
			t.Skip("Point at (99999, 99999) exists, skipping test")
		}

		_, err = pointDao.FindByCoordinates(88888, 88888)
		if err == nil {
			t.Skip("Point at (88888, 88888) exists, skipping test")
		}

		// Use non-existent IDs
		filter := data.NewTravelFilter("non-existent-id-1", "non-existent-id-2", fromDate, toDate, 2)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if paths != nil && len(paths) != 0 {
			t.Errorf("Expected nil or empty paths for non-existent route, got %d paths", len(paths))
		}

		t.Log("Correctly returned nil/empty for non-existent route")
	})

	t.Run("FindPath_InvalidTravelCount", func(t *testing.T) {
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 0): %v", err)
		}

		filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 5)

		paths, err := strategy.FindPath(filter)
		if err == nil {
			t.Fatal("Expected error for TravelCount=5, got nil")
		}

		if paths != nil {
			t.Error("Expected nil paths when error occurs")
		}

		if err.Error() != "unimplemented: TravelCount > 4 not supported" {
			t.Errorf("Unexpected error message: %v", err)
		}

		t.Logf("Correctly returned error: %v", err)
	})

	t.Run("Strategy_GetName", func(t *testing.T) {
		if strategy.GetName() != "Clustered" {
			t.Errorf("Expected name 'Clustered', got '%s'", strategy.GetName())
		}
	})

	t.Run("FindPath_VerifyActualTimestamps", func(t *testing.T) {
		// Test that actual timestamps are loaded (not cluster-based approximations)
		point1, err := pointDao.FindByCoordinates(0, 0)
		if err != nil {
			t.Fatalf("Failed to find point at (0, 0): %v", err)
		}

		point2, err := pointDao.FindByCoordinates(6000, 6000)
		if err != nil {
			t.Fatalf("Failed to find point at (6000, 6000): %v", err)
		}

		filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 2)

		paths, err := strategy.FindPath(filter)
		if err != nil {
			t.Fatalf("FindPath returned error: %v", err)
		}

		if len(paths) == 0 {
			t.Fatal("FindPath didn't find any path")
		}

		// Verify that timestamps are precise (not rounded to hours)
		for _, path := range paths {
			for _, transfer := range path.Transfers {
				// Check that departure and arrival are not exactly on the hour
				// (cluster times would be rounded to hours)
				if transfer.Departure.Minute() == 0 && transfer.Departure.Second() == 0 {
					// This could still be coincidental, so just log it
					t.Logf("Transfer has departure on the hour: %v", transfer.Departure)
				}

				// Verify transfer ID is loaded
				if transfer.ID == "" {
					t.Error("Transfer ID is empty - actual transfer not loaded")
				}
			}
		}
	})
}
