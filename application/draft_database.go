package main

import (
	"darbelis.eu/persedimai/env"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	fmt.Println("TODO draft test for database")

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

	fmt.Printf("database dsn: %s\n", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

}

//type EnvMap map[string]string
//
//func (envFile EnvMap) Getenv(key string) string {
//
//	// TODO padaryti, kad jeigu tuščia, tai kreiptis į os
//	return envFile[key]
//}
//
//func (envFile EnvMap) Read(filename string) (envMap EnvMap, err error) {
//	return godotenv.Read(filename)
//}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM albums WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
