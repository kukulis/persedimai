package main

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")

	client := aviation_edge.NewAviationEdgeApiClient(apiKey)

	fmt.Println("=== Analyzing JFK Current Timetable Time Range ===\n")

	// Get current timetable
	schedules, err := client.GetAirportSchedule("JFK")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total schedules: %d\n\n", len(schedules))

	// Parse all times and find min/max
	var minTime, maxTime time.Time
	layout := "2006-01-02T15:04:05.000"

	for _, schedule := range schedules {
		var schedTime time.Time
		var err error

		// Use departure time for departure type, arrival time for arrival type
		if schedule.Type == "departure" && schedule.Departure.ScheduledTime != "" {
			schedTime, err = time.Parse(layout, schedule.Departure.ScheduledTime)
		} else if schedule.Type == "arrival" && schedule.Arrival.ScheduledTime != "" {
			schedTime, err = time.Parse(layout, schedule.Arrival.ScheduledTime)
		}

		if err != nil {
			continue
		}

		if minTime.IsZero() || schedTime.Before(minTime) {
			minTime = schedTime
		}
		if maxTime.IsZero() || schedTime.After(maxTime) {
			maxTime = schedTime
		}
	}

	now := time.Now()

	fmt.Printf("Current time (server): %s\n", now.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Earliest schedule:     %s\n", minTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Latest schedule:       %s\n", maxTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("\nTime range coverage:   %s\n", maxTime.Sub(minTime).Round(time.Hour))
	fmt.Printf("Hours from now to earliest: %.1f hours\n", minTime.Sub(now).Hours())
	fmt.Printf("Hours from now to latest:   %.1f hours\n", maxTime.Sub(now).Hours())

	// Count by status
	statusCount := make(map[string]int)
	for _, schedule := range schedules {
		statusCount[schedule.Status]++
	}

	fmt.Printf("\nSchedules by status:\n")
	for status, count := range statusCount {
		fmt.Printf("  %s: %d\n", status, count)
	}

	// Sample some schedules from different times
	fmt.Printf("\n=== Sample Schedules ===\n")
	samples := []struct {
		desc string
		time time.Time
	}{
		{"Earliest", minTime},
		{"Mid-range", minTime.Add(maxTime.Sub(minTime) / 2)},
		{"Latest", maxTime},
	}

	for _, sample := range samples {
		fmt.Printf("\n%s (around %s):\n", sample.desc, sample.time.Format("2006-01-02 15:04"))
		count := 0
		for _, schedule := range schedules {
			var schedTime time.Time
			if schedule.Type == "departure" && schedule.Departure.ScheduledTime != "" {
				schedTime, _ = time.Parse(layout, schedule.Departure.ScheduledTime)
			} else if schedule.Type == "arrival" && schedule.Arrival.ScheduledTime != "" {
				schedTime, _ = time.Parse(layout, schedule.Arrival.ScheduledTime)
			}

			if schedTime.IsZero() {
				continue
			}

			diff := schedTime.Sub(sample.time)
			if diff < 0 {
				diff = -diff
			}
			if diff < 30*time.Minute {
				fmt.Printf("  %s %s: %s -> %s (%s)\n",
					schedule.Flight.IataNumber,
					schedule.Type,
					schedule.Departure.IataCode,
					schedule.Arrival.IataCode,
					schedule.Status)
				count++
				if count >= 3 {
					break
				}
			}
		}
	}
}
