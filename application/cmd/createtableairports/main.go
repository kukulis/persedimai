package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/migrations"
	"flag"
	"fmt"
	"os"
)

func main() {
	var environment string

	flag.StringVar(&environment, "env", "dev", "Database environment (dev, test, prod)")

	flag.Parse()

	fmt.Printf("Connecting to database environment: %s\n", environment)
	db, err := di.NewDatabase(environment)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating airports table...")
	err = migrations.CreateAirportsTable(db)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Table created successfully")

	fmt.Println("Done!")
}
