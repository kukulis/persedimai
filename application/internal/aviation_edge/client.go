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

// handleUnmarshalError handles JSON unmarshal failures with a consistent error handling strategy:
// 1. Check if it's an API error response (200 status with error JSON)
// 2. Check if it's valid JSON but unexpected structure (logs full response)
// 3. Return the original unmarshal error
// Parameters:
//   - body: the response body bytes
//   - unmarshalErr: the original unmarshal error
//   - context: string identifier for logging (e.g., "GetFlightTracker")
//   - responseType: description for error message (e.g., "flight tracker")
func (c *AviationEdgeApiClient) handleUnmarshalError(body []byte, unmarshalErr error, context, responseType string) error {
	// Check if it's actually an error response (API returns 200 with error JSON sometimes)
	if errResp := c.checkForErrorResponse(body); errResp != nil {
		return errResp
	}

	// Check if it's valid JSON but unexpected structure
	var anyJSON map[string]interface{}
	if jsonErr := json.Unmarshal(body, &anyJSON); jsonErr == nil {
		// Valid JSON but unexpected structure - log it
		logPath := logUnexpectedResponse(body, 200, context)
		return fmt.Errorf("API returned unexpected JSON format, response logged to %s", logPath)
	}

	return fmt.Errorf("failed to parse %s response: %w", responseType, unmarshalErr)
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
		return nil, c.handleUnmarshalError(body, err, "GetFlightTracker", "flight tracker")
	}

	return flights, nil
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
		return nil, c.handleUnmarshalError(body, err, "GetFlightSchedules", "schedule")
	}

	return schedules, nil
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
		return nil, c.handleUnmarshalError(body, err, "GetHistoricalSchedules", "historical schedule")
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
		return nil, c.handleUnmarshalError(body, err, "GetFutureSchedules", "future schedule")
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
		return nil, c.handleUnmarshalError(body, err, "GetAirlineRoutes", "route")
	}

	return routes, nil
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
		return nil, c.handleUnmarshalError(body, err, "GetAirports", "airport")
	}

	return airports, nil
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
		return nil, c.handleUnmarshalError(body, err, "GetAirlines", "airline")
	}

	return airlines, nil
}

// Other Methods

// GetAutocomplete queries cities, airports, railway and bus stations
func (c *AviationEdgeApiClient) GetAutocomplete(params AutocompleteParams) ([]byte, error) {
	urlStr := c.buildURL("autocomplete", toMap(params))

	return c.doRequest(urlStr)
}
