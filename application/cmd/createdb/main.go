package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/database"
	"flag"
	"fmt"
	"log"
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

	// Get DBConfig to retrieve user credentials for GRANT PRIVILEGES
	fmt.Println("Retrieving database configuration for granting privileges...")
	dbConfig, err := di.NewDbConfig(dbName)
	if err != nil {
		fmt.Printf("Warning: Failed to retrieve database configuration: %v\n", err)
		fmt.Println("Database created, but privileges were not granted.")
		fmt.Println("You may need to grant privileges manually.")
		os.Exit(0)
	}

	fmt.Printf("Granting privileges to user '%s'...\n", dbConfig.Username)
	err = grantPrivileges(db, dbName, dbConfig)
	if err != nil {
		fmt.Printf("Error granting privileges: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Privileges granted successfully to user '%s' on database '%s'!\n", dbConfig.Username, dbName)
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

func grantPrivileges(db *database.Database, dbName string, dbConfig *database.DBConfig) error {
	conn, err := db.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Sanitize parameters to prevent SQL injection
	escapedDbName := database.MysqlRealEscapeString(dbName)
	escapedUsername := database.MysqlRealEscapeString(dbConfig.Username)

	// GRANT ALL PRIVILEGES ON database.* TO 'user'@'host'
	query := fmt.Sprintf("GRANT ALL ON `%s`.* TO '%s'", escapedDbName, escapedUsername)
	log.Println("The grant query is " + query)

	_, err = conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute GRANT PRIVILEGES: %w", err)
	}

	// Execute FLUSH PRIVILEGES to ensure changes take effect
	_, err = conn.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return fmt.Errorf("failed to execute FLUSH PRIVILEGES: %w", err)
	}

	return nil
}
