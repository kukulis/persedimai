package dao

import (
	"context"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/tables"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// MAX_CLUSTERED_CONNECTION_TIME_RANGE defines available options for maximum connection time in hours for clustered search
var MAX_CLUSTERED_CONNECTION_TIME_RANGE = []int{2, 4, 8, 16, 32}

type TravelDao struct {
	database *database.Database
	Timeout  time.Duration // Query timeout (0 = no timeout)
}

func NewTravelDao(database *database.Database) *TravelDao {
	return &TravelDao{
		database: database,
		Timeout:  0, // No timeout by default
	}
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

	sqlQuery := "insert into travels (ID, from_point, to_point, departure, arrival) values " + valuesSubSql

	_, err = connection.Exec(sqlQuery)

	if err != nil {
		return errors.New(err.Error() + " for sqlQuery " + sqlQuery)
	}

	return nil
}

// SelectAll loads all travels from db. Should be avoided to call unless for testing purposes.
func (td *TravelDao) SelectAll() ([]*tables.Transfer, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sqlQuery := "SELECT id, from_point, to_point, departure, arrival FROM travels"
	rows, err := connection.Query(sqlQuery)
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

	sqlQuery := "SELECT COUNT(*) FROM travels"
	var count int
	err = connection.QueryRow(sqlQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (td *TravelDao) FindByID(id string) (*tables.Transfer, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sqlQuery := "SELECT id, from_point, to_point, departure, arrival FROM travels WHERE id = ?"
	travel := &tables.Transfer{}
	err = connection.QueryRow(sqlQuery, id).Scan(&travel.ID, &travel.From, &travel.To, &travel.Departure, &travel.Arrival)
	if err != nil {
		return nil, err
	}

	return travel, nil
}

func (td *TravelDao) FindByIDs(ids []string) ([]*tables.Transfer, error) {
	if len(ids) == 0 {
		return []*tables.Transfer{}, nil
	}

	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Build escaped ID values for IN clause
	escapedIDs := make([]string, len(ids))
	for i, id := range ids {
		escapedIDs[i] = fmt.Sprintf("'%s'", database.MysqlRealEscapeString(id))
	}

	sqlQuery := fmt.Sprintf("SELECT id, from_point, to_point, departure, arrival FROM travels WHERE id IN (%s)",
		strings.Join(escapedIDs, ","))

	rows, err := connection.Query(sqlQuery)
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

func (td *TravelDao) Insert(t *tables.Transfer) {
	// TODO
}

func (td *TravelDao) Upsert([]*tables.Transfer) int {
	// TODO
	return 0
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

	sqlQuery := `SELECT id, from_point, to_point, departure, arrival
	        FROM travels
	        WHERE from_point = ?
	          AND to_point = ?
	          AND arrival >= ?
	          AND arrival <= ?
	        ORDER BY departure ASC
	        LIMIT ?`

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery, filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.Limit)
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

	sqlQuery := fmt.Sprintf(`SELECT
	            t1.id, t1.from_point, t1.to_point, t1.departure, t1.arrival,
	            t2.id, t2.from_point, t2.to_point, t2.departure, t2.arrival
	        FROM travels t1
	        INNER JOIN travels t2 ON t1.to_point = t2.from_point
	        WHERE t1.from_point = '%s'
	          AND t2.to_point = '%s'
	          AND t2.departure >= t1.arrival
	          AND t2.departure <= DATE_ADD(t1.arrival, INTERVAL %d HOUR)
	          AND t2.arrival >= '%s'
	          AND t2.arrival <= '%s'
	        ORDER BY t2.arrival ASC
	        LIMIT %d`,
		database.MysqlRealEscapeString(filter.Source),
		database.MysqlRealEscapeString(filter.Destination),
		filter.MaxWaitHoursBetweenTransits,
		filter.ArrivalTimeFrom.Format(time.DateTime),
		filter.ArrivalTimeTo.Format(time.DateTime),
		filter.Limit)
	//// TODO remove after debug
	//log.Println("FindPathSimple2: sqlQuery = " + sqlQuery)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery)
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

// FindPathSimple3 finds paths with two intermediate stops (3 transfers)
// Returns all matching paths ordered by final arrival time (earliest first)
func (td *TravelDao) FindPathSimple3(filter *data.TravelFilter) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sqlQuery := fmt.Sprintf(`SELECT
	            t1.id, t1.from_point, t1.to_point, t1.departure, t1.arrival,
	            t2.id, t2.from_point, t2.to_point, t2.departure, t2.arrival,
	            t3.id, t3.from_point, t3.to_point, t3.departure, t3.arrival
	        FROM travels t1
	        INNER JOIN travels t2 ON t1.to_point = t2.from_point
	        INNER JOIN travels t3 ON t2.to_point = t3.from_point
	        WHERE t1.from_point = '%s'
	          AND t3.to_point = '%s'
	          AND t2.departure >= t1.arrival
	          AND t3.departure >= t2.arrival
	          AND t3.arrival >= '%s'
	          AND t3.arrival <= '%s'
	        ORDER BY t3.arrival ASC`,
		database.MysqlRealEscapeString(filter.Source),
		database.MysqlRealEscapeString(filter.Destination),
		filter.ArrivalTimeFrom.Format(time.DateTime),
		filter.ArrivalTimeTo.Format(time.DateTime))

	//// TODO remove after debug
	//log.Println("FindPathSimple3: sql = " + sqlQuery)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		transfer1 := &tables.Transfer{}
		transfer2 := &tables.Transfer{}
		transfer3 := &tables.Transfer{}

		err := rows.Scan(
			&transfer1.ID, &transfer1.From, &transfer1.To, &transfer1.Departure, &transfer1.Arrival,
			&transfer2.ID, &transfer2.From, &transfer2.To, &transfer2.Departure, &transfer2.Arrival,
			&transfer3.ID, &transfer3.From, &transfer3.To, &transfer3.Departure, &transfer3.Arrival,
		)
		if err != nil {
			return nil, err
		}

		// Create sequence with all three transfers in order
		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2, transfer3})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// FindPathClustered2 finds paths with one intermediate stop (2 transfers) using clustered data
// Returns all matching paths from the clustered_arrival_travels32 table
func (td *TravelDao) FindPathClustered2(fromPointID, toPointID string, arrivalTimeFrom, arrivalTimeTo time.Time, maxConnectionTimeHours int) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Calculate time clusters from dates using: floor(unix_timestamp(date)/3600)
	minCluster := arrivalTimeFrom.Unix() / 3600
	maxCluster := arrivalTimeTo.Unix() / 3600

	// Build table name dynamically based on max connection time
	tableName := fmt.Sprintf("clustered_arrival_travels%d", maxConnectionTimeHours)

	sqlQuery := fmt.Sprintf(`SELECT
	            c1.travel_id, c1.from_point, c1.to_point, c1.departure_cl, c1.arrival_cl,
	            c2.travel_id, c2.from_point, c2.to_point, c2.departure_cl, c2.arrival_cl
	        FROM %s c1
	        JOIN %s c2
	            ON c1.to_point = c2.from_point
	            AND c1.arrival_cl = c2.departure_cl
	        WHERE c1.from_point = ?
	          AND c2.to_point = ?
	          AND c2.arrival_cl >= ?
	          AND c2.arrival_cl <= ?`, tableName, tableName)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery, fromPointID, toPointID, minCluster, maxCluster)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		var t1ID, t1From, t1To string
		var t1DepCl, t1ArrCl int64
		var t2ID, t2From, t2To string
		var t2DepCl, t2ArrCl int64

		err := rows.Scan(&t1ID, &t1From, &t1To, &t1DepCl, &t1ArrCl,
			&t2ID, &t2From, &t2To, &t2DepCl, &t2ArrCl)
		if err != nil {
			return nil, err
		}

		// Convert cluster times back to time.Time (cluster * 3600 seconds)
		transfer1 := &tables.Transfer{
			ID:        t1ID,
			From:      t1From,
			To:        t1To,
			Departure: time.Unix(t1DepCl*3600, 0),
			Arrival:   time.Unix(t1ArrCl*3600, 0),
		}
		transfer2 := &tables.Transfer{
			ID:        t2ID,
			From:      t2From,
			To:        t2To,
			Departure: time.Unix(t2DepCl*3600, 0),
			Arrival:   time.Unix(t2ArrCl*3600, 0),
		}

		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// FindPathClustered3 finds paths with two intermediate stops (3 transfers) using clustered data
// Returns all matching paths from the clustered_arrival_travels32 table
func (td *TravelDao) FindPathClustered3(fromPointID, toPointID string, arrivalTimeFrom, arrivalTimeTo time.Time, maxConnectionTimeHours int) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Calculate time clusters from dates using: floor(unix_timestamp(date)/3600)
	minCluster := arrivalTimeFrom.Unix() / 3600
	maxCluster := arrivalTimeTo.Unix() / 3600

	// Build table name dynamically based on max connection time
	tableName := fmt.Sprintf("clustered_arrival_travels%d", maxConnectionTimeHours)

	sqlQuery := fmt.Sprintf(`SELECT
	            c1.travel_id, c1.from_point, c1.to_point, c1.departure_cl, c1.arrival_cl,
	            c2.travel_id, c2.from_point, c2.to_point, c2.departure_cl, c2.arrival_cl,
	            c3.travel_id, c3.from_point, c3.to_point, c3.departure_cl, c3.arrival_cl
	        FROM %s c1
	        JOIN %s c2
	            ON c1.to_point = c2.from_point
	            AND c1.arrival_cl = c2.departure_cl
	        JOIN %s c3
	            ON c2.to_point = c3.from_point
	            AND c2.arrival_cl = c3.departure_cl
	        WHERE c1.from_point = ?
	          AND c3.to_point = ?
	          AND c3.arrival_cl >= ?
	          AND c3.arrival_cl <= ?`, tableName, tableName, tableName)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery, fromPointID, toPointID, minCluster, maxCluster)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		var t1ID, t1From, t1To string
		var t1DepCl, t1ArrCl int64
		var t2ID, t2From, t2To string
		var t2DepCl, t2ArrCl int64
		var t3ID, t3From, t3To string
		var t3DepCl, t3ArrCl int64

		err := rows.Scan(&t1ID, &t1From, &t1To, &t1DepCl, &t1ArrCl,
			&t2ID, &t2From, &t2To, &t2DepCl, &t2ArrCl,
			&t3ID, &t3From, &t3To, &t3DepCl, &t3ArrCl)
		if err != nil {
			return nil, err
		}

		// Convert cluster times back to time.Time (cluster * 3600 seconds)
		transfer1 := &tables.Transfer{
			ID:        t1ID,
			From:      t1From,
			To:        t1To,
			Departure: time.Unix(t1DepCl*3600, 0),
			Arrival:   time.Unix(t1ArrCl*3600, 0),
		}
		transfer2 := &tables.Transfer{
			ID:        t2ID,
			From:      t2From,
			To:        t2To,
			Departure: time.Unix(t2DepCl*3600, 0),
			Arrival:   time.Unix(t2ArrCl*3600, 0),
		}
		transfer3 := &tables.Transfer{
			ID:        t3ID,
			From:      t3From,
			To:        t3To,
			Departure: time.Unix(t3DepCl*3600, 0),
			Arrival:   time.Unix(t3ArrCl*3600, 0),
		}

		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2, transfer3})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// FindPathClustered4 finds paths with three intermediate stops (4 transfers) using clustered data
