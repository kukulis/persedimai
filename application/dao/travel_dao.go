package dao

import (
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
)

type TravelDao struct {
	database *database.Database
}

func NewTravelDao(database *database.Database) *TravelDao {
	return &TravelDao{database: database}
}

func (td *TravelDao) Insert(t *tables.Travel) {
	// TODO
}

func (td *TravelDao) Upsert([]*tables.Travel) int {
	// TODO
	return 0
}

func (td *TravelDao) Search(filter *data.TravelFilter) []tables.Travel {
	// TODO build sql
	return nil
}
