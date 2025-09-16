package web

import (
	"darbelis.eu/persedimai/database"
	"github.com/gin-gonic/gin"
)

// GetRouter TODO pass DI factory through parameters
func GetRouter() *gin.Engine {

	// TODO get from DI factory later
	db := &database.Database{}
	flightsSearchController := &FlightsSearchController{database: db}

	router := gin.Default()

	router.GET("/", HomeController)

	flightsGroup := router.Group("/flights")

	flightsGroup.GET("/search", func(c *gin.Context) { flightsSearchController.SearchForm(c) })
	flightsGroup.POST("/search", func(c *gin.Context) { flightsSearchController.SearchResult(c) })

	return router
}
