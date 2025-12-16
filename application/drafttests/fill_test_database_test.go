//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"log"
	"testing"
)

func TestFillTestDatabase(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	dbFiller := &integration_tests.DatabaseFiller{}
	log.Println("=== Starting FillDatabase ===")
	err = dbFiller.FillDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	err = dbFiller.LogResults()
	if err != nil {
		t.Fatal(err)
	}
}
