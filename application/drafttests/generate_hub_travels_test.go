//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"log"
	"testing"
)

func TestFillTestDatabase_WithHubPoints(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	log.Println("=== Starting FillDatabase ===")
	dbFiller := &integration_tests.DatabaseFiller{}
	err = dbFiller.FillDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("=== Filling hubs travels ===")
	err = dbFiller.FillHubsTravels()
	if err != nil {
		t.Fatal(err)
	}

	err = dbFiller.LogResults()
	if err != nil {
		t.Fatal(err)
	}

}
