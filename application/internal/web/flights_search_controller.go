package web

// @deprecated

import (
	"darbelis.eu/persedimai/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FlightsSearchController struct {
	database *database.Database
	// more DI later
}

func (controller *FlightsSearchController) SearchForm(c *gin.Context) {
	c.HTML(http.StatusOK, "search-flights-form.html", gin.H{})
}

func (controller *FlightsSearchController) SearchResult(c *gin.Context) {

	// 1) Get search parameters
	// 2) Get result from the service by the given parameters
	// 3) Pass result to the themplate
	c.HTML(http.StatusOK, "search-flights-result.html", gin.H{})
}
