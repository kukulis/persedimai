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

func TestFindPathSimple2(t *testing.T) {
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

	// Generate test data - 3x3 grid (9 points)
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

	// Test FindPathSimple2: 2-transfer path from point 1 to point 9
	// Point 1 (0,0) -> Point 5 (2000,2000) -> Point 9 (4000,4000)
	// This is a diagonal path through the center
	filter := &data.TravelFilter{
		Source:          "1",
		Destination:     "9",
		ArrivalTimeFrom: fromDate,
		ArrivalTimeTo:   toDate,
		TravelCount:     2,
	}

	sequences, err := travelDao.FindPathSimple2(filter)
	if err != nil {
		t.Fatalf("FindPathSimple2 returned error: %v", err)
	}

	if sequences == nil || len(sequences) == 0 {
		t.Fatal("FindPathSimple2 returned no sequences")
	}

	t.Logf("Found %d possible paths from point 1 to point 9", len(sequences))

	// Get the first (best) sequence
	sequence := sequences[0]

	if sequence.Count() != 2 {
		t.Errorf("Expected 2 transfers in sequence, got %d", sequence.Count())
	}

	// Verify sequence validity
	if !sequence.IsValid() {
		t.Error("Sequence is not valid")
	}

	transfer1 := sequence.First()
	transfer2 := sequence.Last()

	// Verify the path starts from point 1
	if transfer1.From != "1" {
		t.Errorf("Expected first transfer From='1', got From='%s'", transfer1.From)
	}

	// Verify the path ends at point 9
	if transfer2.To != "9" {
		t.Errorf("Expected last transfer To='9', got To='%s'", transfer2.To)
	}

	// Verify connection: first transfer's destination = second transfer's source
	if transfer1.To != transfer2.From {
		t.Errorf("Transfers not connected: transfer1.To='%s' != transfer2.From='%s'", transfer1.To, transfer2.From)
	}

	// Verify second transfer departs after first arrives
	if transfer2.Departure.Before(transfer1.Arrival) {
		t.Errorf("Second transfer departs before first arrives: %v < %v", transfer2.Departure, transfer1.Arrival)
	}

	// Verify arrival is within time window
	if transfer2.Arrival.Before(filter.ArrivalTimeFrom) {
		t.Errorf("Arrival %v is before ArrivalTimeFrom %v", transfer2.Arrival, filter.ArrivalTimeFrom)
	}
	if transfer2.Arrival.After(filter.ArrivalTimeTo) {
		t.Errorf("Arrival %v is after ArrivalTimeTo %v", transfer2.Arrival, filter.ArrivalTimeTo)
	}

	// Calculate connection time
	connectionTime := sequence.ConnectionTime(0)
	totalConnectionTime := sequence.TotalConnectionTime()

	t.Logf("Found path: %s -> %s -> %s", transfer1.From, transfer1.To, transfer2.To)
	t.Logf("  Transfer 1: depart %v, arrive %v", transfer1.Departure, transfer1.Arrival)
	t.Logf("  Connection time: %v", connectionTime)
	t.Logf("  Transfer 2: depart %v, arrive %v", transfer2.Departure, transfer2.Arrival)
	t.Logf("  Total duration: %v", sequence.TotalDuration())
	t.Logf("  Total connection time: %v", totalConnectionTime)
}

func TestFindPathSimple2_NoPath(t *testing.T) {
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
		TravelCount:     2,
	}

	sequences, err := travelDao.FindPathSimple2(filter)
	if err != nil {
		t.Fatalf("FindPathSimple2 returned error: %v", err)
	}

	if len(sequences) != 0 {
		t.Errorf("Expected no sequences for non-existent route, got %d", len(sequences))
	}

	t.Log("Correctly returned empty result for non-existent route")
}
