package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/generator"
	"darbelis.eu/persedimai/internal/migrations"
	"fmt"
	"testing"
	"time"
)

func TestGenerateTravelsToDB(t *testing.T) {
	gf := generator.GeneratorFactory{}

	var idGenerator generator.IdGenerator
	idGenerator = &generator.SimpleIdGenerator{}
	g := gf.CreateGenerator(5, 1000, 0, idGenerator)

	db, err := di.NewDatabase("test")
	if err != nil {
		t.Error(err)
	}

	// Create both points and travels tables
	err = migrations.CreatePointsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	err = migrations.CreateTravelsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	if !ClearTestDatabase(db, "points") {
		fmt.Println("Failed to clear points table")
	}

	if !ClearTestDatabase(db, "travels") {
		fmt.Println("Failed to clear travels table")
	}

	// Generate and insert points first
	pointDao := dao.NewPointDao(db)
	pointDbConsumer := generator.NewPointConsumer(pointDao, 10)
	err = g.GeneratePoints(pointDbConsumer)
	if err != nil {
		t.Error(err)
	}

	err = pointDbConsumer.Flush()
	if err != nil {
		t.Error(err)
	}

	// Retrieve points from database
	points, err := pointDao.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// Generate and insert travels
	travelDao := dao.NewTravelDao(db)
	travelDbConsumer := generator.NewTravelConsumer(travelDao, 10)

	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	speed := 1000.0
	restHours := 2

	err = g.GenerateTravels(points, fromDate, toDate, speed, restHours, travelDbConsumer)
	if err != nil {
		t.Error(err)
	}

	err = travelDbConsumer.Flush()
	if err != nil {
		t.Error(err)
	}

	// Select points from database and assert count
	pointsFromDb, err := pointDao.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	expectedPointsCount := 9
	actualPointsCount := len(pointsFromDb)
	if actualPointsCount != expectedPointsCount {
		t.Errorf("Expected %d points in database, got %d", expectedPointsCount, actualPointsCount)
	}

	// Select travels from database and assert count
	travelsFromDb, err := travelDao.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	expectedTravelsCount := 112
	actualTravelsCount := len(travelsFromDb)
	if actualTravelsCount != expectedTravelsCount {
		t.Errorf("Expected %d travels in database, got %d", expectedTravelsCount, actualTravelsCount)
	}

	t.Logf("Successfully generated and verified: %d points and %d travels in database", actualPointsCount, actualTravelsCount)
}
