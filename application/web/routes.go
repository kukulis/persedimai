package web

import "github.com/gin-gonic/gin"

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", HomeController)

	flightsGroup := router.Group("/flights")

	flightsGroup.GET("/search", FlightsSearchFormController)
	flightsGroup.POST("/search", FlightsSearchHandleController)

	return router
}
