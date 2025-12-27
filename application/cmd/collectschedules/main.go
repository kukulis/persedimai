package main

import (
	"darbelis.eu/persedimai/di"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	var airportCode string
	var startDate string
	var endDate string
	var environment string

	flag.StringVar(&environment, "env", "dev", "Database environment (dev, test, prod)")
	flag.StringVar(&airportCode, "airport", "", "Airport IATA code (e.g., VNO, JFK)")
	flag.StringVar(&startDate, "start", "", "Start date in YYYY-MM-DD format")
	flag.StringVar(&endDate, "end", "", "End date in YYYY-MM-DD format")
	flag.Parse()

	// Validate required parameters
	if airportCode == "" {
		fmt.Println("Error: airport parameter is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  collectschedules -airport VNO -start 2025-12-27 -end 2025-12-30")
		os.Exit(1)
	}

	if startDate == "" {
		fmt.Println("Error: start date parameter is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if endDate == "" {
		fmt.Println("Error: end date parameter is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	di.InitializeSingletons(environment)

	collector := di.DataCollectorLoader()

	// Collect schedules
	fmt.Printf("Collecting departure schedules for airport %s from %s to %s\n", airportCode, startDate, endDate)
	err = collector.CollectDepartureSchedules(airportCode, startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to collect schedules: %v", err)
	}

	//fmt.Printf("\nCollection completed! Total schedules collected: %d\n", consumer.TotalCount)
}
