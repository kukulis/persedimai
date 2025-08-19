package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	fmt.Println("TODO draft test for database")

	cfg := mysql.NewConfig()

	envFile, _ := EnvMap{}.Read(".env")

	// vietoj 'os', padaryta 'envFile'

	cfg.User = envFile.Getenv("DBUSER")
	cfg.Passwd = envFile.Getenv("DBPASS")
	host := envFile.Getenv("DBHOST")
	port := envFile.Getenv("DBPORT")
	cfg.DBName = envFile.Getenv("DBNAME")

	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port

	// Get a database handle.
	//var err error
	dsn := cfg.FormatDSN()

	fmt.Printf("database dsn: %s\n", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

type EnvMap map[string]string

func (envFile EnvMap) Getenv(key string) string {

	// TODO padaryti, kad jeigu tuščia, tai kreiptis į os
	return envFile[key]
}

func (envFile EnvMap) Read(filename string) (envMap EnvMap, err error) {
	return godotenv.Read(filename)
}
