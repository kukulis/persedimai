package main

import (
	"darbelis.eu/persedimai/tables"
	"fmt"
)

func main() {
	fmt.Printf("TODO persedimai\n")

	f := tables.Flight{
		1, 2, "2025-01-01", "2025-02-02",
	}

	fmt.Printf("Flight : %v\n", f)
}
