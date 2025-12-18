//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"log"
	"testing"
	"time"
)

func TestSlowQueryTimeout(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.CloseConnection()

	//conn, err := db.GetConnection()
	//if err != nil {
	//	t.Fatal(err)
	//}

	pointDao := dao.NewPointDao(db)

	point1, err := pointDao.FindByCoordinates(6000, 6000)
	if err != nil {
		t.Fatal(err)
	}
	point2, err := pointDao.FindByCoordinates(60000, 60000)
	if err != nil {
		t.Fatal(err)
	}

	travelDao := dao.NewTravelDao(db)

	// Set timeout on the DAO - this will cancel queries that take too long
	secondsTimeout := 30
	travelDao.Timeout = time.Duration(secondsTimeout) * time.Second
	log.Printf("Set TravelDao timeout to: %v", travelDao.Timeout)

	var sequences []*tables.TransferSequence

	start := time.Now()

	elapsed := func() time.Duration {
		return time.Since(start).Round(time.Millisecond)
	}

	stopChan := make(chan bool)
	errChan := make(chan error, 1)

	go func() {
		sequences, err = travelDao.FindPathSimple3(&data.TravelFilter{
			Source:                      point1.ID,
			Destination:                 point2.ID,
			MaxWaitHoursBetweenTransits: 1,
			ArrivalTimeFrom:             util.ParseDate("2027-01-01"),
			ArrivalTimeTo:               util.ParseDate("2027-06-01"),
		})
		errChan <- err
		stopChan <- true
		close(stopChan)
	}()

	finished := false

	for i := 0; i < secondsTimeout+5; i++ {
		select {
		case finished = <-stopChan:
		default:
			log.Printf("[%6s]     .\n", elapsed())
			time.Sleep(1 * time.Second)
		}
		if finished {
			break
		}
	}

	finish := time.Now()
	duration := finish.Sub(start)

	if finished {
		queryErr := <-errChan
		if queryErr != nil {
			log.Printf("Query finished with error after %v: %v", duration, queryErr)
			if queryErr.Error() == "context deadline exceeded" {
				log.Println("✅ Query was successfully cancelled by timeout!")
				t.Logf("SUCCESS: Query timed out as expected after ~%d seconds", secondsTimeout)
			}
		} else {
			log.Printf("Query finished successfully after %v", duration)
		}
	} else {
		log.Printf("Query took too long: expected ~%d seconds, was %.2f seconds", secondsTimeout, duration.Seconds())
		t.Errorf("Query did not timeout as expected")
	}

	if sequences != nil {
		log.Printf("Sequences found: %d", len(sequences))
	} else {
		log.Println("No sequences found")
	}
}
