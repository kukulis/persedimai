package main

import (
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/web"
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("TODO persedimai\n")

	f := tables.Flight{
		1, 2, "2025-01-01", "2025-02-02",
	}

	fmt.Printf("Flight : %v\n", f)

	//gin.BasicAuth(gin.Accounts{
	//	"foo": "bar",
	//})

	fmt.Printf("Ok : %v\n", http.StatusOK)

	// TODO

	web.GetRouter().Run(":8080")

}
