package dao

import (
	"database/sql"
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/util"
	"errors"
	"fmt"
	"strings"
)

type AviationEdgeFlightSchedulesDao struct {
	database *database.Database
}

func NewAviationEdgeFlightSchedulesDao(database *database.Database) *AviationEdgeFlightSchedulesDao {
	return &AviationEdgeFlightSchedulesDao{
		database: database,
	}
}

func (*AviationEdgeFlightSchedulesDao) GetTableFields() []string {
	return []string{
		"type",
		"status",
		"dep_iata_code",
		"dep_icao_code",
		"dep_terminal",
		"dep_gate",
		"dep_delay",
		"dep_scheduled_time",
		"dep_estimated_time",
		"dep_actual_time",
		"dep_estimated_runway",
		"dep_actual_runway",
		"arr_iata_code",
		"arr_icao_code",
		"arr_terminal",
		"arr_gate",
		"arr_baggage",
		"arr_delay",
		"arr_scheduled_time",
		"arr_estimated_time",
		"arr_actual_time",
		"arr_estimated_runway",
		"arr_actual_runway",
		"airline_name",
		"airline_iata_code",
		"airline_icao_code",
		"flight_number",
		"flight_iata_number",
		"flight_icao_number",
		"aircraft_reg_number",
		"aircraft_icao_code",
		"aircraft_icao_code_hex",
		"aircraft_iata_code",
		"aircraft_model_code",
		"aircraft_model_text",
		"codeshared_airline_name",
		"codeshared_airline_iata",
		"codeshared_airline_icao",
		"codeshared_flight_number",
		"codeshared_flight_iata",
		"codeshared_flight_icao",
	}
}

// UpsertFlightSchedules inserts or updates flight schedules in the database
// Builds a single SQL statement with multiple value rows
func (dao *AviationEdgeFlightSchedulesDao) UpsertFlightSchedules(schedules []*aviation_edge.ScheduleResponse) error {
	if len(schedules) == 0 {
		return nil
	}

	connection, err := dao.database.GetConnection()
	if err != nil {
		return err
	}

	// Build value lines for each schedule
	lines := make([]string, len(schedules))
	for i, schedule := range schedules {
		line := dao.buildValueLine(schedule)
		lines[i] = line
	}

	valuesSubSql := strings.Join(lines, ",\n")
	fieldsSubSql := strings.Join(dao.GetTableFields(), ",\n")
	updatesArray := util.ArrayMap(dao.GetTableFields(), func(column string) string { return fmt.Sprintf("%s = VALUES(%s)", column, column) })
	updatesSubSql := strings.Join(updatesArray, ",\n")

	sqlQuery := `INSERT INTO flight_schedules ( ` + fieldsSubSql + ` ) VALUES ` +
		valuesSubSql + ` ON DUPLICATE KEY UPDATE ` + updatesSubSql + `,
		updated_at = CURRENT_TIMESTAMP`

	_, err = connection.Exec(sqlQuery)

	if err != nil {
		return errors.New(err.Error() + " for sqlQuery " + sqlQuery)
	}

	return nil
}

// buildValueLine constructs a single value row for the INSERT statement
func (dao *AviationEdgeFlightSchedulesDao) buildValueLine(schedule *aviation_edge.ScheduleResponse) string {
	quotedValuesArray := util.ArrayMap(dao.ToArray(schedule), util.QuoteStringOrNull)

	return "(" + strings.Join(quotedValuesArray, ",") + ")"
}

func (dao *AviationEdgeFlightSchedulesDao) ToArray(sr *aviation_edge.ScheduleResponse) []string {
	return []string{
		sr.Type,
		sr.Status,

		// Departure fields
		sr.Departure.IataCode,
		sr.Departure.IcaoCode,
		sr.Departure.Terminal,
		sr.Departure.Gate,
		sr.Departure.Delay.Value,
		sr.Departure.ScheduledTime,
		sr.Departure.EstimatedTime,
		sr.Departure.ActualTime,
		sr.Departure.EstimatedRunway,
		sr.Departure.ActualRunway,

		// Arrival fields
		sr.Arrival.IataCode,
		sr.Arrival.IcaoCode,
		sr.Arrival.Terminal,
		sr.Arrival.Gate,
		sr.Arrival.Baggage,
		sr.Arrival.Delay.Value,
		sr.Arrival.ScheduledTime,
		sr.Arrival.EstimatedTime,
		sr.Arrival.ActualTime,
		sr.Arrival.EstimatedRunway,
		sr.Arrival.ActualRunway,

		// Airline fields
		sr.Airline.Name,
		sr.Airline.IataCode,
		sr.Airline.IcaoCode,

		// Flight fields
		sr.Flight.Number,
		sr.Flight.IataNumber,
		sr.Flight.IcaoNumber,

		// Aircraft fields
		sr.Aircraft.RegNumber,
		sr.Aircraft.IcaoCode,
		sr.Aircraft.IcaoCodeHex,
		sr.Aircraft.IataCode,
		sr.Aircraft.ModelCode,
		sr.Aircraft.ModelText,

		// Codeshared fields
		sr.GetAirlineName(),
		sr.GetAirlineIataCode(),
		sr.GetAirlineIcaoCode(),
		sr.GetFlightNumber(),
		sr.GetFlightIataNumber(),
		sr.GetFlightIcaoNumber(),
	}
}

