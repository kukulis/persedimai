package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/database"
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

	fmt.Println("Creating clusters...")
	err = createClusters(db)
	if err != nil {
		fmt.Printf("Error creating clusters: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Clusters created successfully!")
}

func createClusters(db *database.Database) error {
	clustersCreator := migrations.NewClustersCreator(db)

	fmt.Println("Updating clusters on travels...")
	err := clustersCreator.UpdateClustersOnTravels()
	if err != nil {
		return fmt.Errorf("failed to update clusters on travels: %w", err)
	}

	fmt.Println("Creating cluster tables...")
	err = clustersCreator.CreateClustersTables()
	if err != nil {
		return fmt.Errorf("failed to create cluster tables: %w", err)
	}

	fmt.Println("Inserting cluster data...")
	err = clustersCreator.InsertClustersDatas()
	if err != nil {
		return fmt.Errorf("failed to insert cluster data: %w", err)
	}

	return nil
}
