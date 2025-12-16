package performance_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/travel_finder"
	"log"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkFindPaths(b *testing.B) {
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

	points, err := pointDao.SelectAll()

	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 3, 1, 0, 0, 0, 0, time.UTC)

	for b.Loop() {
		point1 := points[rand.Intn(len(points))]
		point2 := points[rand.Intn(len(points))]

		filter := data.TravelFilter{
			Source:          point1.ID,
			Destination:     point2.ID,
			ArrivalTimeFrom: fromDate,
			ArrivalTimeTo:   toDate,
			TravelCount:     2,
		}

		paths, err := strategy.FindPaths(&filter)
		if err != nil {
			b.Fatal(err)
		}
		if len(paths) == 0 {
			//b.Errorf("Did not find any path between points  %s and %s",
			//	point1.BuildLocationKey(),
			//	point2.BuildLocationKey(),
			//)
			log.Printf("Did not find any path between points  %s and %s",
				point1.BuildLocationKey(),
				point2.BuildLocationKey(),
			)
		}
	}
}
