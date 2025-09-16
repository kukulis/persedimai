package web

import (
	"darbelis.eu/persedimai/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

//// TODO structure with functions instead
//func FlightsSearchFormController(c *gin.Context) {
//	c.HTML(http.StatusOK, "search-flights-form.html", gin.H{})
//}
//
//// TODO structure with functions instead
//func FlightsSearchHandleController(c *gin.Context) {
//	c.String(200, "Flights Search Handle")
//}

type FlightsSearchController struct {
	database *database.Database
	// more DI later
}

func (controller *FlightsSearchController) SearchForm(c *gin.Context) {
	c.HTML(http.StatusOK, "search-flights-form.html", gin.H{})
}

func (controller *FlightsSearchController) SearchResult(c *gin.Context) {
	c.String(200, "Flights Search Handle")
}
