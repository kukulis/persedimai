package main

import (
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/util"
	"darbelis.eu/persedimai/web"
	"fmt"
	"net/http"
)

func main() {
	// -- pseudo test
	f := tables.Travel{
		ID: 0, From: 1, To: 2, Departure: util.ParseDate("2025-01-01"), Arrival: util.ParseDateTime("2025-02-02 11:30:00"),
	}

	fmt.Printf("Travel : %v\n", f)
	fmt.Printf("Ok : %v\n", http.StatusOK)
	// --

	router := web.GetRouter()
	router.LoadHTMLGlob("templates/*")
	err := router.Run(":8080")

	fmt.Println(err)
}