// Returns all matching paths from the clustered_arrival_travels32 table
func (td *TravelDao) FindPathClustered4(fromPointID, toPointID string, arrivalTimeFrom, arrivalTimeTo time.Time, maxConnectionTimeHours int) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Calculate time clusters from dates using: floor(unix_timestamp(date)/3600)
	minCluster := arrivalTimeFrom.Unix() / 3600
	maxCluster := arrivalTimeTo.Unix() / 3600

	// Build table name dynamically based on max connection time
	tableName := fmt.Sprintf("clustered_arrival_travels%d", maxConnectionTimeHours)

	sqlQuery := fmt.Sprintf(`SELECT
	            c1.travel_id, c1.from_point, c1.to_point, c1.departure_cl, c1.arrival_cl,
	            c2.travel_id, c2.from_point, c2.to_point, c2.departure_cl, c2.arrival_cl,
	            c3.travel_id, c3.from_point, c3.to_point, c3.departure_cl, c3.arrival_cl,
	            c4.travel_id, c4.from_point, c4.to_point, c4.departure_cl, c4.arrival_cl
	        FROM %s c1
	        JOIN %s c2
	            ON c1.to_point = c2.from_point
	            AND c1.arrival_cl = c2.departure_cl
	        JOIN %s c3
	            ON c2.to_point = c3.from_point
	            AND c2.arrival_cl = c3.departure_cl
	        JOIN %s c4
	            ON c3.to_point = c4.from_point
	            AND c3.arrival_cl = c4.departure_cl
	        WHERE c1.from_point = ?
	          AND c4.to_point = ?
	          AND c4.arrival_cl >= ?
	          AND c4.arrival_cl <= ?`, tableName, tableName, tableName, tableName)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)
	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery, fromPointID, toPointID, minCluster, maxCluster)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		var t1ID, t1From, t1To string
		var t1DepCl, t1ArrCl int64
		var t2ID, t2From, t2To string
		var t2DepCl, t2ArrCl int64
		var t3ID, t3From, t3To string
		var t3DepCl, t3ArrCl int64
		var t4ID, t4From, t4To string
		var t4DepCl, t4ArrCl int64

		err := rows.Scan(&t1ID, &t1From, &t1To, &t1DepCl, &t1ArrCl,
			&t2ID, &t2From, &t2To, &t2DepCl, &t2ArrCl,
			&t3ID, &t3From, &t3To, &t3DepCl, &t3ArrCl,
			&t4ID, &t4From, &t4To, &t4DepCl, &t4ArrCl)
		if err != nil {
			return nil, err
		}

		// Convert cluster times back to time.Time (cluster * 3600 seconds)
		transfer1 := &tables.Transfer{
			ID:        t1ID,
			From:      t1From,
			To:        t1To,
			Departure: time.Unix(t1DepCl*3600, 0),
			Arrival:   time.Unix(t1ArrCl*3600, 0),
		}
		transfer2 := &tables.Transfer{
			ID:        t2ID,
			From:      t2From,
			To:        t2To,
			Departure: time.Unix(t2DepCl*3600, 0),
			Arrival:   time.Unix(t2ArrCl*3600, 0),
		}
		transfer3 := &tables.Transfer{
			ID:        t3ID,
			From:      t3From,
			To:        t3To,
			Departure: time.Unix(t3DepCl*3600, 0),
			Arrival:   time.Unix(t3ArrCl*3600, 0),
		}
		transfer4 := &tables.Transfer{
			ID:        t4ID,
			From:      t4From,
			To:        t4To,
			Departure: time.Unix(t4DepCl*3600, 0),
			Arrival:   time.Unix(t4ArrCl*3600, 0),
		}

		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2, transfer3, transfer4})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// FindPathClustered5 finds paths with four intermediate stops (5 transfers) using clustered data
