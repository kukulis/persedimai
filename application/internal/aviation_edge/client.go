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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// Flight Tracker Methods

// GetFlightTracker retrieves real-time flight information
// Supports filtering by flight number, airline, departure/arrival airports
func (c *AviationEdgeApiClient) GetFlightTracker(params map[string]string) ([]FlightTrackerResponse, error) {
	urlStr := c.buildURL("flights", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var flights []FlightTrackerResponse
	if err := json.Unmarshal(body, &flights); err != nil {
		fmt.Println(string(body))
		return nil, fmt.Errorf("failed to parse flight tracker response: %w", err)
	}

	return flights, nil
}

// GetFlightByNumber retrieves a specific flight by its IATA flight number
func (c *AviationEdgeApiClient) GetFlightByNumber(flightIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(map[string]string{
		"flight_iata": flightIata,
	})
}

// GetFlightsByAirline retrieves all flights for a specific airline
func (c *AviationEdgeApiClient) GetFlightsByAirline(airlineIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(map[string]string{
		"airline_iata": airlineIata,
	})
}

// GetFlightsByDepartureAirport retrieves all flights departing from a specific airport
func (c *AviationEdgeApiClient) GetFlightsByDepartureAirport(airportIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(map[string]string{
		"dep_iata": airportIata,
	})
}

// GetFlightsByArrivalAirport retrieves all flights arriving at a specific airport
func (c *AviationEdgeApiClient) GetFlightsByArrivalAirport(airportIata string) ([]FlightTrackerResponse, error) {
	return c.GetFlightTracker(map[string]string{
		"arr_iata": airportIata,
	})
}

// Schedule Methods

// GetFlightSchedules retrieves airport arrival/departure schedules
func (c *AviationEdgeApiClient) GetFlightSchedules(params map[string]string) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("timetable", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		return nil, fmt.Errorf("failed to parse schedule response: %w", err)
	}

	return schedules, nil
}

// GetAirportSchedule retrieves schedules for a specific airport
func (c *AviationEdgeApiClient) GetAirportSchedule(airportIata string) ([]ScheduleResponse, error) {
	return c.GetFlightSchedules(map[string]string{
		"iataCode": airportIata,
	})
}

// GetAirportScheduleByType retrieves schedules for a specific airport filtered by type (arrival/departure)
func (c *AviationEdgeApiClient) GetAirportScheduleByType(airportIata, scheduleType string) ([]ScheduleResponse, error) {
	return c.GetFlightSchedules(map[string]string{
		"iataCode": airportIata,
		"type":     scheduleType,
	})
}

// GetHistoricalSchedules retrieves historical schedules for past dates
// Parameters: code (airport IATA), type (departure/arrival), date_from, date_to (YYYY-MM-DD)
func (c *AviationEdgeApiClient) GetHistoricalSchedules(params map[string]string) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("flightsHistory", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		return nil, fmt.Errorf("failed to parse historical schedule response: %w", err)
	}

	return schedules, nil
}

// GetFutureSchedules retrieves future schedules based on future dates
// Parameters: iataCode (airport), type (departure/arrival), date (YYYY-MM-DD)
func (c *AviationEdgeApiClient) GetFutureSchedules(params map[string]string) ([]ScheduleResponse, error) {
	urlStr := c.buildURL("flightsFuture", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleResponse
	if err := json.Unmarshal(body, &schedules); err != nil {
		return nil, fmt.Errorf("failed to parse future schedule response: %w", err)
	}

	return schedules, nil
}

// Route Methods

// GetAirlineRoutes retrieves aviation routes filtered by flight, airline, or departure airport
func (c *AviationEdgeApiClient) GetAirlineRoutes(params map[string]string) ([]RouteResponse, error) {
	urlStr := c.buildURL("routes", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var routes []RouteResponse
	if err := json.Unmarshal(body, &routes); err != nil {
		return nil, fmt.Errorf("failed to parse route response: %w", err)
	}

	return routes, nil
}

// GetRoutesByAirline retrieves all routes for a specific airline
func (c *AviationEdgeApiClient) GetRoutesByAirline(airlineIata string) ([]RouteResponse, error) {
	return c.GetAirlineRoutes(map[string]string{
		"airlineIata": airlineIata,
	})
}

// GetRoutesByDepartureAirport retrieves all routes from a specific departure airport
func (c *AviationEdgeApiClient) GetRoutesByDepartureAirport(airportIata string) ([]RouteResponse, error) {
	return c.GetAirlineRoutes(map[string]string{
		"departureIata": airportIata,
	})
}

// Airport Methods

// GetAirports retrieves airport database information
func (c *AviationEdgeApiClient) GetAirports(params map[string]string) ([]AirportResponse, error) {
	urlStr := c.buildURL("airportDatabase", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var airports []AirportResponse
	if err := json.Unmarshal(body, &airports); err != nil {
		return nil, fmt.Errorf("failed to parse airport response: %w", err)
	}

	return airports, nil
}

// GetAirportByIataCode retrieves airport information by IATA code
func (c *AviationEdgeApiClient) GetAirportByIataCode(iataCode string) ([]AirportResponse, error) {
	return c.GetAirports(map[string]string{
		"codeIataAirport": iataCode,
	})
}

// GetAirportsByCountry retrieves all airports in a specific country
func (c *AviationEdgeApiClient) GetAirportsByCountry(countryIso2 string) ([]AirportResponse, error) {
	return c.GetAirports(map[string]string{
		"codeIso2Country": countryIso2,
	})
}

// Airline Methods

// GetAirlines retrieves airline database information
func (c *AviationEdgeApiClient) GetAirlines(params map[string]string) ([]AirlineResponse, error) {
	urlStr := c.buildURL("airlineDatabase", params)

	body, err := c.doRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var airlines []AirlineResponse
	if err := json.Unmarshal(body, &airlines); err != nil {
		return nil, fmt.Errorf("failed to parse airline response: %w", err)
	}

	return airlines, nil
}

// GetAirlineByIataCode retrieves airline information by IATA code
func (c *AviationEdgeApiClient) GetAirlineByIataCode(iataCode string) ([]AirlineResponse, error) {
	return c.GetAirlines(map[string]string{
		"codeIataAirline": iataCode,
	})
}

// GetAirlinesByCountry retrieves all airlines from a specific country
func (c *AviationEdgeApiClient) GetAirlinesByCountry(countryIso2 string) ([]AirlineResponse, error) {
	return c.GetAirlines(map[string]string{
		"codeIso2Country": countryIso2,
	})
}

// Other Methods

// GetAutocomplete queries cities, airports, railway and bus stations
func (c *AviationEdgeApiClient) GetAutocomplete(query string) ([]byte, error) {
	urlStr := c.buildURL("autocomplete", map[string]string{
		"query": query,
	})

	return c.doRequest(urlStr)
}
