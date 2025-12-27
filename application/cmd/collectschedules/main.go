package main

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// Command-line flags
	var countryCode string
	var airportCode string
	var startDate string
	var endDate string
	var outputFile string
	var includeDepartures bool
	var includeArrivals bool
	var rateLimitSeconds int
	var printOnly bool
	var useCurrent bool

	flag.StringVar(&countryCode, "country", "", "Country code (e.g., US, GB, FR)")
	flag.StringVar(&airportCode, "airport", "", "Single airport IATA code (e.g., JFK, LAX) - uses current timetable")
	flag.StringVar(&startDate, "start", "", "Start date in format YYYY-MM-DD (default: today)")
	flag.StringVar(&endDate, "end", "", "End date in format YYYY-MM-DD (default: today)")
	flag.StringVar(&outputFile, "output", "schedules.json", "Output JSON file path")
	flag.BoolVar(&includeDepartures, "departures", true, "Include departure schedules")
	flag.BoolVar(&includeArrivals, "arrivals", true, "Include arrival schedules")
	flag.IntVar(&rateLimitSeconds, "rate-limit", 1, "Delay between API calls in seconds")
	flag.BoolVar(&printOnly, "print", false, "Print to stdout instead of saving to file")
	flag.BoolVar(&useCurrent, "current", false, "Use current timetable instead of historical/future")

	flag.Parse()

	// Validate required parameters
	if countryCode == "" && airportCode == "" {
		fmt.Println("Error: either -country or -airport is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  collectschedules -country US -start 2025-12-20 -end 2025-12-22")
		fmt.Println("  collectschedules -country GB -start 2025-12-27 -output uk_schedules.json")
		fmt.Println("  collectschedules -airport JFK -print")
		fmt.Println("  collectschedules -airport LAX -output lax.json -departures=false")
		os.Exit(1)
	}

	if airportCode != "" && countryCode != "" {
		fmt.Println("Error: cannot specify both -country and -airport")
		os.Exit(1)
	}

	// Set default dates if not provided
	if startDate == "" {
		startDate = time.Now().Format(time.DateOnly)
	}
	if endDate == "" {
		endDate = startDate
	}

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, looking for AVIATION_EDGE_API_KEY environment variable")
	}

	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: AVIATION_EDGE_API_KEY environment variable not set")
		fmt.Println("Please set it in your .env file or export it:")
		fmt.Println("  export AVIATION_EDGE_API_KEY=your-api-key")
		os.Exit(1)
	}

	// Display configuration
	fmt.Println("=== Aviation Edge Schedule Collection ===")
	if airportCode != "" {
		fmt.Printf("Airport Code: %s (Current Timetable)\n", airportCode)
	} else {
		fmt.Printf("Country Code: %s\n", countryCode)
		fmt.Printf("Date Range: %s to %s\n", startDate, endDate)
		fmt.Printf("Rate Limit: %d second(s) between calls\n", rateLimitSeconds)
	}
	fmt.Printf("Include Departures: %t\n", includeDepartures)
	fmt.Printf("Include Arrivals: %t\n", includeArrivals)
	if printOnly {
		fmt.Println("Output Mode: Print to stdout")
	} else {
		fmt.Printf("Output File: %s\n", outputFile)
	}
	fmt.Println("========================================")
	fmt.Println()

	// Initialize API client
	log.Println("Initializing Aviation Edge API client...")
	apiClient := aviation_edge.NewAviationEdgeApiClient(apiKey)

	// Create data collector with dependency injection
	collector := aviation_edge.NewDataCollector(apiClient)

	// Create consumer based on mode
	var consumer aviation_edge.ScheduleConsumer
	if printOnly {
		consumer = &aviation_edge.PrintScheduleConsumer{}
	} else {
		consumer = aviation_edge.NewFileScheduleConsumer(outputFile, true)
	}

	// Collect schedules
	log.Println("Starting schedule collection...")
	startTime := time.Now()

	if airportCode != "" {
		// Single airport mode - use current timetable
		err = collectSingleAirport(apiClient, airportCode, consumer, includeDepartures, includeArrivals)
		if err != nil {
			log.Fatalf("Error collecting schedules: %v", err)
		}
	} else {
		// Country mode - use data collector
		err = collector.CollectSchedules(aviation_edge.CollectSchedulesParams{
			CountryCode:       countryCode,
			StartDate:         startDate,
			EndDate:           endDate,
			IncludeDepartures: includeDepartures,
			IncludeArrivals:   includeArrivals,
			Consumer:          consumer,
			RateLimitDelay:    time.Duration(rateLimitSeconds) * time.Second,
		})

		if err != nil {
			log.Fatalf("Error collecting schedules: %v", err)
		}
	}

	// Flush file consumer if used
	if !printOnly {
		if fileConsumer, ok := consumer.(*aviation_edge.FileScheduleConsumer); ok {
			if err := fileConsumer.Flush(); err != nil {
				log.Fatalf("Error flushing schedules to file: %v", err)
			}
			log.Printf("Successfully saved %d schedules to %s", len(fileConsumer.Schedules), outputFile)
		}
	} else {
		if printConsumer, ok := consumer.(*aviation_edge.PrintScheduleConsumer); ok {
			log.Printf("Total schedules printed: %d", printConsumer.TotalCount)
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("Collection completed successfully!\n")
	fmt.Printf("Total time: %s\n", elapsed.Round(time.Second))
	fmt.Println("========================================")
}

func collectSingleAirport(apiClient *aviation_edge.AviationEdgeApiClient, airportCode string, consumer aviation_edge.ScheduleConsumer, includeDepartures, includeArrivals bool) error {
	log.Printf("Collecting current schedules for airport: %s", airportCode)

	var allSchedules []aviation_edge.ScheduleResponse

	// Collect departure schedules
	if includeDepartures {
		log.Printf("Fetching departure schedules...")
		schedules, err := apiClient.GetFlightSchedules(aviation_edge.FlightSchedulesParams{
			IataCode: airportCode,
			Type:     "departure",
		})
		if err != nil {
			return fmt.Errorf("failed to get departure schedules: %w", err)
		}
		log.Printf("Found %d departure schedules", len(schedules))
		allSchedules = append(allSchedules, schedules...)
	}

	// Collect arrival schedules
	if includeArrivals {
		log.Printf("Fetching arrival schedules...")
		schedules, err := apiClient.GetFlightSchedules(aviation_edge.FlightSchedulesParams{
			IataCode: airportCode,
			Type:     "arrival",
		})
		if err != nil {
			return fmt.Errorf("failed to get arrival schedules: %w", err)
		}
		log.Printf("Found %d arrival schedules", len(schedules))
		allSchedules = append(allSchedules, schedules...)
	}

	// Consume all collected schedules
	if len(allSchedules) > 0 {
		if err := consumer.Consume(allSchedules); err != nil {
			return fmt.Errorf("consumer failed: %w", err)
		}
		log.Printf("Total schedules collected: %d", len(allSchedules))
	} else {
		log.Println("No schedules found")
	}

	return nil
}
