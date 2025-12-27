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

	// Test future schedules for VNO
	fmt.Println("=== Testing Future Schedules for VNO (Vilnius) ===")

	dates := []string{"2025-12-28", "2025-12-30", "2026-01-05", "2026-01-10"}

	for _, date := range dates {
		fmt.Printf("\nDate: %s\n", date)

		// Departures
		depSchedules, err := client.GetFutureSchedules(aviation_edge.FutureSchedulesParams{
			IataCode: "VNO",
			Type:     "departure",
			Date:     date,
		})
		if err != nil {
			fmt.Printf("  Departures: Error - %v\n", err)
		} else {
			fmt.Printf("  Departures: %d flights\n", len(depSchedules))
			if len(depSchedules) > 0 {
				fmt.Printf("    Sample: %s to %s at %s\n",
					depSchedules[0].Flight.IataNumber,
					depSchedules[0].Arrival.IataCode,
					depSchedules[0].Departure.ScheduledTime)
			}
		}

		// Arrivals
		arrSchedules, err := client.GetFutureSchedules(aviation_edge.FutureSchedulesParams{
			IataCode: "VNO",
			Type:     "arrival",
			Date:     date,
		})
		if err != nil {
			fmt.Printf("  Arrivals: Error - %v\n", err)
		} else {
			fmt.Printf("  Arrivals: %d flights\n", len(arrSchedules))
			if len(arrSchedules) > 0 {
				fmt.Printf("    Sample: %s from %s at %s\n",
					arrSchedules[0].Flight.IataNumber,
					arrSchedules[0].Departure.IataCode,
					arrSchedules[0].Arrival.ScheduledTime)
			}
		}
	}
}
