package main

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	apiKey := os.Getenv("AVIATION_EDGE_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: AVIATION_EDGE_API_KEY is not set - API calls will fail")
		fmt.Println("Set it in .env file or export AVIATION_EDGE_API_KEY=your_key")
	}
	aviation_edge.ExampleUsage(apiKey)
}
