package migrations

import "darbelis.eu/persedimai/internal/database"

// CreateFlightSchedulesTable creates a table to store flight schedule data from Aviation Edge API
// The table structure matches the ScheduleResponse JSON format from vno_future_2026.json
func CreateFlightSchedulesTable(db *database.Database) error {
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	defer func() { _ = db.CloseConnection() }()

	sql := `CREATE TABLE IF NOT EXISTS flight_schedules (
		-- Primary key
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		-- Top-level fields
		type VARCHAR(32) NOT NULL COMMENT 'departure or arrival',
		status VARCHAR(32) NOT NULL COMMENT 'flight status (e.g., active, scheduled, cancelled)',

		-- Departure information
		dep_iata_code VARCHAR(3) NOT NULL COMMENT 'departure airport IATA code',
		dep_icao_code VARCHAR(4) COMMENT 'departure airport ICAO code',
		dep_terminal VARCHAR(16) COMMENT 'departure terminal',
		dep_gate VARCHAR(16) COMMENT 'departure gate',
		dep_delay VARCHAR(16) COMMENT 'departure delay in minutes',
		dep_scheduled_time DATETIME NOT NULL COMMENT 'scheduled departure time',
		dep_estimated_time DATETIME COMMENT 'estimated departure time',
		dep_actual_time DATETIME COMMENT 'actual departure time',
		dep_estimated_runway DATETIME COMMENT 'estimated runway departure time',
		dep_actual_runway DATETIME COMMENT 'actual runway departure time',

		-- Arrival information
		arr_iata_code VARCHAR(3) NOT NULL COMMENT 'arrival airport IATA code',
		arr_icao_code VARCHAR(4) COMMENT 'arrival airport ICAO code',
		arr_terminal VARCHAR(16) COMMENT 'arrival terminal',
		arr_gate VARCHAR(16) COMMENT 'arrival gate',
		arr_baggage VARCHAR(16) COMMENT 'baggage claim area',
		arr_delay VARCHAR(16) COMMENT 'arrival delay in minutes',
		arr_scheduled_time DATETIME NOT NULL COMMENT 'scheduled arrival time',
		arr_estimated_time DATETIME COMMENT 'estimated arrival time',
		arr_actual_time DATETIME COMMENT 'actual arrival time',
		arr_estimated_runway DATETIME COMMENT 'estimated runway arrival time',
		arr_actual_runway DATETIME COMMENT 'actual runway arrival time',

		-- Airline information
		airline_name VARCHAR(128) NOT NULL COMMENT 'airline name',
		airline_iata_code VARCHAR(2) NOT NULL COMMENT 'airline IATA code',
		airline_icao_code VARCHAR(3) COMMENT 'airline ICAO code',

		-- Flight information
		flight_number VARCHAR(16) NOT NULL COMMENT 'flight number',
		flight_iata_number VARCHAR(16) NOT NULL COMMENT 'flight IATA number (airline code + number)',
		flight_icao_number VARCHAR(16) COMMENT 'flight ICAO number',

		-- Aircraft information
		aircraft_reg_number VARCHAR(16) COMMENT 'aircraft registration number',
		aircraft_icao_code VARCHAR(8) COMMENT 'aircraft type ICAO code',
		aircraft_icao_code_hex VARCHAR(16) COMMENT 'aircraft ICAO code hex',
		aircraft_iata_code VARCHAR(8) COMMENT 'aircraft type IATA code',
		aircraft_model_code VARCHAR(16) COMMENT 'aircraft model code',
		aircraft_model_text VARCHAR(128) COMMENT 'aircraft model text/name',

		-- Codeshared flight information (optional)
		codeshared_airline_name VARCHAR(128) COMMENT 'codeshared airline name',
		codeshared_airline_iata VARCHAR(2) COMMENT 'codeshared airline IATA code',
		codeshared_airline_icao VARCHAR(3) COMMENT 'codeshared airline ICAO code',
		codeshared_flight_number VARCHAR(16) COMMENT 'codeshared flight number',
		codeshared_flight_iata VARCHAR(16) COMMENT 'codeshared flight IATA number',
		codeshared_flight_icao VARCHAR(16) COMMENT 'codeshared flight ICAO number',

		-- Metadata
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'record creation timestamp',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record update timestamp',

		-- Unique constraint to prevent duplicate schedules
		UNIQUE KEY unique_schedule (flight_iata_number, dep_scheduled_time, dep_iata_code, arr_iata_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Flight schedules from Aviation Edge API'`

	_, err = conn.Exec(sql)

	return err
}
