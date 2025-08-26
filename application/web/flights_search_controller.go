package web

import "github.com/gin-gonic/gin"

func FlightsSearchFormController(c *gin.Context) {
	c.String(200, "Flights Search Form")
}

func FlightsSearchHandleController(c *gin.Context) {
	c.String(200, "Flights Search Handle")
}
