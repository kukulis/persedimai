package main

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/aviation_edge"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	var environment string
	flag.StringVar(&environment, "env", "dev", "Database environment (dev, test, prod)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	di.InitializeSingletons(environment)

	client := di.Wrap(di.GetAviationEdgeClient)

	airports, err := client.GetAirports(
		aviation_edge.AirportsParams{})

	if err != nil {
		fmt.Println("Failed to get airports: ", err)
		return
	}

	//for _, airport := range airports {
	//	fmt.Printf("The extracted airport %v \n", airport)
	//}
	fmt.Printf("Airports amount= %d\n", len(airports))

	airportsDao := di.GetAirportsDao()

	stepSize := 100
	for step := 0; step <= len(airports); step += stepSize {
		airportPtrs := []*aviation_edge.AirportResponse{}

		for i := step; i < min(len(airports), step+stepSize); i++ {
			if airports[i].CodeIso2Country == "" {
				airports[i].CodeIso2Country = "-"
			}
			if airports[i].NameCountry == "" {
				airports[i].NameCountry = "-"
			}
			airportPtrs = append(airportPtrs, &airports[i])
		}

		err = airportsDao.Upsert(airportPtrs)
		if err != nil {
			fmt.Println("Failed to update airports: ", err)
			return
		}

	}

	fmt.Printf("Airports storing finished \n")

}
