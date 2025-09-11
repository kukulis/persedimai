package main

import (
	"darbelis.eu/persedimai/database"
	"fmt"
)

func main() {

	// just testing only
	testDatabase()
	testDatabase()
}

func testDatabase() {
	db := database.GetDatabase()
	rows, err := db.GetConnection().Query("SELECT * FROM albums")

	if err != nil {
		fmt.Errorf("Error: %v", err)
		return
	}
	defer rows.Close()

	var albums []database.Album

	for rows.Next() {
		var alb database.Album
		err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		if err != nil {
			fmt.Errorf("albumsByArtist error: %v", err)
		}
		albums = append(albums, alb)
	}

	fmt.Printf("Albums found: %v\n", albums)

}
