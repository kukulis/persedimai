package drafttests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/tables"
	"fmt"
	"testing"
)

func TestInsertPoint(t *testing.T) {
	db := database.GetDatabase("test")

	if !integration_tests.ClearTestDatabase(db, "points") {
		fmt.Println("Failed to clear database")
	}
	pointDao := dao.NewPointDao(db)

	err := pointDao.InsertMany([]*tables.Point{
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
