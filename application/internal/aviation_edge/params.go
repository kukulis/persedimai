package aviation_edge

import (
	"reflect"
)

// FlightTrackerParams contains parameters for real-time flight tracking queries
type FlightTrackerParams struct {
	FlightIata  string `url:"flight_iata"`  // Specific flight IATA number (optional)
	AirlineIata string `url:"airline_iata"` // Filter by airline IATA code (optional)
	DepIata     string `url:"dep_iata"`     // Departure airport IATA code (optional)
	ArrIata     string `url:"arr_iata"`     // Arrival airport IATA code (optional)
}

// FlightSchedulesParams contains parameters for current timetable queries (timetable endpoint)
type FlightSchedulesParams struct {
	IataCode string `url:"iataCode"` // Airport IATA code (required)
	Type     string `url:"type"`     // "departure" or "arrival" (optional)
}

// HistoricalSchedulesParams contains parameters for historical schedule queries (flightsHistory endpoint)
type HistoricalSchedulesParams struct {
	Code     string `url:"code"`      // Airport IATA code (required)
	Type     string `url:"type"`      // "departure" or "arrival" (required)
	DateFrom string `url:"date_from"` // Start date YYYY-MM-DD format (required)
	DateTo   string `url:"date_to"`   // End date YYYY-MM-DD format (optional, defaults to DateFrom)
}

// FutureSchedulesParams contains parameters for future schedule queries (flightsFuture endpoint)
// Note: The API requires dates to be at least 1 week in the future
type FutureSchedulesParams struct {
	IataCode    string `url:"iataCode"`     // Airport IATA code (required)
	Type        string `url:"type"`         // "departure" or "arrival" (required)
	Date        string `url:"date"`         // Future date YYYY-MM-DD format (required, must be > 1 week ahead)
	FlightNum   string `url:"flight_num"`   // Specific flight number (optional)
	ArrIataCode string `url:"arr_iataCode"` // Filter by arrival airport IATA code (optional)
	DepIataCode string `url:"dep_iataCode"` // Filter by departure airport IATA code (optional)
}

// AirlineRoutesParams contains parameters for airline route queries
type AirlineRoutesParams struct {
	AirlineIata   string `url:"airlineIata"`   // Airline IATA code (optional)
	DepartureIata string `url:"departureIata"` // Departure airport IATA code (optional)
}

// AirportsParams contains parameters for airport database queries
type AirportsParams struct {
	CodeIataAirport string `url:"codeIataAirport"` // Airport IATA code (optional)
	CodeIso2Country string `url:"codeIso2Country"` // Country ISO2 code (optional)
}

// AirlinesParams contains parameters for airline database queries
type AirlinesParams struct {
	CodeIataAirline string `url:"codeIataAirline"` // Airline IATA code (optional)
	CodeIso2Country string `url:"codeIso2Country"` // Country ISO2 code (optional)
}

// AutocompleteParams contains parameters for autocomplete queries
type AutocompleteParams struct {
	Query string `url:"query"` // Search query (required)
}

// toMap converts a parameter struct to map[string]string for URL building
// It uses reflection to iterate over struct fields and their 'url' tags
// Empty string values are skipped to keep URLs clean
func toMap(params interface{}) map[string]string {
	result := make(map[string]string)

	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	// Handle nil or invalid values
	if !v.IsValid() {
		return result
	}

	// If it's a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Only process struct types
	if v.Kind() != reflect.Struct {
		return result
	}

	// Iterate over all fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get the url tag value
		urlTag := field.Tag.Get("url")
		if urlTag == "" {
			// If no url tag, skip this field
			continue
		}

		// Only handle string fields
		if value.Kind() == reflect.String {
			stringValue := value.String()
			// Skip empty strings
			if stringValue != "" {
				result[urlTag] = stringValue
			}
		}
	}

	return result
}
