package dao

import (
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
)

type TravelDao struct {
	database *database.Database
}

func NewTravelDao(database *database.Database) *TravelDao {
	return &TravelDao{database: database}
}

func (td *TravelDao) insert(t *tables.Travel) {
	// TODO
}

func (td *TravelDao) upsert([]*tables.Travel) int {
	// TODO
	return 0
}
