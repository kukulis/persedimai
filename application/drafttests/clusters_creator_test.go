package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/migrations"
	"testing"
)

func TestClustersCreator(t *testing.T) {

	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}
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

	clustersCreator := migrations.NewClustersCreator(db)

	err = clustersCreator.CreateClustersTables()
	if err != nil {
		t.Fatal(err)
	}
	err = clustersCreator.InsertClustersDatas()
	if err != nil {
		t.Fatal(err)
	}
}
