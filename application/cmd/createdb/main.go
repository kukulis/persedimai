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
	var dbName string

	flag.StringVar(&environment, "env", "dev", "Database environment (dev, test, prod)")
	flag.StringVar(&dbName, "name", "", "Database name to create")
	flag.Parse()

	if dbName == "" {
		fmt.Println("Error: database name is required")
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

	fmt.Printf("Creating database: %s\n", dbName)
	err = createDatabase(db, dbName)
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database '%s' created successfully!\n", dbName)
}

func createDatabase(db *database.Database, dbName string) error {
	conn, err := db.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Sanitize database name to prevent SQL injection
	escapedDbName := database.MysqlRealEscapeString(dbName)

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", escapedDbName)
	_, err = conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute CREATE DATABASE: %w", err)
	}

	return nil
}
