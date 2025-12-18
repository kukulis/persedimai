package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// sql connection wrapper

type Database struct {
	connection *sql.DB
	dbConfig   *DBConfig
	version    string
}

func NewDatabase(config *DBConfig) *Database {
	return &Database{
		dbConfig: config,
		version:  "",
	}
}

func (db *Database) connect() error {
	var err error
	dsn, err := db.dbConfig.FormatDsn()

	if err != nil {
		return err
	}

	// is this correct?
	log.Printf("database dsn: %s\n", dsn)

	db.connection, err = sql.Open(db.dbConfig.DbType, dsn)

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

func (db *Database) GetConnection() (*sql.DB, error) {
	if db.connection == nil {
		log.Println("GetConnection: Connecting to database...")
		err := db.connect()
		if err != nil {
			return nil, err
		}
	} else {
		//log.Println("GetConnection: already connected!")
	}

	return db.connection, nil
}

func (db *Database) CloseConnection() error {
	var err error = nil
	if db.connection != nil {
		err = db.connection.Close()
	}

	db.connection = nil

	if err != nil {
		fmt.Printf("Error closing database connection: %v\n", err)
	}

	return err
}

func MysqlRealEscapeString(value string) string {
	// Order matters: backslash must be escaped first to avoid double-escaping
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\x00", "\\0")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "'", "\\'")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\x1a", "\\Z")

	return value
}

func (db *Database) GetDatabaseName() string {
	return db.dbConfig.Dbname
}

func (db *Database) CheckVersion() string {
	if db.version != "" {
		return db.version
	}

	if db.dbConfig.DbType != "mysql" {
		return "CheckVersion for db type " + db.dbConfig.DbType + " not implemented."
	}

	dbConn, err := db.GetConnection()
	if err != nil {
		return err.Error()
	}

	err = dbConn.QueryRow("SELECT VERSION()").Scan(&db.version)
	if err != nil {
		return err.Error()
	}

	return db.version
}

func (db *Database) AddTimeoutToQuery(baseQuery string, timeout time.Duration) string {
	if timeout == 0 {
		return baseQuery
	}

	version := db.CheckVersion()

	if strings.Contains(strings.ToLower(version), "mariadb") {
		// MariaDB: use SET STATEMENT
		timeoutSec := int(timeout.Seconds())
		return fmt.Sprintf("SET STATEMENT max_statement_time=%d FOR %s",
			timeoutSec, baseQuery)
	} else {
		// MySQL: use optimizer hint
		timeoutMs := int(timeout.Milliseconds())
		return fmt.Sprintf("SELECT /*+ MAX_EXECUTION_TIME(%d) */ %s",
			timeoutMs,
			strings.TrimPrefix(baseQuery, "SELECT"))
	}
}
