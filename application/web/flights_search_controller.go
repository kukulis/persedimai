package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FlightsSearchFormController(c *gin.Context) {
	c.HTML(http.StatusOK, "search-flights-form.html", gin.H{})
}

func FlightsSearchHandleController(c *gin.Context) {
	c.String(200, "Flights Search Handle")
}
