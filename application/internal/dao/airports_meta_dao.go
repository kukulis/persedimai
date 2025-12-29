package dao

import (
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/tables"
	"database/sql"
)

type AirportsMetaDao struct {
	database *database.Database
}

func NewAirportsMetaDao(database *database.Database) *AirportsMetaDao {
	return &AirportsMetaDao{database: database}
}

// CreateTable creates the airports_meta table if it doesn't exist
func (dao *AirportsMetaDao) CreateTable() error {
	conn, err := dao.database.GetConnection()
	if err != nil {
		return err
	}

	sqlQuery := `CREATE TABLE IF NOT EXISTS airports_meta (
		-- Primary key
		airport_code VARCHAR(3) PRIMARY KEY COMMENT 'IATA airport code',

		-- Import tracking
		imported_from DATETIME COMMENT 'Start date of imported data',
		imported_to DATETIME COMMENT 'End date of imported data',

		-- Metadata
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'record creation timestamp',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record update timestamp'
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Metadata tracking for airport data imports'`

	_, err = conn.Exec(sqlQuery)

	return err
}

// Upsert inserts or updates an airport metadata record
// updateDates controls whether imported_from and imported_to should be updated on duplicate key
func (dao *AirportsMetaDao) Upsert(meta *tables.AirportMeta, updateDates bool) error {
	conn, err := dao.database.GetConnection()
	if err != nil {
		return err
	}

	var sqlQuery string
	if updateDates {
		sqlQuery = `INSERT INTO airports_meta (airport_code, imported_from, imported_to)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE
				imported_from = VALUES(imported_from),
				imported_to = VALUES(imported_to),
				updated_at = CURRENT_TIMESTAMP`
	} else {
		sqlQuery = `INSERT INTO airports_meta (airport_code, imported_from, imported_to)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE
				updated_at = CURRENT_TIMESTAMP`
	}

	_, err = conn.Exec(sqlQuery, meta.AirportCode, meta.ImportedFrom, meta.ImportedTo)

	return err
}

// Get retrieves airport metadata by airport code
func (dao *AirportsMetaDao) Get(airportCode string) (*tables.AirportMeta, error) {
	conn, err := dao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sqlQuery := `SELECT airport_code, imported_from, imported_to
		FROM airports_meta
		WHERE airport_code = ?
		LIMIT 1`

	row := conn.QueryRow(sqlQuery, airportCode)

	meta := &tables.AirportMeta{}
	var importedFrom, importedTo sql.NullTime

	err = row.Scan(
		&meta.AirportCode,
		&importedFrom,
		&importedTo,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Assign nullable time fields
	if importedFrom.Valid {
		meta.ImportedFrom = &importedFrom.Time
	}
	if importedTo.Valid {
		meta.ImportedTo = &importedTo.Time
	}

	return meta, nil
}
