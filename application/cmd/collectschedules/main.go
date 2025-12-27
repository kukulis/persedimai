package main

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
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

	// Get API key from environment
	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: AVIATION_EDGE_API_KEY is not set")
		fmt.Println("Set it in .env file or export AVIATION_EDGE_API_KEY=your_key")
		os.Exit(1)
	}

	// Create API client
	apiClient := aviation_edge.NewAviationEdgeApiClient(apiKey)

	// Create DataCollector
	collector := aviation_edge.NewDataCollector(apiClient)

	// Create PrintScheduleConsumer
	consumer := &aviation_edge.PrintScheduleConsumer{}

	// Collect schedules
	fmt.Printf("Collecting departure schedules for airport %s from %s to %s\n", airportCode, startDate, endDate)
	err = collector.CollectDepartureSchedules(airportCode, startDate, endDate, consumer)
	if err != nil {
		log.Fatalf("Failed to collect schedules: %v", err)
	}

	fmt.Printf("\nCollection completed! Total schedules collected: %d\n", consumer.TotalCount)
}
