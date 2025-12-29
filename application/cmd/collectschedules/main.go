package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

// ImportPlan contains the plan for importing airport data
type ImportPlan struct {
	SkipImport       bool
	ImportStartDate  string
	ImportEndDate    string
	MetaStartDate    time.Time // Date to save to metadata
	MetaEndDate      time.Time // Date to save to metadata
	AlreadyCoveredBy string    // Description of existing coverage
}

// calculateImportPlan determines what date range needs to be imported based on existing metadata
func calculateImportPlan(startDate, endDate string, meta *tables.AirportMeta) *ImportPlan {
	requestedStart := util.ParseDate(startDate)
	requestedEnd := util.ParseDate(endDate)

	plan := &ImportPlan{}

	// Check if data is already imported
	if meta != nil && meta.ImportedFrom != nil && meta.ImportedTo != nil {
		// Check if requested range is fully covered
		if !meta.ImportedFrom.After(requestedStart) && !meta.ImportedTo.Before(requestedEnd) {
			plan.SkipImport = true
			plan.AlreadyCoveredBy = fmt.Sprintf("%s to %s",
				meta.ImportedFrom.Format(time.DateOnly),
				meta.ImportedTo.Format(time.DateOnly))
			return plan
		}

		// Partially covered - import the whole range (including any gaps)
		// Calculate the full range to import
		plan.ImportStartDate = startDate
		plan.ImportEndDate = endDate

		// Metadata should cover the entire range (merge with existing)
		if meta.ImportedFrom.Before(requestedStart) {
			plan.MetaStartDate = *meta.ImportedFrom
		} else {
			plan.MetaStartDate = requestedStart
		}

		if meta.ImportedTo.After(requestedEnd) {
			plan.MetaEndDate = *meta.ImportedTo
		} else {
			plan.MetaEndDate = requestedEnd
		}
	} else {
		// No existing metadata or no import dates - import the full range
		plan.ImportStartDate = startDate
		plan.ImportEndDate = endDate
		plan.MetaStartDate = requestedStart
		plan.MetaEndDate = requestedEnd
	}

	return plan
}

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

	collector := di.Wrap(di.GetDataCollector)

	airportsMetaDao := di.Wrap(di.GetAirportsMetaDao)
	err = airportsMetaDao.CreateTable()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = collector.InitializeEuropeanAirportsMeta()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get existing import metadata for the airport
	meta, err := airportsMetaDao.Get(airportCode)
	if err != nil {
		log.Fatalf("Failed to get airport metadata: %v", err)
	}

	// Calculate what needs to be imported
	plan := calculateImportPlan(startDate, endDate, meta)

	// Check if import can be skipped
	if plan.SkipImport {
		fmt.Printf("Airport %s is already imported in the range %s (requested: %s to %s)\n",
			airportCode,
			plan.AlreadyCoveredBy,
			startDate, endDate)
		fmt.Println("No import needed.")
		return
	}

	// Execute import for the date range
	fmt.Printf("Collecting departure schedules for airport %s from %s to %s\n",
		airportCode, plan.ImportStartDate, plan.ImportEndDate)
	err = collector.CollectDepartureSchedules(airportCode, plan.ImportStartDate, plan.ImportEndDate)
	if err != nil {
		log.Fatalf("Failed to collect schedules: %v", err)
	}

	// Update metadata with import information
	if meta == nil {
		meta = &tables.AirportMeta{
			AirportCode: airportCode,
		}
	}
	meta.ImportedFrom = &plan.MetaStartDate
	meta.ImportedTo = &plan.MetaEndDate

	err = airportsMetaDao.Upsert(meta, true)
	if err != nil {
		log.Fatalf("Failed to update airport metadata: %v", err)
	}

	fmt.Printf("\nImport completed! Airport %s metadata updated with range %s to %s\n",
		airportCode,
		plan.MetaStartDate.Format(time.DateOnly),
		plan.MetaEndDate.Format(time.DateOnly))

	//fmt.Printf("\nCollection completed! Total schedules collected: %d\n", consumer.TotalCount)
}
