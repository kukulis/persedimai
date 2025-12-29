package migrations

import "darbelis.eu/persedimai/internal/database"

// CreateAirportsTable creates a table to store airport data from Aviation Edge API
// The table structure matches the AirportResponse JSON format
func CreateAirportsTable(db *database.Database) error {
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	defer func() { _ = db.CloseConnection() }()

	sql := `CREATE TABLE IF NOT EXISTS airports (
		-- Primary key
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		-- Airport identification
		airport_id INT NOT NULL COMMENT 'Aviation Edge airport ID',
		name_airport VARCHAR(255) NOT NULL COMMENT 'airport name',
		code_iata_airport VARCHAR(3) NOT NULL COMMENT 'airport IATA code (3 letters)',
		code_icao_airport VARCHAR(4) COMMENT 'airport ICAO code (4 letters)',

		-- Translations (stored as JSON)
		name_translations JSON COMMENT 'airport name translations in different languages',

		-- Geographic coordinates
		latitude_airport DECIMAL(10, 7) NOT NULL COMMENT 'airport latitude',
		longitude_airport DECIMAL(10, 7) NOT NULL COMMENT 'airport longitude',
		geoname_id VARCHAR(16) COMMENT 'GeoNames database ID',

		-- Time zone information
		timezone VARCHAR(64) COMMENT 'timezone identifier (e.g., Europe/Vilnius)',
		gmt VARCHAR(16) COMMENT 'GMT offset (e.g., +2)',

		-- Contact information
		phone VARCHAR(64) COMMENT 'airport phone number',

		-- Location information
		name_country VARCHAR(128) NOT NULL COMMENT 'country name',
		code_iso2_country VARCHAR(3) NOT NULL COMMENT 'ISO 3166-1 alpha-2 country code',
		code_iata_city VARCHAR(3) COMMENT 'IATA code of the city',

		-- Metadata
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'record creation timestamp',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record update timestamp',

		-- Unique constraints
		UNIQUE KEY unique_airport_id (airport_id),
		UNIQUE KEY unique_iata_code (code_iata_airport),
		UNIQUE KEY unique_icao_code (code_icao_airport),

		-- Indexes for common queries
		INDEX idx_country (code_iso2_country),
		INDEX idx_city (code_iata_city),
		INDEX idx_coordinates (latitude_airport, longitude_airport)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Airport data from Aviation Edge API'`

	_, err = conn.Exec(sql)

	return err
}
