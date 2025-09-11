package database

import (
	"darbelis.eu/persedimai/env"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

// sql connection wrapper

type Database struct {
	connection *sql.DB
	dbConfig   *DBConfig
}

var singletonDatabase *Database

func GetDatabase() *Database {
	if singletonDatabase == nil {
		singletonDatabase = &Database{}
		singletonDatabase.dbConfig = &DBConfig{}
		singletonDatabase.dbConfig.LoadFromEnv("../.env")
		singletonDatabase.connection = nil
	}

	return singletonDatabase
}

func (db *Database) connect() error {
	cfg := mysql.NewConfig()

	envFile, _ := env.EnvMap{}.Read(".env")

	// vietoj 'os', padaryta 'envFile'

	cfg.User = envFile.Getenv("DBUSER")
	cfg.Passwd = envFile.Getenv("DBPASS")
	host := envFile.Getenv("DBHOST")
	port := envFile.Getenv("DBPORT")
	cfg.DBName = envFile.Getenv("DBNAME")

	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port

	// Get a database handle.
	var err error
	dsn := cfg.FormatDSN()

	// is this correct?
	log.Printf("database dsn: %s\n", dsn)

	db.connection, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.connection.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Connected!")

	return err
}

func (db *Database) GetConnection() *sql.DB {

	if db.connection == nil {
		log.Println("GetConnection: Connecting to database...")
		err := db.connect()
		if err != nil {
			log.Fatalf("Error connecting to database: %s", err)
		}
	} else {
		log.Println("GetConnection: already connected!")
	}

	return db.connection
}
