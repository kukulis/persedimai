package integration_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/migrations"
	"testing"
	"time"
)

func TestFindPathSimple1(t *testing.T) {
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

	// Test FindPathSimple1: direct connection from point 1 to point 2
	filter := &data.TravelFilter{
		Source:          "1",
		Destination:     "2",
		ArrivalTimeFrom: fromDate,
		ArrivalTimeTo:   toDate,
		TravelCount:     1,
	}

	sequences, err := travelDao.FindPathSimple1(filter)
	if err != nil {
		t.Fatalf("FindPathSimple1 returned error: %v", err)
	}

	if sequences == nil || len(sequences) == 0 {
		t.Fatal("FindPathSimple1 returned no sequences")
	}

	// Get the first (best) sequence
	sequence := sequences[0]

	if sequence.Count() != 1 {
		t.Errorf("Expected 1 transfer in sequence, got %d", sequence.Count())
	}

	transfer := sequence.First()

	// Verify the transfer
	if transfer.From != "1" {
		t.Errorf("Expected From='1', got From='%s'", transfer.From)
	}
	if transfer.To != "2" {
		t.Errorf("Expected To='2', got To='%s'", transfer.To)
	}

	if transfer.Arrival.Before(filter.ArrivalTimeFrom) {
		t.Errorf("Arrival %v is before ArrivalTimeFrom %v", transfer.Arrival, filter.ArrivalTimeFrom)
	}
	if transfer.Arrival.After(filter.ArrivalTimeTo) {
		t.Errorf("Arrival %v is after ArrivalTimeTo %v", transfer.Arrival, filter.ArrivalTimeTo)
	}

	// Expected: Point 1 (0,0) to Point 2 (2000,0), distance 2000, speed 1000 = 2 hours
	expectedDuration := 2 * time.Hour
	actualDuration := sequence.TotalDuration()
	if actualDuration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, actualDuration)
	}

	// Verify sequence is valid
	if !sequence.IsValid() {
		t.Error("Sequence is not valid")
	}

	t.Logf("Found direct path: %s -> %s, departure: %v, arrival: %v, duration: %v",
		transfer.From, transfer.To, transfer.Departure, transfer.Arrival, actualDuration)
}

func TestFindPathSimple1_NoPath(t *testing.T) {
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
	ClearTestDatabase(db, "travels")

	travelDao := dao.NewTravelDao(db)

	// Test with non-existent route
	filter := &data.TravelFilter{
		Source:          "999",
		Destination:     "888",
		ArrivalTimeFrom: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		ArrivalTimeTo:   time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC),
		TravelCount:     1,
	}

	sequences, err := travelDao.FindPathSimple1(filter)
	if err != nil {
		t.Fatalf("FindPathSimple1 returned error: %v", err)
	}

	if len(sequences) != 0 {
		t.Errorf("Expected no sequences for non-existent route, got %d", len(sequences))
	}

	t.Log("Correctly returned empty result for non-existent route")
}
