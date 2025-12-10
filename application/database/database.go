package database

import (
	"darbelis.eu/persedimai/env"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

// sql connection wrapper

type Database struct {
	connection *sql.DB
	dbConfig   *DBConfig
}

// TODO move to di package
var singletonDatabase *Database = nil

func GetDatabase(environment string) *Database {

	// TODO depending on environment load different env file
	if singletonDatabase == nil {
		singletonDatabase = &Database{}
		singletonDatabase.dbConfig = &DBConfig{}
		singletonDatabase.dbConfig.LoadFromEnv("../.env")
		singletonDatabase.connection = nil
	}

	return singletonDatabase
}

// TODO move to di package till here

func (db *Database) connect() error {
	cfg := mysql.NewConfig()

	var err error

	envFile, err := env.EnvMap{}.Read(".env")

	if err != nil {
		panic(err)
	}

	// vietoj 'os', padaryta 'envFile'

	cfg.User = envFile.Getenv("DBUSER")
	cfg.Passwd = envFile.Getenv("DBPASS")
	host := envFile.Getenv("DBHOST")
	port := envFile.Getenv("DBPORT")
	cfg.DBName = envFile.Getenv("DBNAME")

	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port

	// Get a database handle.

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

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func (database *Database) GetDatabaseName() string {
	return database.dbConfig.Dbname()
}
