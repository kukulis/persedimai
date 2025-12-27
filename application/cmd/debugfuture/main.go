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

	// Test near-term future date that fails
	url := fmt.Sprintf("https://aviation-edge.com/v2/public/flightsFuture?key=%s&iataCode=VNO&type=departure&date=2025-12-28", apiKey)

	fmt.Println("=== Testing Future Schedules API (Near-term date) ===")
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

	fmt.Println("=== Raw Response ===")
	fmt.Println(string(body))
}