// Returns all matching paths from the clustered_arrival_travels table
func (td *TravelDao) FindPathClustered5(fromPointID, toPointID string, arrivalTimeFrom, arrivalTimeTo time.Time, maxConnectionTimeHours int) ([]*tables.TransferSequence, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Calculate time clusters from dates using: floor(unix_timestamp(date)/3600)
	minCluster := arrivalTimeFrom.Unix() / 3600
	maxCluster := arrivalTimeTo.Unix() / 3600

	// Build table name dynamically based on max connection time
	tableName := fmt.Sprintf("clustered_arrival_travels%d", maxConnectionTimeHours)

	sqlQuery := fmt.Sprintf(`SELECT
	            c1.travel_id, c1.from_point, c1.to_point, c1.departure_cl, c1.arrival_cl,
	            c2.travel_id, c2.from_point, c2.to_point, c2.departure_cl, c2.arrival_cl,
	            c3.travel_id, c3.from_point, c3.to_point, c3.departure_cl, c3.arrival_cl,
	            c4.travel_id, c4.from_point, c4.to_point, c4.departure_cl, c4.arrival_cl,
	            c5.travel_id, c5.from_point, c5.to_point, c5.departure_cl, c5.arrival_cl
	        FROM %s c1
	        JOIN %s c2
	            ON c1.to_point = c2.from_point
	            AND c1.arrival_cl = c2.departure_cl
	        JOIN %s c3
	            ON c2.to_point = c3.from_point
	            AND c2.arrival_cl = c3.departure_cl
	        JOIN %s c4
	            ON c3.to_point = c4.from_point
	            AND c3.arrival_cl = c4.departure_cl
	        JOIN %s c5
	            ON c4.to_point = c5.from_point
	            AND c4.arrival_cl = c5.departure_cl
	        WHERE c1.from_point = '%s'
	          AND c5.to_point = '%s'
	          AND c5.arrival_cl >= %d
	          AND c5.arrival_cl <= %d`,
		tableName, tableName, tableName, tableName, tableName,
		database.MysqlRealEscapeString(fromPointID),
		database.MysqlRealEscapeString(toPointID),
		minCluster,
		maxCluster)

	log.Printf("FindPathClustered5 sql: %s", sqlQuery)

	// Add server-side timeout hint and execute query
	sqlQuery = td.database.AddTimeoutToQuery(sqlQuery, td.Timeout+2*time.Second)

	rows, err := td.executeQueryWithConfiguration(connection, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*tables.TransferSequence
	for rows.Next() {
		var t1ID, t1From, t1To string
		var t1DepCl, t1ArrCl int64
		var t2ID, t2From, t2To string
		var t2DepCl, t2ArrCl int64
		var t3ID, t3From, t3To string
		var t3DepCl, t3ArrCl int64
		var t4ID, t4From, t4To string
		var t4DepCl, t4ArrCl int64
		var t5ID, t5From, t5To string
		var t5DepCl, t5ArrCl int64

		err := rows.Scan(&t1ID, &t1From, &t1To, &t1DepCl, &t1ArrCl,
			&t2ID, &t2From, &t2To, &t2DepCl, &t2ArrCl,
			&t3ID, &t3From, &t3To, &t3DepCl, &t3ArrCl,
			&t4ID, &t4From, &t4To, &t4DepCl, &t4ArrCl,
			&t5ID, &t5From, &t5To, &t5DepCl, &t5ArrCl)
		if err != nil {
			return nil, err
		}

		// Convert cluster times back to time.Time (cluster * 3600 seconds)
		transfer1 := &tables.Transfer{
			ID:        t1ID,
			From:      t1From,
			To:        t1To,
			Departure: time.Unix(t1DepCl*3600, 0),
			Arrival:   time.Unix(t1ArrCl*3600, 0),
		}
		transfer2 := &tables.Transfer{
			ID:        t2ID,
			From:      t2From,
			To:        t2To,
			Departure: time.Unix(t2DepCl*3600, 0),
			Arrival:   time.Unix(t2ArrCl*3600, 0),
		}
		transfer3 := &tables.Transfer{
			ID:        t3ID,
			From:      t3From,
			To:        t3To,
			Departure: time.Unix(t3DepCl*3600, 0),
			Arrival:   time.Unix(t3ArrCl*3600, 0),
		}
		transfer4 := &tables.Transfer{
			ID:        t4ID,
			From:      t4From,
			To:        t4To,
			Departure: time.Unix(t4DepCl*3600, 0),
			Arrival:   time.Unix(t4ArrCl*3600, 0),
		}
		transfer5 := &tables.Transfer{
			ID:        t5ID,
			From:      t5From,
			To:        t5To,
			Departure: time.Unix(t5DepCl*3600, 0),
			Arrival:   time.Unix(t5ArrCl*3600, 0),
		}

		sequence := tables.NewTransferSequence([]*tables.Transfer{transfer1, transfer2, transfer3, transfer4, transfer5})
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}

// executeQueryWithConfiguration executes a query with optional timeout (both client and server side)
// Supports both direct SQL and parameterized queries
func (td *TravelDao) executeQueryWithConfiguration(connection *sql.DB, sqlQuery string, args ...interface{}) (*sql.Rows, error) {
	skipTimeout := true
	if skipTimeout || td.Timeout <= 0 {
		return connection.Query(sqlQuery, args...)
	}

	// Create context with timeout (client-side safety net)
	ctx, cancel := context.WithTimeout(context.Background(), td.Timeout)
	defer cancel()
	//log.Printf("Executing query with timeout: %v", td.Timeout)

	return connection.QueryContext(ctx, sqlQuery, args...)
}

// TravelTimeBounds represents the min/max time boundaries for travels
type TravelTimeBounds struct {
	MinDeparture time.Time
	MaxDeparture time.Time
	MinArrival   time.Time
	MaxArrival   time.Time
}

// GetMinMaxDeparture returns the minimum and maximum departure times
func (td *TravelDao) GetMinMaxDeparture() (minDeparture, maxDeparture time.Time, err error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	sqlQuery := "SELECT MIN(departure), MAX(departure) FROM travels"
	err = connection.QueryRow(sqlQuery).Scan(&minDeparture, &maxDeparture)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return minDeparture, maxDeparture, nil
}

// GetMinMaxArrival returns the minimum and maximum arrival times
func (td *TravelDao) GetMinMaxArrival() (minArrival, maxArrival time.Time, err error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	sqlQuery := "SELECT MIN(arrival), MAX(arrival) FROM travels"
	err = connection.QueryRow(sqlQuery).Scan(&minArrival, &maxArrival)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return minArrival, maxArrival, nil
}

// GetTimeBounds returns all time boundaries in a single query (more efficient)
func (td *TravelDao) GetTimeBounds() (*TravelTimeBounds, error) {
	connection, err := td.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sqlQuery := "SELECT MIN(departure), MAX(departure), MIN(arrival), MAX(arrival) FROM travels"
	bounds := &TravelTimeBounds{}
	err = connection.QueryRow(sqlQuery).Scan(&bounds.MinDeparture, &bounds.MaxDeparture, &bounds.MinArrival, &bounds.MaxArrival)
	if err != nil {
		return nil, err
	}

	return bounds, nil
}
