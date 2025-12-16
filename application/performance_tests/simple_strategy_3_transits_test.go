package performance_tests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/travel_finder"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestSimpleStrategy3Transits(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	travelDao := dao.NewTravelDao(db)

	skipFilling := false

	if !skipFilling {
		dbFiller := integration_tests.DatabaseFiller{}

		t.Log("Filling test database...")
		err = dbFiller.FillDatabase(db)
		if err != nil {
			t.Fatal(err)
		}

		err = dbFiller.FillHubsTravels()
		if err != nil {
			t.Fatal(err)
		}

		err = dbFiller.LogResults()
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Filling test database finished.")
	}

	pointDao := dao.NewPointDao(db)
	strategy := travel_finder.NewSimpleTravelSearchStrategy(travelDao)

	points, err := pointDao.SelectAll()
	pointMap := make(map[string]*tables.Point)
	for _, point := range points {
		pointMap[point.ID] = point
	}

	fromDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)

	point1 := points[rand.Intn(len(points))]
	point2 := points[rand.Intn(len(points))]

	filter := data.NewTravelFilter(point1.ID, point2.ID, fromDate, toDate, 2)

	/*
		2025/12/16 15:03:36 Found path: Travel Path (2 transfer(s), Duration: 3280h6m52s)
		  1. 126000.00000_90000.00000 → 0.00000_90000.00000
		     Depart: 2027-01-13 04:56:26
		     Arrive: 2027-01-18 05:16:11
		     Duration: 120h19m45s
		  2. 0.00000_90000.00000 → 0.00000_186000.00000
		     Depart: 2027-05-26 00:41:36
		     Arrive: 2027-05-29 21:03:18
		     Duration: 92h21m42s

	*/

	//point1id := "b13659eb-703a-41c9-bbbf-3cbe043ca0d1"
	//point2id := "6613d042-2a79-4ea7-9fa6-f86090fc948d"
	//filter := data.TravelFilter{
	//	Source:          point1id,
	//	Destination:     point2id,
	//	ArrivalTimeFrom: fromDate,
	//	ArrivalTimeTo:   toDate,
	//	TravelCount:     2,
	//}

	paths, err := strategy.FindPaths(filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) == 0 {
		t.Fatal("no paths found")
	}

	pointGetter := data.NewMapPointGetter(pointMap)

	log.Printf("Found %d paths", len(paths))

	for _, path := range paths {
		log.Println("Found path: " + path.ToString(pointGetter))
	}

}
