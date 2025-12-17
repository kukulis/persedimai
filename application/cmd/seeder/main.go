package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"darbelis.eu/persedimai/internal/database"
	"flag"
	"fmt"
	"os"
)

func main() {
	var environment string
	var strategy string

	flag.StringVar(&environment, "env", "dev", "Database environment (dev, test, prod)")
	flag.StringVar(&strategy, "strategy", "", "Data generation strategy")
	flag.Parse()

	if strategy == "" {
		fmt.Println("Error: strategy parameter is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Connecting to database environment: %s\n", environment)
	db, err := di.NewDatabase(environment)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Executing strategy: %s\n", strategy)
	executeStrategy(db, strategy)

	fmt.Println("Data generation completed successfully!")
}

func executeStrategy(db *database.Database, strategy string) {
	switch strategy {
	case "simple":
		generateSimpleData(db)
	case "complex":
		generateComplexData(db)
	case "normal":
		generateNormalData(db)
	case "complete":
		generateCompleteData(db)
	default:
		fmt.Printf("Error: unknown strategy '%s'\n", strategy)
		fmt.Println("Available strategies: simple, complex, normal, complete")
		os.Exit(1)
	}
}

func generateSimpleData(db *database.Database) {
	fmt.Println("Generating simple data...")
	// TODO: Implement simple data generation
}

func generateComplexData(db *database.Database) {
	fmt.Println("Generating complex data...")
	// TODO: Implement complex data generation
}

func generateNormalData(db *database.Database) {
	fmt.Println("Generating normal data...")

	dbFiller := integration_tests.DatabaseFiller{}

	fmt.Println("Filling test database...")
	err := dbFiller.FillDatabase(db)
	if err != nil {
		fmt.Printf("Error filling database: %v\n", err)
		os.Exit(1)
	}

	err = dbFiller.LogResults()
	if err != nil {
		fmt.Printf("Error logging results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Normal data generation finished!")
}

func generateCompleteData(db *database.Database) {
	fmt.Println("Generating complete data...")

	dbFiller := integration_tests.DatabaseFiller{}

	fmt.Println("Filling test database...")
	err := dbFiller.FillDatabase(db)
	if err != nil {
		fmt.Printf("Error filling database: %v\n", err)
		os.Exit(1)
	}

	err =
		()
	if err != nil {
		fmt.Printf("Error filling hubs travels: %v\n", err)
		os.Exit(1)
	}

	err = dbFiller.LogResults()
	if err != nil {
		fmt.Printf("Error logging results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Complete data generation finished!")
}
