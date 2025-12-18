package main

import (
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"darbelis.eu/persedimai/internal/web"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	// -- pseudo test
	f := tables.Transfer{
		ID: "0", From: "1", To: "2", Departure: util.ParseDate("2025-01-01"), Arrival: util.ParseDateTime("2025-02-02 11:30:00"),
	}

	fmt.Printf("Transfer : %v\n", f)
	fmt.Printf("Ok : %v\n", http.StatusOK)
	// --

	router := web.GetRouter()

	port := strings.ToLower(os.Getenv("PORT"))
	if port == "" {
		port = "7080"
	}

	err = router.Run(":" + port)

	fmt.Println(err)
}