// GetAll retrieves all flight schedules from the database
func (dao *AviationEdgeFlightSchedulesDao) GetAll() ([]*aviation_edge.ScheduleResponse, error) {
	connection, err := dao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	fields := dao.GetTableFields()
	fieldsSubSql := strings.Join(fields, ", ")
	sqlQuery := fmt.Sprintf("SELECT %s FROM flight_schedules ORDER BY id", fieldsSubSql)

	rows, err := connection.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*aviation_edge.ScheduleResponse
	for rows.Next() {
		schedule := &aviation_edge.ScheduleResponse{}

		// Use sql.NullString for nullable fields
		var (
			depIcaoCode, depTerminal, depGate, depDelay                              sql.NullString
			depEstimatedTime, depActualTime, depEstimatedRunway, depActualRunway    sql.NullString
			arrIcaoCode, arrTerminal, arrGate, arrBaggage, arrDelay                  sql.NullString
			arrEstimatedTime, arrActualTime, arrEstimatedRunway, arrActualRunway    sql.NullString
			airlineIcaoCode                                                          sql.NullString
			flightIcaoNumber                                                         sql.NullString
			aircraftRegNumber, aircraftIcaoCode, aircraftIcaoCodeHex                 sql.NullString
			aircraftIataCode, aircraftModelCode, aircraftModelText                   sql.NullString
			codesharedAirlineName, codesharedAirlineIata, codesharedAirlineIcao      sql.NullString
			codesharedFlightNumber, codesharedFlightIata, codesharedFlightIcao       sql.NullString
		)

		err := rows.Scan(
			&schedule.Type,
			&schedule.Status,
			&schedule.Departure.IataCode,
			&depIcaoCode,
			&depTerminal,
			&depGate,
			&depDelay,
			&schedule.Departure.ScheduledTime,
			&depEstimatedTime,
			&depActualTime,
			&depEstimatedRunway,
			&depActualRunway,
			&schedule.Arrival.IataCode,
			&arrIcaoCode,
			&arrTerminal,
			&arrGate,
			&arrBaggage,
			&arrDelay,
			&schedule.Arrival.ScheduledTime,
			&arrEstimatedTime,
			&arrActualTime,
			&arrEstimatedRunway,
			&arrActualRunway,
			&schedule.Airline.Name,
			&schedule.Airline.IataCode,
			&airlineIcaoCode,
			&schedule.Flight.Number,
			&schedule.Flight.IataNumber,
			&flightIcaoNumber,
			&aircraftRegNumber,
			&aircraftIcaoCode,
			&aircraftIcaoCodeHex,
			&aircraftIataCode,
			&aircraftModelCode,
			&aircraftModelText,
			&codesharedAirlineName,
			&codesharedAirlineIata,
			&codesharedAirlineIcao,
			&codesharedFlightNumber,
			&codesharedFlightIata,
			&codesharedFlightIcao,
		)
		if err != nil {
			return nil, err
		}

		// Assign nullable fields
		schedule.Departure.IcaoCode = depIcaoCode.String
		schedule.Departure.Terminal = depTerminal.String
		schedule.Departure.Gate = depGate.String
		schedule.Departure.Delay.Value = depDelay.String
		schedule.Departure.EstimatedTime = depEstimatedTime.String
		schedule.Departure.ActualTime = depActualTime.String
		schedule.Departure.EstimatedRunway = depEstimatedRunway.String
		schedule.Departure.ActualRunway = depActualRunway.String

		schedule.Arrival.IcaoCode = arrIcaoCode.String
		schedule.Arrival.Terminal = arrTerminal.String
		schedule.Arrival.Gate = arrGate.String
		schedule.Arrival.Baggage = arrBaggage.String
		schedule.Arrival.Delay.Value = arrDelay.String
		schedule.Arrival.EstimatedTime = arrEstimatedTime.String
		schedule.Arrival.ActualTime = arrActualTime.String
		schedule.Arrival.EstimatedRunway = arrEstimatedRunway.String
		schedule.Arrival.ActualRunway = arrActualRunway.String

		schedule.Airline.IcaoCode = airlineIcaoCode.String
		schedule.Flight.IcaoNumber = flightIcaoNumber.String

		schedule.Aircraft.RegNumber = aircraftRegNumber.String
		schedule.Aircraft.IcaoCode = aircraftIcaoCode.String
		schedule.Aircraft.IcaoCodeHex = aircraftIcaoCodeHex.String
		schedule.Aircraft.IataCode = aircraftIataCode.String
		schedule.Aircraft.ModelCode = aircraftModelCode.String
		schedule.Aircraft.ModelText = aircraftModelText.String

		// Populate Codeshared if any codeshared field is non-empty
		if codesharedAirlineName.Valid || codesharedAirlineIata.Valid || codesharedFlightIata.Valid {
			schedule.Codeshared = &aviation_edge.Codeshared{
				Airline: aviation_edge.Airline{
					Name:     codesharedAirlineName.String,
					IataCode: codesharedAirlineIata.String,
					IcaoCode: codesharedAirlineIcao.String,
				},
				Flight: aviation_edge.Flight{
					Number:     codesharedFlightNumber.String,
					IataNumber: codesharedFlightIata.String,
					IcaoNumber: codesharedFlightIcao.String,
				},
			}
		}

		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
