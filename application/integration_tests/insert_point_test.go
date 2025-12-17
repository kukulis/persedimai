package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/migrations"
	"darbelis.eu/persedimai/internal/tables"
	"fmt"
	"testing"
)

func TestInsertPoint(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	if db.GetDatabaseName() != "test" {
		t.Errorf("db name should be test")
	}

	err = migrations.CreatePointsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	if !ClearTestDatabase(db, "points") {
		fmt.Println("Failed to clear database")
	}
	pointDao := dao.NewPointDao(db)

	err = pointDao.InsertMany([]*tables.Point{
		{
			ID:   "1",
			X:    10,
			Y:    10,
			Name: "pirmas",
		},
		{
			ID:   "2",
			X:    10,
			Y:    10,
			Name: "antras",
		},
	})

	if err != nil {
		t.Error(err)
	}
}
