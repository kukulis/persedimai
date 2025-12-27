package aviation_edge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://aviation-edge.com/v2/public"
)

type AviationEdgeApiClient struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

func NewAviationEdgeApiClient(apiKey string) *AviationEdgeApiClient {
	return &AviationEdgeApiClient{
		APIKey:  apiKey,
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Helper methods

func (c *AviationEdgeApiClient) buildURL(endpoint string, params map[string]string) string {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", c.BaseURL, endpoint))
	q := u.Query()
	q.Set("key", c.APIKey)

	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (c *AviationEdgeApiClient) doRequest(urlStr string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Always read body - needed for both success and error responses
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes with error response parsing
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(body, resp.StatusCode, urlStr)
	}

	return body, nil
}

// handleErrorResponse processes non-200 API responses
// Attempts to parse as error JSON, falls back to logging unexpected formats
func (c *AviationEdgeApiClient) handleErrorResponse(body []byte, statusCode int, urlStr string) error {
	// Strategy: Try multiple parsing approaches in order of likelihood

	// 1. Try to parse as API error response
	var apiError ErrorResponse
	if err := json.Unmarshal(body, &apiError); err == nil {
		// Check if we got meaningful error data
		if apiError.IsError() {
			return fmt.Errorf("API error (%d): %s", statusCode, apiError.ErrorMessage())
		}
	}

	// 2. Check if it's valid JSON but not in error format
	var anyJSON map[string]interface{}
	if err := json.Unmarshal(body, &anyJSON); err == nil {
		// Valid JSON but unexpected structure
		logPath := logUnexpectedResponse(body, statusCode, urlStr)
		return fmt.Errorf("API returned unexpected JSON format (status %d), response logged to %s",
			statusCode, logPath)
	}

	// 3. Not JSON at all (plain text, HTML, etc.)
	logPath := logUnexpectedResponse(body, statusCode, urlStr)

	// Include snippet in error for immediate visibility
	snippet := string(body)
	if len(snippet) > 100 {
		snippet = snippet[:100] + "..."
	}

	return fmt.Errorf("API returned non-JSON response (status %d): %q, full response logged to %s",
		statusCode, snippet, logPath)
}

// checkForErrorResponse checks if the body contains an API error response
// even when status code is 200 (API returns 200 with error JSON for some errors)
// Returns error if the body contains an error response, nil otherwise
func (c *AviationEdgeApiClient) checkForErrorResponse(body []byte) error {
	var apiError ErrorResponse
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.IsError() {
		return fmt.Errorf("API error: %s", apiError.ErrorMessage())
	}
	return nil
}

// Flight Tracker Methods

// GetFlightTracker retrieves real-time flight information
// Supports filtering by flight number, airline, departure/arrival airports
func (c *AviationEdgeApiClient) GetFlightTracker(params FlightTrackerParams) ([]FlightTrackerResponse, error) {
	urlStr := c.buildURL("flights", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var flights []FlightTrackerResponse
	if err := json.Unmarshal(body, &flights); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetFlightTracker")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse flight tracker response: %w", err)
	}

	return flights, nil
}

// GetFlightByNumber retrieves a specific flight by its IATA flight number
func (c *AviationEdgeApiClient) GetFlightByNumber(flightIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(FlightTrackerParams{
		FlightIata: flightIata,
	})
}

// GetFlightsByAirline retrieves all flights for a specific airline
func (c *AviationEdgeApiClient) GetFlightsByAirline(airlineIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(FlightTrackerParams{
		AirlineIata: airlineIata,
	})
}

// GetFlightsByDepartureAirport retrieves all flights departing from a specific airport
func (c *AviationEdgeApiClient) GetFlightsByDepartureAirport(airportIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(FlightTrackerParams{
		DepIata: airportIata,
	})
}

// GetFlightsByArrivalAirport retrieves all flights arriving at a specific airport
func (c *AviationEdgeApiClient) GetFlightsByArrivalAirport(airportIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(FlightTrackerParams{
		ArrIata: airportIata,
	})
}

// Schedule Methods

// GetFlightSchedules retrieves airport arrival/departure schedules
func (c *AviationEdgeApiClient) GetFlightSchedules(params FlightSchedulesParams) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("timetable", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetFlightSchedules")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse schedule response: %w", err)
	}

	return schedules, nil
}

// GetAirportSchedule retrieves schedules for a specific airport
func (c *AviationEdgeApiClient) GetAirportSchedule(airportIata string) ([]ScheduleResponse, error) {
	return c.GetFlightSchedules(FlightSchedulesParams{
		IataCode: airportIata,
	})
}

