package dao

import (
	"database/sql"
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type AirportsDao struct {
	database *database.Database
}

func NewAirportsDao(database *database.Database) *AirportsDao {
	return &AirportsDao{
		database: database,
	}
}

func (*AirportsDao) GetTableFields() []string {
	return []string{
		"airport_id",
		"name_airport",
		"code_iata_airport",
		"code_icao_airport",
		"name_translations",
		"latitude_airport",
		"longitude_airport",
		"geoname_id",
		"timezone",
		"gmt",
		"phone",
		"name_country",
		"code_iso2_country",
		"code_iata_city",
	}
}

// Upsert inserts or updates airports in the database
// Builds a single SQL statement with multiple value rows
func (dao *AirportsDao) Upsert(airports []*aviation_edge.AirportResponse) error {
	if len(airports) == 0 {
		return nil
	}

	connection, err := dao.database.GetConnection()
	if err != nil {
		return err
	}

	// Build value lines for each airport
	lines := make([]string, len(airports))
	for i, airport := range airports {
		line := dao.buildValueLine(airport)
		lines[i] = line
	}

	valuesSubSql := strings.Join(lines, ",\n")
	fieldsSubSql := strings.Join(dao.GetTableFields(), ",\n")
	updatesArray := util.ArrayMap(dao.GetTableFields(), func(column string) string {
		return fmt.Sprintf("%s = VALUES(%s)", column, column)
	})
	updatesSubSql := strings.Join(updatesArray, ",\n")

	sqlQuery := `INSERT INTO airports ( ` + fieldsSubSql + ` ) VALUES ` +
		valuesSubSql + ` ON DUPLICATE KEY UPDATE ` + updatesSubSql + `,
		updated_at = CURRENT_TIMESTAMP`

	_, err = connection.Exec(sqlQuery)

	if err != nil {
		return errors.New(err.Error() + " for sqlQuery " + sqlQuery)
	}

	return nil
}

// buildValueLine constructs a single value row for the INSERT statement
func (dao *AirportsDao) buildValueLine(airport *aviation_edge.AirportResponse) string {
	quotedValuesArray := util.ArrayMap(dao.ToArray(airport), util.QuoteStringOrNull)

	return "(" + strings.Join(quotedValuesArray, ",") + ")"
}

func (dao *AirportsDao) ToArray(ar *aviation_edge.AirportResponse) []string {
	// Convert name_translations map to JSON string
	nameTranslationsJSON := ""
	if ar.NameTranslations != nil && len(ar.NameTranslations) > 0 {
		jsonBytes, err := json.Marshal(ar.NameTranslations)
		if err == nil {
			nameTranslationsJSON = string(jsonBytes)
		}
	}

	return []string{
		strconv.Itoa(ar.AirportID),
		ar.NameAirport,
		ar.CodeIataAirport,
		ar.CodeIcaoAirport,
		nameTranslationsJSON,
		fmt.Sprintf("%f", ar.LatitudeAirport),
		fmt.Sprintf("%f", ar.LongitudeAirport),
		ar.GeonameID,
		ar.Timezone,
		ar.GMT,
		ar.Phone,
		ar.NameCountry,
		ar.CodeIso2Country,
		ar.CodeIataCity,
	}
}

// Get retrieves a single airport by IATA code
func (dao *AirportsDao) Get(iataCode string) (*aviation_edge.AirportResponse, error) {
	connection, err := dao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	fields := dao.GetTableFields()
	fieldsSubSql := strings.Join(fields, ", ")
	sqlQuery := fmt.Sprintf("SELECT %s FROM airports WHERE code_iata_airport = ? LIMIT 1", fieldsSubSql)

	row := connection.QueryRow(sqlQuery, iataCode)

	airport := &aviation_edge.AirportResponse{}

	// Use sql.NullString for nullable fields
	var (
		codeIcaoAirport   sql.NullString
		nameTranslations  sql.NullString
		geonameID         sql.NullString
		timezone          sql.NullString
		gmt               sql.NullString
		phone             sql.NullString
		codeIataCity      sql.NullString
	)

	err = row.Scan(
		&airport.AirportID,
		&airport.NameAirport,
		&airport.CodeIataAirport,
		&codeIcaoAirport,
		&nameTranslations,
		&airport.LatitudeAirport,
		&airport.LongitudeAirport,
		&geonameID,
		&timezone,
		&gmt,
		&phone,
		&airport.NameCountry,
		&airport.CodeIso2Country,
		&codeIataCity,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Assign nullable fields
	airport.CodeIcaoAirport = codeIcaoAirport.String
	airport.GeonameID = geonameID.String
	airport.Timezone = timezone.String
	airport.GMT = gmt.String
	airport.Phone = phone.String
	airport.CodeIataCity = codeIataCity.String

	// Parse name_translations JSON
	if nameTranslations.Valid && nameTranslations.String != "" {
		var translations map[string]string
		if err := json.Unmarshal([]byte(nameTranslations.String), &translations); err == nil {
			airport.NameTranslations = translations
		}
	}

	return airport, nil
}

// GetAll retrieves all airports from the database
func (dao *AirportsDao) GetAll() ([]*aviation_edge.AirportResponse, error) {
	connection, err := dao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	fields := dao.GetTableFields()
	fieldsSubSql := strings.Join(fields, ", ")
	sqlQuery := fmt.Sprintf("SELECT %s FROM airports ORDER BY code_iata_airport", fieldsSubSql)

	rows, err := connection.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []*aviation_edge.AirportResponse
	for rows.Next() {
		airport := &aviation_edge.AirportResponse{}

		// Use sql.NullString for nullable fields
		var (
			codeIcaoAirport   sql.NullString
			nameTranslations  sql.NullString
			geonameID         sql.NullString
			timezone          sql.NullString
			gmt               sql.NullString
			phone             sql.NullString
			codeIataCity      sql.NullString
		)

		err := rows.Scan(
			&airport.AirportID,
			&airport.NameAirport,
			&airport.CodeIataAirport,
			&codeIcaoAirport,
			&nameTranslations,
			&airport.LatitudeAirport,
			&airport.LongitudeAirport,
			&geonameID,
			&timezone,
			&gmt,
			&phone,
			&airport.NameCountry,
			&airport.CodeIso2Country,
			&codeIataCity,
		)
		if err != nil {
			return nil, err
		}

		// Assign nullable fields
		airport.CodeIcaoAirport = codeIcaoAirport.String
		airport.GeonameID = geonameID.String
		airport.Timezone = timezone.String
		airport.GMT = gmt.String
		airport.Phone = phone.String
		airport.CodeIataCity = codeIataCity.String

		// Parse name_translations JSON
		if nameTranslations.Valid && nameTranslations.String != "" {
			var translations map[string]string
			if err := json.Unmarshal([]byte(nameTranslations.String), &translations); err == nil {
				airport.NameTranslations = translations
			}
		}

		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airports, nil
}
