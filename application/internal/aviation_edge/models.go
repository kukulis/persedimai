package aviation_edge

import (
	"encoding/json"
	"fmt"
)

// DelayValue handles both string and number delay values from different API endpoints
type DelayValue struct {
	Value string
}

func (d *DelayValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		d.Value = s
		return nil
	}

	// Try as number
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		d.Value = fmt.Sprintf("%.0f", n)
		return nil
	}

	// If both fail, set empty string
	d.Value = ""
	return nil
}

func (d DelayValue) String() string {
	return d.Value
}

// Flight Tracker Response Models

type FlightTrackerResponse struct {
	Geography Geography `json:"geography"`
	Speed     Speed     `json:"speed"`
	Departure Departure `json:"departure"`
	Arrival   Arrival   `json:"arrival"`
	Aircraft  Aircraft  `json:"aircraft"`
	Airline   Airline   `json:"airline"`
	Flight    Flight    `json:"flight"`
	System    System    `json:"system"`
	Status    string    `json:"status"`
}

type Geography struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Direction float64 `json:"direction"`
}

type Speed struct {
	Horizontal float64 `json:"horizontal"`
	IsGround   float64 `json:"isGround"`
	Vertical   float64 `json:"vertical"`
}

type Departure struct {
	IataCode        string     `json:"iataCode"`
	IcaoCode        string     `json:"icaoCode"`
	Terminal        string     `json:"terminal"`
	Gate            string     `json:"gate"`
	Delay           DelayValue `json:"delay"` // Handles both string and number
	ScheduledTime   string     `json:"scheduledTime"`
	EstimatedTime   string     `json:"estimatedTime"`
	ActualTime      string     `json:"actualTime"`
	EstimatedRunway string     `json:"estimatedRunway"`
	ActualRunway    string     `json:"actualRunway"`
}

type Arrival struct {
	IataCode        string     `json:"iataCode"`
	IcaoCode        string     `json:"icaoCode"`
	Terminal        string     `json:"terminal"`
	Gate            string     `json:"gate"`
	Baggage         string     `json:"baggage"`
	Delay           DelayValue `json:"delay"` // Handles both string and number
	ScheduledTime   string     `json:"scheduledTime"`
	EstimatedTime   string     `json:"estimatedTime"`
	ActualTime      string     `json:"actualTime"`
	EstimatedRunway string     `json:"estimatedRunway"`
	ActualRunway    string     `json:"actualRunway"`
}

type Aircraft struct {
	RegNumber   string `json:"regNumber"`
	IcaoCode    string `json:"icaoCode"`
	IcaoCodeHex string `json:"icaoCodeHex"`
	IataCode    string `json:"iataCode"`
	ModelCode   string `json:"modelCode"`
	ModelText   string `json:"modelText"`
}

type Airline struct {
	Name     string `json:"name"`
	IataCode string `json:"iataCode"`
	IcaoCode string `json:"icaoCode"`
}

type Flight struct {
	Number     string `json:"number"`
	IataNumber string `json:"iataNumber"`
	IcaoNumber string `json:"icaoNumber"`
}

type System struct {
	Updated int64  `json:"updated"`
	Squawk  string `json:"squawk"`
}

// Schedule Response Models

type ScheduleResponse struct {
	Type       string      `json:"type"`
	Status     string      `json:"status"`
	Departure  Departure   `json:"departure"`
	Arrival    Arrival     `json:"arrival"`
	Airline    Airline     `json:"airline"`
	Flight     Flight      `json:"flight"`
	Aircraft   Aircraft    `json:"aircraft,omitempty"`
	Codeshared *Codeshared `json:"codeshared,omitempty"`
}

func (sr *ScheduleResponse) GetAirlineName() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Airline.Name
}

func (sr *ScheduleResponse) GetAirlineIataCode() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Airline.IataCode
}

func (sr *ScheduleResponse) GetAirlineIcaoCode() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Airline.IcaoCode
}

func (sr *ScheduleResponse) GetFlightNumber() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Flight.Number
}

func (sr *ScheduleResponse) GetFlightIataNumber() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Flight.IataNumber
}

func (sr *ScheduleResponse) GetFlightIcaoNumber() string {
	if sr.Codeshared == nil {
		return ""
	}
	return sr.Codeshared.Flight.IataNumber
}

type Codeshared struct {
	Airline Airline `json:"airline"`
	Flight  Flight  `json:"flight"`
}

// Route Response Models

type RouteResponse struct {
	DepartureIata     string   `json:"departureIata"`
	DepartureIcao     string   `json:"departureIcao"`
	DepartureTerminal string   `json:"departureTerminal"`
	DepartureTime     string   `json:"departureTime"`
	ArrivalIata       string   `json:"arrivalIata"`
	ArrivalIcao       string   `json:"arrivalIcao"`
	ArrivalTerminal   string   `json:"arrivalTerminal"`
	ArrivalTime       string   `json:"arrivalTime"`
	AirlineIata       string   `json:"airlineIata"`
	AirlineIcao       string   `json:"airlineIcao"`
	FlightNumber      string   `json:"flightNumber"`
	Codeshares        []string `json:"codeshares,omitempty"`
}

// Airport Database Models

type AirportResponse struct {
	AirportID        int               `json:"airportId"`
	NameAirport      string            `json:"nameAirport"`
	CodeIataAirport  string            `json:"codeIataAirport"`
	CodeIcaoAirport  string            `json:"codeIcaoAirport"`
	NameTranslations map[string]string `json:"nameTranslations,omitempty"`
	LatitudeAirport  float64           `json:"latitudeAirport"`
	LongitudeAirport float64           `json:"longitudeAirport"`
	GeonameID        string            `json:"geonameId"`
	Timezone         string            `json:"timezone"`
	GMT              string            `json:"GMT"`
	Phone            string            `json:"phone"`
	NameCountry      string            `json:"nameCountry"`
	CodeIso2Country  string            `json:"codeIso2Country"`
	CodeIataCity     string            `json:"codeIataCity"`
}

// Airline Database Models

type AirlineResponse struct {
	AirlineID       int     `json:"airlineId"`
	NameAirline     string  `json:"nameAirline"`
	CodeIataAirline string  `json:"codeIataAirline"`
	CodeIcaoAirline string  `json:"codeIcaoAirline"`
	CallSign        string  `json:"callsign"`
	StatusAirline   string  `json:"statusAirline"`
	Type            string  `json:"type"`
	SizeAirline     int     `json:"sizeAirline"`
	AgeFleet        float64 `json:"ageFleet"`
	Founding        int     `json:"founding"`
	CodeHub         string  `json:"codeHub"`
	NameCountry     string  `json:"nameCountry"`
	CodeIso2Country string  `json:"codeIso2Country"`
}

// Error Response Models

// ErrorResponse represents an error response from the Aviation Edge API
// The API returns different error formats, this handles common patterns
type ErrorResponse struct {
	Error   string `json:"error"`   // Error type or short message
	Message string `json:"message"` // Detailed error message (most common)
	Code    int    `json:"code"`    // Error code (optional)
	Success *bool  `json:"success"` // Success flag (pointer to distinguish false from missing)
}

// IsError returns true if this looks like an error response
func (e ErrorResponse) IsError() bool {
	// If success is explicitly false, it's an error
	if e.Success != nil && !*e.Success {
		return true
	}
	// If error or message is set, it's an error
	return e.Error != "" || e.Message != ""
}

// ErrorMessage returns the best available error message
func (e ErrorResponse) ErrorMessage() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Error != "" {
		return e.Error
	}
	return "Unknown error"
}
