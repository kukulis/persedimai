package web

import (
	"darbelis.eu/persedimai/database"
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
	c.HTML(http.StatusOK, "search-flights-result.html", gin.H{})
}
