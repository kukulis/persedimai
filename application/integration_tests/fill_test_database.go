package integration_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/generator"
	"darbelis.eu/persedimai/migrations"
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/util"
	"log"
	"time"
)

type DatabaseFiller struct {
	db               *database.Database
	generator        *generator.Generator
	pointDao         *dao.PointDao
	pointDbConsumer  *generator.PointDbConsumer
	travelDao        *dao.TravelDao
	travelDbConsumer *generator.TravelDbConsumer

	fromDate   *time.Time
	toDate     *time.Time
	speed      float64
	restHours  int
	squareSize float64
	randFactor float64
	n          int
}

func (d *DatabaseFiller) FillDatabase(db *database.Database) error {
	d.db = db
	// Create tables
	log.Println("Creating tables...")
	err := migrations.CreatePointsTable(d.db)
	if err != nil {
		return err
	}

	err = migrations.CreateTravelsTable(d.db)
	if err != nil {
		return err
	}

	// Clear existing data
	log.Println("Clearing existing data...")
	ClearTestDatabase(d.db, "points")
	ClearTestDatabase(d.db, "travels")

	// Setup generator to generate about 1000 points
	// With n=63 and skip pattern (i%2==0, j%2==0), we get (63/2+1)^2 = 32^2 = 1024 points
	gf := generator.GeneratorFactory{}
	idGenerator := &generator.UUIDGenerator{}

	// squareSize = 3000 means distance between valid points is 6000 units
	// At speed 1000 units/hour, this gives ~6 hours travel time
	d.n = 63
	d.squareSize = 3000.0
	d.randFactor = 0.1

	d.generator = gf.CreateGenerator(d.n, d.squareSize, d.randFactor, idGenerator)

	// Generate and insert points
	log.Println("Generating points...")
	d.pointDao = dao.NewPointDao(d.db)
	d.pointDbConsumer = generator.NewPointConsumer(d.pointDao, 100)

	err = d.generator.GeneratePoints(d.pointDbConsumer)
	if err != nil {
		return err
	}

	err = d.pointDbConsumer.Flush()
	if err != nil {
		return err
	}

	// Retrieve points from database for travel generation
	log.Println("Retrieving points from database...")
	points, err := d.pointDao.SelectAll()
	if err != nil {
		return err
	}

	// Generate travels
	// Time interval: 2 months (2027-01-01 to 2027-03-01)
	// Rest time: 1 day (24 hours)
	// Speed: 1000 units/hour (gives ~6 hour travel for 6000 unit distance)
	log.Println("Generating travels...")
	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	d.fromDate = &fromDate
	toDate := time.Date(2027, 12, 1, 0, 0, 0, 0, time.UTC) // 2 months later
	d.toDate = &toDate
	d.speed = 1000.0
	d.restHours = 24

	d.travelDao = dao.NewTravelDao(d.db)
	d.travelDbConsumer = generator.NewTravelConsumer(d.travelDao, 500)

	err = d.generator.GenerateTravels(points, *d.fromDate, *d.toDate, d.speed, d.restHours, d.travelDbConsumer)
	if err != nil {
		return err
	}

	err = d.travelDbConsumer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (d *DatabaseFiller) LogResults() error {
	// Count and log results
	log.Println("Counting results...")
	pointsCount, err := d.pointDao.Count()
	if err != nil {
		return err
	}

	travelsCount, err := d.travelDao.Count()
	if err != nil {
		return err
	}

	duration := d.toDate.Sub(*d.fromDate)
	durationDays := int(duration.Hours() / 24)
	averageTravelTime := int(d.squareSize * 2 / d.speed)

	log.Printf("=== FillDatabase Complete ===")
	log.Printf("Total Points:  %d", pointsCount)
	log.Printf("Total Travels: %d", travelsCount)
	log.Printf("Time Period:   %s to %s (%d days)", d.fromDate.Format(time.DateOnly), d.toDate.Format(time.DateOnly), durationDays)
	log.Printf("Travel Config: ~%d hours travel time, %d hours rest", averageTravelTime, d.restHours)
	log.Printf("Randomness:    ±%d%% variation", int(d.randFactor*100))

	return nil
}

// FillHubsTravels may be called after FillDatabase only
func (d *DatabaseFiller) FillHubsTravels() error {
	allPoints, err := d.pointDao.SelectAll()
	if err != nil {
		return err
	}
	hubPoints := util.ArrayFilter(allPoints, func(p *tables.Point) bool { return p.X == 0 })
	log.Println("== Filling between hubs travels ==")

	err = d.generator.GenerateTravelsBetweenHubPoints(hubPoints, *d.fromDate, *d.toDate, d.speed, d.restHours, d.travelDbConsumer)
	if err != nil {
		return err
	}
	log.Println("== Filling between hubs and non hubs travels ==")
	err = d.generator.GenerateTravelsFromHubToNonHubPoints(hubPoints, allPoints, *d.fromDate, *d.toDate, d.speed, d.restHours, d.travelDbConsumer)
	if err != nil {
		return err
	}

	err = d.travelDbConsumer.Flush()
	if err != nil {
		return err
	}

	return nil
}
