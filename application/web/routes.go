package web

import "github.com/gin-gonic/gin"

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	flightsGroup := router.Group("/flights")

	flightsGroup.GET("/search", func(c *gin.Context) {
		c.String(200, "TODO flights form")
	})
	flightsGroup.POST("/search", func(c *gin.Context) {
		c.String(200, "TODO handle flights form")
	})

	return router
}
