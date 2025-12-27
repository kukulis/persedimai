package aviation_edge

import (
	"fmt"
	"log"
)

// Example usage of the AviationEdgeApiClient
// This file demonstrates how to use the API client methods

func ExampleUsage(apiKey string) {
	// Initialize the client with your API key
	client := NewAviationEdgeApiClient(apiKey)

	// Example 1: Get real-time flight information by flight number
	flights, err := client.GetFlightTracker(FlightTrackerParams{
		FlightIata: "AA100",
	})
	if err != nil {
		log.Printf("Error getting flight: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 1: Flight by Number ===\n")
		fmt.Printf("Found %d flights for AA100\n", len(flights))
		for i, flight := range flights {
			if i >= 5 {
				fmt.Printf("And Other ... \n")
				break
			}
			fmt.Printf("  Flight: %s | From: %s To: %s | Status: %s\n",
				flight.Flight.IataNumber,
				flight.Departure.IataCode,
				flight.Arrival.IataCode,
				flight.Status)
			fmt.Printf("  Altitude: %.0f ft | Speed: %.0f km/h\n",
				flight.Geography.Altitude,
				flight.Speed.Horizontal)
		}
	}

	short := true
	if short {
		return
	}

	// Example 2: Get all flights from a specific airline
	airlineFlights, err := client.GetFlightTracker(FlightTrackerParams{
		AirlineIata: "AA",
	})
	if err != nil {
		log.Printf("Error getting airline flights: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 2: Flights by Airline ===\n")
		fmt.Printf("Found %d American Airlines flights\n", len(airlineFlights))
		for i, flight := range airlineFlights {
			if i >= 5 {
				fmt.Printf("  ... and %d more flights\n", len(airlineFlights)-5)
				break
			}
			fmt.Printf("  %s: %s -> %s (Status: %s)\n",
				flight.Flight.IataNumber,
				flight.Departure.IataCode,
				flight.Arrival.IataCode,
				flight.Status)
		}
	}

	// Example 3: Get all departing flights from an airport
	departingFlights, err := client.GetFlightTracker(FlightTrackerParams{
		DepIata: "JFK",
	})
	if err != nil {
		log.Printf("Error getting departing flights: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 3: Departing Flights from JFK ===\n")
		fmt.Printf("Found %d departing flights\n", len(departingFlights))
		for i, flight := range departingFlights {
			if i >= 5 {
				fmt.Printf("  ... and %d more flights\n", len(departingFlights)-5)
				break
			}
			fmt.Printf("  %s to %s - Departure: %s\n",
				flight.Flight.IataNumber,
				flight.Arrival.IataCode,
				flight.Departure.ScheduledTime)
		}
	}

	// Example 4: Get all arriving flights at an airport
	arrivingFlights, err := client.GetFlightTracker(FlightTrackerParams{
		ArrIata: "LAX",
	})
	if err != nil {
		log.Printf("Error getting arriving flights: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 4: Arriving Flights at LAX ===\n")
		fmt.Printf("Found %d arriving flights\n", len(arrivingFlights))
		for i, flight := range arrivingFlights {
			if i >= 5 {
				fmt.Printf("  ... and %d more flights\n", len(arrivingFlights)-5)
				break
			}
			fmt.Printf("  %s from %s - Arrival: %s\n",
				flight.Flight.IataNumber,
				flight.Departure.IataCode,
				flight.Arrival.ScheduledTime)
		}
	}

	// Example 5: Get airport schedule
	schedule, err := client.GetFlightSchedules(FlightSchedulesParams{
		IataCode: "JFK",
	})
	if err != nil {
		log.Printf("Error getting schedule: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 5: Airport Schedule ===\n")
		fmt.Printf("Found %d scheduled flights at JFK\n", len(schedule))
		for i, sched := range schedule {
			if i >= 5 {
				fmt.Printf("  ... and %d more scheduled flights\n", len(schedule)-5)
				break
			}
			fmt.Printf("  %s %s: %s -> %s (Type: %s, Status: %s)\n",
				sched.Airline.IataCode,
				sched.Flight.IataNumber,
				sched.Departure.IataCode,
				sched.Arrival.IataCode,
				sched.Type,
				sched.Status)
		}
	}

	// Example 6: Get airport schedule by type (departure or arrival)
	departureSchedule, err := client.GetFlightSchedules(FlightSchedulesParams{
		IataCode: "JFK",
		Type:     "departure",
	})
	if err != nil {
		log.Printf("Error getting departure schedule: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 6: Departure Schedule ===\n")
		fmt.Printf("Found %d departure flights at JFK\n", len(departureSchedule))
		for i, sched := range departureSchedule {
			if i >= 3 {
				fmt.Printf("  ... and %d more departures\n", len(departureSchedule)-3)
				break
			}
			fmt.Printf("  %s - Scheduled: %s, Gate: %s\n",
				sched.Flight.IataNumber,
				sched.Departure.ScheduledTime,
				sched.Departure.Gate)
		}
	}

	// Example 7: Get routes by airline
	routes, err := client.GetAirlineRoutes(AirlineRoutesParams{
		AirlineIata: "AA",
	})
	if err != nil {
		log.Printf("Error getting routes: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 7: Routes by Airline ===\n")
		fmt.Printf("Found %d routes for American Airlines\n", len(routes))
		for i, route := range routes {
			if i >= 5 {
				fmt.Printf("  ... and %d more routes\n", len(routes)-5)
				break
			}
			fmt.Printf("  %s: %s -> %s (Flight: %s)\n",
				route.AirlineIata,
				route.DepartureIata,
				route.ArrivalIata,
				route.FlightNumber)
		}
	}

	// Example 8: Get routes from a departure airport
	airportRoutes, err := client.GetAirlineRoutes(AirlineRoutesParams{
		DepartureIata: "JFK",
	})
	if err != nil {
		log.Printf("Error getting airport routes: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 8: Routes from Airport ===\n")
		fmt.Printf("Found %d routes departing from JFK\n", len(airportRoutes))
		for i, route := range airportRoutes {
			if i >= 5 {
				fmt.Printf("  ... and %d more routes\n", len(airportRoutes)-5)
				break
			}
			fmt.Printf("  JFK -> %s via %s %s\n",
				route.ArrivalIata,
				route.AirlineIata,
				route.FlightNumber)
		}
	}

	// Example 9: Get airport information by IATA code
	airports, err := client.GetAirports(AirportsParams{
		CodeIataAirport: "JFK",
	})
	if err != nil {
		log.Printf("Error getting airport: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 9: Airport Information ===\n")
		for _, airport := range airports {
			fmt.Printf("Airport: %s (%s)\n", airport.NameAirport, airport.CodeIataAirport)
			fmt.Printf("  ICAO: %s\n", airport.CodeIcaoAirport)
			fmt.Printf("  Location: %.4f, %.4f\n", airport.LatitudeAirport, airport.LongitudeAirport)
			fmt.Printf("  Country: %s (%s)\n", airport.NameCountry, airport.CodeIso2Country)
			fmt.Printf("  Timezone: %s (GMT%s)\n", airport.Timezone, airport.GMT)
		}
	}

	// Example 10: Get all airports in a country
	countryAirports, err := client.GetAirports(AirportsParams{
		CodeIso2Country: "US",
	})
	if err != nil {
		log.Printf("Error getting country airports: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 10: Airports in Country ===\n")
		fmt.Printf("Found %d airports in the US\n", len(countryAirports))
		for i, airport := range countryAirports {
			if i >= 10 {
				fmt.Printf("  ... and %d more airports\n", len(countryAirports)-10)
				break
			}
			fmt.Printf("  %s - %s (%s)\n",
				airport.CodeIataAirport,
				airport.NameAirport,
				airport.CodeIataCity)
		}
	}

	// Example 11: Get airline information by IATA code
	airlines, err := client.GetAirlines(AirlinesParams{
		CodeIataAirline: "AA",
	})
	if err != nil {
		log.Printf("Error getting airline: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 11: Airline Information ===\n")
		for _, airline := range airlines {
			fmt.Printf("Airline: %s (%s)\n", airline.NameAirline, airline.CodeIataAirline)
			fmt.Printf("  ICAO: %s\n", airline.CodeIcaoAirline)
			fmt.Printf("  Call Sign: %s\n", airline.CallSign)
			fmt.Printf("  Country: %s\n", airline.NameCountry)
			fmt.Printf("  Fleet Size: %d | Avg Age: %.1f years\n", airline.SizeAirline, airline.AgeFleet)
			fmt.Printf("  Founded: %d | Status: %s\n", airline.Founding, airline.StatusAirline)
		}
	}

	// Example 12: Get all airlines from a country
	countryAirlines, err := client.GetAirlines(AirlinesParams{
		CodeIso2Country: "US",
	})
	if err != nil {
		log.Printf("Error getting country airlines: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 12: Airlines in Country ===\n")
		fmt.Printf("Found %d airlines in the US\n", len(countryAirlines))
		for i, airline := range countryAirlines {
			if i >= 10 {
				fmt.Printf("  ... and %d more airlines\n", len(countryAirlines)-10)
				break
			}
			fmt.Printf("  %s - %s (Fleet: %d)\n",
				airline.CodeIataAirline,
				airline.NameAirline,
				airline.SizeAirline)
		}
	}

	// Example 13: Get historical schedules with custom parameters
	historicalSchedules, err := client.GetHistoricalSchedules(HistoricalSchedulesParams{
		Code:     "JFK",
		Type:     "departure",
		DateFrom: "2025-12-20",
		DateTo:   "2025-12-20",
	})
	if err != nil {
		log.Printf("Error getting historical schedules: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 13: Historical Schedules ===\n")
		fmt.Printf("Found %d historical departure flights from JFK on 2025-12-20\n", len(historicalSchedules))
		for i, sched := range historicalSchedules {
			if i >= 3 {
				fmt.Printf("  ... and %d more historical flights\n", len(historicalSchedules)-3)
				break
			}
			fmt.Printf("  %s to %s - Scheduled: %s (Status: %s)\n",
				sched.Flight.IataNumber,
				sched.Arrival.IataCode,
				sched.Departure.ScheduledTime,
				sched.Status)
		}
	}

	// Example 14: Get future schedules with custom parameters
	futureSchedules, err := client.GetFutureSchedules(FutureSchedulesParams{
		IataCode: "JFK",
		Type:     "arrival",
		Date:     "2025-12-30",
	})
	if err != nil {
		log.Printf("Error getting future schedules: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 14: Future Schedules ===\n")
		fmt.Printf("Found %d future arrival flights to JFK on 2025-12-30\n", len(futureSchedules))
		for i, sched := range futureSchedules {
			if i >= 3 {
				fmt.Printf("  ... and %d more future flights\n", len(futureSchedules)-3)
				break
			}
			fmt.Printf("  %s from %s - Scheduled: %s\n",
				sched.Flight.IataNumber,
				sched.Departure.IataCode,
				sched.Arrival.ScheduledTime)
		}
	}

	// Example 15: Use the generic GetFlightTracker method with custom parameters
	customFlights, err := client.GetFlightTracker(FlightTrackerParams{
		DepIata:     "JFK",
		ArrIata:     "LAX",
		AirlineIata: "AA",
	})
	if err != nil {
		log.Printf("Error getting custom flights: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 15: Custom Flight Search ===\n")
		fmt.Printf("Found %d American Airlines flights from JFK to LAX\n", len(customFlights))
		for _, flight := range customFlights {
			fmt.Printf("  %s - Departure: %s, Arrival: %s\n",
				flight.Flight.IataNumber,
				flight.Departure.ScheduledTime,
				flight.Arrival.ScheduledTime)
			fmt.Printf("    Position: %.4f, %.4f at %.0f ft\n",
				flight.Geography.Latitude,
				flight.Geography.Longitude,
				flight.Geography.Altitude)
		}
	}

	// Example 16: Get airline routes with custom parameters
	customRoutes, err := client.GetAirlineRoutes(AirlineRoutesParams{
		AirlineIata:   "AA",
		DepartureIata: "JFK",
	})
	if err != nil {
		log.Printf("Error getting custom routes: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 16: Custom Route Search ===\n")
		fmt.Printf("Found %d American Airlines routes from JFK\n", len(customRoutes))
		for i, route := range customRoutes {
			if i >= 5 {
				fmt.Printf("  ... and %d more routes\n", len(customRoutes)-5)
				break
			}
			fmt.Printf("  %s -> %s (Departure: %s, Arrival: %s)\n",
				route.DepartureIata,
				route.ArrivalIata,
				route.DepartureTime,
				route.ArrivalTime)
		}
	}

	// Example 17: Autocomplete search
	autocompleteData, err := client.GetAutocomplete(AutocompleteParams{Query: "New York"})
	if err != nil {
		log.Printf("Error getting autocomplete: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 17: Autocomplete Search ===\n")
		fmt.Printf("Autocomplete results for 'New York': %d bytes received\n", len(autocompleteData))
		fmt.Printf("Raw data (first 200 chars): %s...\n", string(autocompleteData[:min(200, len(autocompleteData))]))
	}

	// Example 18: Get airports with custom parameters
	customAirports, err := client.GetAirports(AirportsParams{
		CodeIataAirport: "JFK",
		CodeIso2Country: "US",
	})
	if err != nil {
		log.Printf("Error getting custom airports: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 18: Custom Airport Search ===\n")
		fmt.Printf("Found %d airports matching JFK in US\n", len(customAirports))
		for _, airport := range customAirports {
			fmt.Printf("  %s - %s (City: %s)\n",
				airport.CodeIataAirport,
				airport.NameAirport,
				airport.CodeIataCity)
		}
	}

	// Example 19: Get airlines with custom parameters
	customAirlines, err := client.GetAirlines(AirlinesParams{
		CodeIataAirline: "AA",
		CodeIso2Country: "US",
	})
	if err != nil {
		log.Printf("Error getting custom airlines: %v\n", err)
	} else {
		fmt.Printf("\n=== Example 19: Custom Airline Search ===\n")
		fmt.Printf("Found %d airlines matching AA in US\n", len(customAirlines))
		for _, airline := range customAirlines {
			fmt.Printf("  %s - %s (Type: %s, Status: %s)\n",
				airline.CodeIataAirline,
				airline.NameAirline,
				airline.Type,
				airline.StatusAirline)
		}
	}

	fmt.Printf("\n=== All Examples Completed ===\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
