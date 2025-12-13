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

func (td *TravelDao) InsertMany(travels []*tables.Transfer) error {
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

func (td *TravelDao) SelectAll() ([]*tables.Transfer, error) {
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

	var travels []*tables.Transfer
	for rows.Next() {
		travel := &tables.Transfer{}
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

func (td *TravelDao) Insert(t *tables.Transfer) {
	// TODO
}

func (td *TravelDao) Upsert([]*tables.Transfer) int {
	// TODO
	return 0
}

func (td *TravelDao) Search(filter *data.TravelFilter) []tables.Transfer {
	// TODO build sql
	return nil
}

// FindPathSimple1 finds direct paths (1 transfer) from source to destination
// Returns all matching paths ordered by departure time (earliest first)
func (td *TravelDao) FindPathSimple1(filter *data.TravelFilter) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Convert int IDs to strings
	//sourceID := fmt.Sprintf("%d", filter.Source)
	//destID := fmt.Sprintf("%d", filter.Destination)

	sql := `SELECT id, from_point, to_point, departure, arrival
	        FROM travels
	        WHERE from_point = ?
	          AND to_point = ?
	          AND arrival >= ?
	          AND arrival <= ?
	        ORDER BY departure ASC`

	rows, err := connection.Query(sql, filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		transfer := &tables.Transfer{}
		err := rows.Scan(&transfer.ID, &transfer.From, &transfer.To, &transfer.Departure, &transfer.Arrival)
		if err != nil {
			return nil, err
		}
		// Each direct connection is its own sequence with 1 transfer
		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// FindPathSimple2 finds paths with one intermediate stop (2 transfers)
// Returns all matching paths ordered by final arrival time (earliest first)
func (td *TravelDao) FindPathSimple2(filter *data.TravelFilter) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Convert int IDs to strings
	//sourceID := fmt.Sprintf("%d", filter.Source)
	//destID := fmt.Sprintf("%d", filter.Destination)

	sql := `SELECT
	            t1.id, t1.from_point, t1.to_point, t1.departure, t1.arrival,
	            t2.id, t2.from_point, t2.to_point, t2.departure, t2.arrival
	        FROM travels t1
	        INNER JOIN travels t2 ON t1.to_point = t2.from_point
	        WHERE t1.from_point = ?
	          AND t2.to_point = ?
	          AND t2.departure >= t1.arrival
	          AND t2.arrival >= ?
	          AND t2.arrival <= ?
	        ORDER BY t2.arrival ASC`

	rows, err := connection.Query(sql, filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		transfer1 := &tables.Transfer{}
		transfer2 := &tables.Transfer{}

		err := rows.Scan(
			&transfer1.ID, &transfer1.From, &transfer1.To, &transfer1.Departure, &transfer1.Arrival,
			&transfer2.ID, &transfer2.From, &transfer2.To, &transfer2.Departure, &transfer2.Arrival,
		)
		if err != nil {
			return nil, err
		}

		// Create sequence with both transfers in order
		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}
