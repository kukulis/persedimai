package dao

import (
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
	"errors"
	"fmt"
	"strings"
)

type TravelDao struct {
	database *database.Database
}

func NewTravelDao(database *database.Database) *TravelDao {
	return &TravelDao{database: database}
}

func (td *TravelDao) InsertMany(travels []*tables.Travel) error {
	connection, err := td.database.GetConnection()
	if err != nil {
		return err
	}

	lines := make([]string, len(travels))
	for i, travel := range travels {
		line := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')",
			database.MysqlRealEscapeString(travel.ID),
			database.MysqlRealEscapeString(travel.From),
			database.MysqlRealEscapeString(travel.To),
			travel.Departure.Format("2006-01-02 15:04:05"),
			travel.Arrival.Format("2006-01-02 15:04:05"))
		lines[i] = line
	}

	valuesSubSql := strings.Join(lines, ",\n")

	sql := "insert into travels (ID, from_point, to_point, departure, arrival) values " + valuesSubSql

	_, err = connection.Exec(sql)

	if err != nil {
		return errors.New(err.Error() + " for sql " + sql)
	}

	return nil
}

func (td *TravelDao) SelectAll() ([]*tables.Travel, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql := "SELECT id, from_point, to_point, departure, arrival FROM travels"
	rows, err := connection.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var travels []*tables.Travel
	for rows.Next() {
		travel := &tables.Travel{}
		err := rows.Scan(&travel.ID, &travel.From, &travel.To, &travel.Departure, &travel.Arrival)
		if err != nil {
			return nil, err
		}
		travels = append(travels, travel)
	}

	return travels, nil
}

func (td *TravelDao) Count() (int, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return 0, err
	}

	sql := "SELECT COUNT(*) FROM travels"
	var count int
	err = connection.QueryRow(sql).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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
