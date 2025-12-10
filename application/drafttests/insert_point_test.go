package drafttests

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
	"testing"
)

func TestInsertPoint(t *testing.T) {

	// TODO reconfigure to connect to the test database
	// TODO clear test database before the test
	// TODO when clearing database additionally check for the database name to avoid clearing live database by accident

	pointDao := dao.NewPointDao(database.GetDatabase())

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
