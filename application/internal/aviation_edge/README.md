# Aviation Edge API Client

A Go client library for interacting with the Aviation Edge API.

## API Key

Your API Key: `***` ( in .env file)
Account Level: Developer
Documentation: https://aviation-edge.com/developers/

## Installation

The client is already part of the project at `application/internal/aviation_edge/`

## Quick Start

```go
import "your-project/internal/aviation_edge"

// Initialize the client
client := aviation_edge.NewAviationEdgeApiClient("***")

// Get real-time flight information
flights, err := client.GetFlightTracker(aviation_edge.FlightTrackerParams{
    FlightIata: "AA100",
})
if err != nil {
    log.Fatal(err)
}

for _, flight := range flights {
    fmt.Printf("Flight %s from %s to %s\n",
        flight.Flight.IataNumber,
        flight.Departure.IataCode,
        flight.Arrival.IataCode)
}
```

## Available Methods

### Real-time Flight Tracking

- **GetFlightTracker(params)** - Get real-time flight information with custom parameters

### Flight Schedules

- **GetFlightSchedules(params)** - Get current timetable schedules with custom parameters
- **GetHistoricalSchedules(params)** - Get past schedules for historical dates
- **GetFutureSchedules(params)** - Get future schedules (must be > 1 week ahead)

### Routes

- **GetAirlineRoutes(params)** - Get airline routes with custom parameters

### Static Data

- **GetAirports(params)** - Get airport database information with custom parameters
- **GetAirlines(params)** - Get airline database information with custom parameters

### Other

- **GetAutocomplete(params)** - Search cities, airports, railway and bus stations

## Response Types

### FlightTrackerResponse
Contains real-time flight information including:
- Geography (latitude, longitude, altitude, direction)
- Speed (horizontal, vertical, ground status)
- Departure details (airport, gate, terminal, times)
- Arrival details (airport, gate, terminal, baggage, times)
- Aircraft information
- Airline information
- Flight number
- Status

### ScheduleResponse
Contains schedule information including:
- Type (departure/arrival)
- Status
- Departure and arrival details
- Airline and flight information
- Aircraft details
- Codeshared flight information

### RouteResponse
Contains route information including:
- Departure airport (IATA, ICAO, terminal, time)
- Arrival airport (IATA, ICAO, terminal, time)
- Airline codes
- Flight number
- Codeshares

### AirportResponse
Contains airport database information including:
- Airport ID and name
- IATA and ICAO codes
- Coordinates (latitude, longitude)
- Timezone and GMT offset
- Country information
- Phone and contact details

### AirlineResponse
Contains airline database information including:
- Airline ID and name
- IATA and ICAO codes
- Call sign
- Fleet size and age
- Founding year
- Hub code
- Country information

## Advanced Usage

### Typed Parameters

All generic methods use strongly-typed parameter structs for better IDE support and compile-time safety:

```go
// Example: Get flights from JFK to LAX on American Airlines
flights, err := client.GetFlightTracker(aviation_edge.FlightTrackerParams{
    DepIata:     "JFK",
    ArrIata:     "LAX",
    AirlineIata: "AA",
})

// Example: Get historical schedules for a specific date
schedules, err := client.GetHistoricalSchedules(aviation_edge.HistoricalSchedulesParams{
    Code:     "JFK",
    Type:     "departure",
    DateFrom: "2025-12-20",
    DateTo:   "2025-12-20",
})

// Example: Get future schedules (must be > 1 week ahead)
futureSchedules, err := client.GetFutureSchedules(aviation_edge.FutureSchedulesParams{
    IataCode: "JFK",
    Type:     "arrival",
    Date:     "2026-01-10",
})

// Example: Get routes by airline and departure airport
routes, err := client.GetAirlineRoutes(aviation_edge.AirlineRoutesParams{
    AirlineIata:   "AA",
    DepartureIata: "JFK",
})
```

### Parameter Structs Reference

**FlightTrackerParams** - Real-time flight tracking
- `FlightIata` - Specific flight IATA number (optional)
- `AirlineIata` - Filter by airline IATA code (optional)
- `DepIata` - Departure airport IATA code (optional)
- `ArrIata` - Arrival airport IATA code (optional)

**FlightSchedulesParams** - Current timetable
- `IataCode` - Airport IATA code (required)
- `Type` - "departure" or "arrival" (optional)

**HistoricalSchedulesParams** - Past schedules
- `Code` - Airport IATA code (required)
- `Type` - "departure" or "arrival" (required)
- `DateFrom` - Start date YYYY-MM-DD format (required)
- `DateTo` - End date YYYY-MM-DD format (optional, defaults to DateFrom)

**FutureSchedulesParams** - Future schedules (must be > 1 week ahead)
- `IataCode` - Airport IATA code (required)
- `Type` - "departure" or "arrival" (required)
- `Date` - Future date YYYY-MM-DD format (required, must be > 1 week ahead)
  - **Important:** The API enforces a minimum date requirement. The date must be at least 7 days in the future from the current date.
  - Example error: "date must be above 2026-01-03" (when called on 2025-12-28)
- `FlightNum` - Specific flight number (optional)
- `ArrIataCode` - Filter by arrival airport IATA code (optional)
- `DepIataCode` - Filter by departure airport IATA code (optional)

**AirlineRoutesParams** - Airline routes
- `AirlineIata` - Airline IATA code (optional)
- `DepartureIata` - Departure airport IATA code (optional)

**AirportsParams** - Airport database queries
- `CodeIataAirport` - Airport IATA code (optional)
- `CodeIso2Country` - Country ISO2 code (optional)

**AirlinesParams** - Airline database queries
- `CodeIataAirline` - Airline IATA code (optional)
- `CodeIso2Country` - Country ISO2 code (optional)

**AutocompleteParams** - Autocomplete search
- `Query` - Search query string (required)

### Error Handling

All methods return an error as the second return value. Always check for errors:

```go
airports, err := client.GetAirports(aviation_edge.AirportsParams{
    CodeIataAirport: "JFK",
})
if err != nil {
    log.Printf("Failed to get airport: %v", err)
    return
}
```

### Custom HTTP Client

You can customize the HTTP client settings:

```go
client := aviation_edge.NewAviationEdgeApiClient("your-api-key")
client.HTTPClient.Timeout = 60 * time.Second
```

## Common IATA Codes

**Major Airports:**
- JFK - John F. Kennedy International Airport (New York)
- LAX - Los Angeles International Airport
- LHR - London Heathrow Airport
- CDG - Charles de Gaulle Airport (Paris)
- DXB - Dubai International Airport

**Major Airlines:**
- AA - American Airlines
- DL - Delta Air Lines
- UA - United Airlines
- BA - British Airways
- LH - Lufthansa

## API Limitations

Developer Account:
- Check your personal dashboard for usage limits
- API endpoint: https://aviation-edge.com/v2/public/
- All requests require the API key parameter

## Support

For API issues or questions:
- Email: support@aviation-edge.com
- Documentation: https://aviation-edge.com/developers/
- Personal Dashboard: Login with giedriustum@gmail.com

## Example File

See `example_usage.go.example` for comprehensive usage examples of all available methods.
