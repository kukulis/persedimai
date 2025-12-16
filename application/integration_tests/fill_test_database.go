package integration_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/migrations"
	"log"
	"time"
)

func FillTestDatabase(db *database.Database) error {
	log.Println("=== Starting FillTestDatabase ===")

	// Create tables
	log.Println("Creating tables...")
	err := migrations.CreatePointsTable(db)
	if err != nil {
		return err
	}

	err = migrations.CreateTravelsTable(db)
	if err != nil {
		return err
	}

	// Clear existing data
	log.Println("Clearing existing data...")
	ClearTestDatabase(db, "points")
	ClearTestDatabase(db, "travels")

	// Setup generator to generate about 1000 points
	// With n=63 and skip pattern (i%2==0, j%2==0), we get (63/2+1)^2 = 32^2 = 1024 points
	gf := generator.GeneratorFactory{}
	idGenerator := &generator.UUIDGenerator{}

	// squareSize = 3000 means distance between valid points is 6000 units
	// At speed 1000 units/hour, this gives ~6 hours travel time
	n := 63
	squareSize := 3000.0
	randFactor := 0.1 // 10% variation for realistic data

	g := gf.CreateGenerator(n, squareSize, randFactor, idGenerator)

	// Generate and insert points
	log.Println("Generating points...")
	pointDao := dao.NewPointDao(db)
	pointDbConsumer := generator.NewPointConsumer(pointDao, 100) // Buffer size 100

	err = g.GeneratePoints(pointDbConsumer)
	if err != nil {
		return err
	}

	err = pointDbConsumer.Flush()
	if err != nil {
		return err
	}

	// Retrieve points from database for travel generation
	log.Println("Retrieving points from database...")
	points, err := pointDao.SelectAll()
	if err != nil {
		return err
	}

	// Generate travels
	// Time interval: 2 months (2027-01-01 to 2027-03-01)
	// Rest time: 1 day (24 hours)
	// Speed: 1000 units/hour (gives ~6 hour travel for 6000 unit distance)
	log.Println("Generating travels...")
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 12, 1, 0, 0, 0, 0, time.UTC) // 2 months later
	speed := 1000.0
	restHours := 24 // 1 day

	travelDao := dao.NewTravelDao(db)
	travelDbConsumer := generator.NewTravelConsumer(travelDao, 500) // Buffer size 100

	err = g.GenerateTravels(points, fromDate, toDate, speed, restHours, travelDbConsumer)
	if err != nil {
		return err
	}

	err = travelDbConsumer.Flush()
	if err != nil {
		return err
	}

	// Count and log results
	log.Println("Counting results...")
	pointsCount, err := pointDao.Count()
	if err != nil {
		return err
	}

	travelsCount, err := travelDao.Count()
	if err != nil {
		return err
	}

	duration := toDate.Sub(fromDate)
	durationDays := int(duration.Hours() / 24)
	averageTravelTime := int(squareSize * 2 / speed)

	log.Printf("=== FillTestDatabase Complete ===")
	log.Printf("Total Points:  %d", pointsCount)
	log.Printf("Total Travels: %d", travelsCount)
	log.Printf("Time Period:   %s to %s (%d days)", fromDate.Format(time.DateOnly), toDate.Format(time.DateOnly), durationDays)
	log.Printf("Travel Config: ~%d hours travel time, %d hours rest", averageTravelTime, restHours)
	log.Printf("Randomness:    ±%d%% variation", int(randFactor*100))

	return nil
}
