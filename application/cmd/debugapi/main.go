package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")

	// Test historical endpoint
	url := fmt.Sprintf("https://aviation-edge.com/v2/public/flightsHistory?key=%s&code=VNO&type=departure&date_from=2025-12-20&date_to=2025-12-20", apiKey)

	fmt.Println("=== Testing Historical Schedules API ===")
	fmt.Printf("URL: %s\n\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Content-Type: %s\n\n", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("=== Raw Response (first 2000 chars) ===")
	if len(body) > 2000 {
		fmt.Println(string(body[:2000]))
		fmt.Printf("\n... and %d more bytes\n", len(body)-2000)
	} else {
		fmt.Println(string(body))
	}
}
