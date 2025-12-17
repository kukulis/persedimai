package main

import (
	"darbelis.eu/persedimai/di"
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
	default:
		fmt.Printf("Error: unknown strategy '%s'\n", strategy)
		fmt.Println("Available strategies: simple, complex")
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