// GetAirportScheduleByType retrieves schedules for a specific airport filtered by type (arrival/departure)
func (c *AviationEdgeApiClient) GetAirportScheduleByType(airportIata, scheduleType string) ([]ScheduleResponse, error) {
	return c.GetFlightSchedules(FlightSchedulesParams{
		IataCode: airportIata,
		Type:     scheduleType,
	})
}

// GetHistoricalSchedules retrieves historical schedules for past dates
// Parameters: code (airport IATA), type (departure/arrival), date_from, date_to (YYYY-MM-DD)
func (c *AviationEdgeApiClient) GetHistoricalSchedules(params HistoricalSchedulesParams) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("flightsHistory", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetHistoricalSchedules")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse historical schedule response: %w", err)
	}

	return schedules, nil
}

// GetFutureSchedules retrieves future schedules based on future dates
// Parameters: iataCode (airport), type (departure/arrival), date (YYYY-MM-DD)
func (c *AviationEdgeApiClient) GetFutureSchedules(params FutureSchedulesParams) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("flightsFuture", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetFutureSchedules")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse future schedule response: %w", err)
	}

	return schedules, nil
}

// Route Methods

// GetAirlineRoutes retrieves aviation routes filtered by flight, airline, or departure airport
func (c *AviationEdgeApiClient) GetAirlineRoutes(params AirlineRoutesParams) ([]RouteResponse, error) {
	urlStr := c.buildURL("routes", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var routes []RouteResponse
	if err := json.Unmarshal(body, &routes); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetAirlineRoutes")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse route response: %w", err)
	}

	return routes, nil
}

// GetRoutesByAirline retrieves all routes for a specific airline
func (c *AviationEdgeApiClient) GetRoutesByAirline(airlineIata string) ([]RouteResponse, error) {
	return c.GetAirlineRoutes(AirlineRoutesParams{
		AirlineIata: airlineIata,
	})
}

// GetRoutesByDepartureAirport retrieves all routes from a specific departure airport
func (c *AviationEdgeApiClient) GetRoutesByDepartureAirport(airportIata string) ([]RouteResponse, error) {
	return c.GetAirlineRoutes(AirlineRoutesParams{
		DepartureIata: airportIata,
	})
}

// Airport Methods

// GetAirports retrieves airport database information
func (c *AviationEdgeApiClient) GetAirports(params AirportsParams) ([]AirportResponse, error) {
	urlStr := c.buildURL("airportDatabase", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var airports []AirportResponse
	if err := json.Unmarshal(body, &airports); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetAirports")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse airport response: %w", err)
	}

	return airports, nil
}

// GetAirportByIataCode retrieves airport information by IATA code
func (c *AviationEdgeApiClient) GetAirportByIataCode(iataCode string) ([]AirportResponse, error) {
	return c.GetAirports(AirportsParams{
		CodeIataAirport: iataCode,
	})
}

// GetAirportsByCountry retrieves all airports in a specific country
func (c *AviationEdgeApiClient) GetAirportsByCountry(countryIso2 string) ([]AirportResponse, error) {
	return c.GetAirports(AirportsParams{
		CodeIso2Country: countryIso2,
	})
}

// Airline Methods

// GetAirlines retrieves airline database information
func (c *AviationEdgeApiClient) GetAirlines(params AirlinesParams) ([]AirlineResponse, error) {
	urlStr := c.buildURL("airlineDatabase", toMap(params))

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var airlines []AirlineResponse
	if err := json.Unmarshal(body, &airlines); err != nil {
		// Check if it's actually an error response (API returns 200 with error JSON sometimes)
		if errResp := c.checkForErrorResponse(body); errResp != nil {
			return nil, errResp
		}

		// Check if it's valid JSON but unexpected structure
		var anyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
			// Valid JSON but unexpected structure - log it
			logPath := logUnexpectedResponse(body, 200, "GetAirlines")
			return nil, fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
		}

		return nil, fmt.Errorf("failed to parse airline response: %w", err)
	}

	return airlines, nil
}

// GetAirlineByIataCode retrieves airline information by IATA code
func (c *AviationEdgeApiClient) GetAirlineByIataCode(iataCode string) ([]AirlineResponse, error) {
	return c.GetAirlines(AirlinesParams{
		CodeIataAirline: iataCode,
	})
}

// GetAirlinesByCountry retrieves all airlines from a specific country
func (c *AviationEdgeApiClient) GetAirlinesByCountry(countryIso2 string) ([]AirlineResponse, error) {
	return c.GetAirlines(AirlinesParams{
		CodeIso2Country: countryIso2,
	})
}

// Other Methods

// GetAutocomplete queries cities, airports, railway and bus stations
func (c *AviationEdgeApiClient) GetAutocomplete(params AutocompleteParams) ([]byte, error) {
	urlStr := c.buildURL("autocomplete", toMap(params))

	return c.doRequest(urlStr)
}
