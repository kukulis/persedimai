package main

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")

	client := aviation_edge.NewAviationEdgeApiClient(apiKey)

	// Test historical schedules for JFK on a recent date
	fmt.Println("=== Testing Historical Schedules for JFK (2025-12-20) ===")
	schedules, err := client.GetHistoricalSchedules(aviation_edge.HistoricalSchedulesParams{
		Code:     "JFK",
		Type:     "departure",
		DateFrom: "2025-12-20",
		DateTo:   "2025-12-20",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Found %d departure schedules\n", len(schedules))
		for i, schedule := range schedules {
			if i >= 5 {
				fmt.Printf("... and %d more schedules\n", len(schedules)-5)
				break
			}
			fmt.Printf("  %s: %s -> %s (Status: %s, Scheduled: %s)\n",
				schedule.Flight.IataNumber,
				schedule.Departure.IataCode,
				schedule.Arrival.IataCode,
				schedule.Status,
				schedule.Departure.ScheduledTime)
		}
	}

	fmt.Println("\n=== Testing Future Schedules for JFK (2026-01-10) ===")
	futureSchedules, err := client.GetFutureSchedules(aviation_edge.FutureSchedulesParams{
		IataCode: "JFK",
		Type:     "departure",
		Date:     "2026-01-10",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Found %d future departure schedules\n", len(futureSchedules))
		for i, schedule := range futureSchedules {
			if i >= 5 {
				fmt.Printf("... and %d more schedules\n", len(futureSchedules)-5)
				break
			}
			fmt.Printf("  %s: %s -> %s (Scheduled: %s)\n",
				schedule.Flight.IataNumber,
				schedule.Departure.IataCode,
				schedule.Arrival.IataCode,
				schedule.Departure.ScheduledTime)
		}
	}

	fmt.Println("\n=== Testing Current Timetable for JFK ===")
	currentSchedules, err := client.GetAirportSchedule("JFK")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Found %d current schedules\n", len(currentSchedules))
		for i, schedule := range currentSchedules {
			if i >= 5 {
				fmt.Printf("... and %d more schedules\n", len(currentSchedules)-5)
				break
			}
			fmt.Printf("  %s: %s -> %s (Type: %s, Status: %s)\n",
				schedule.Flight.IataNumber,
				schedule.Departure.IataCode,
				schedule.Arrival.IataCode,
				schedule.Type,
				schedule.Status)
		}
	}
}
