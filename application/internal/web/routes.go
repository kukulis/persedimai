package web

import (
	"darbelis.eu/persedimai/di"
	"github.com/gin-gonic/gin"
	"html/template"
)

// GetRouter TODO pass DI factory through parameters
func GetRouter() *gin.Engine {

	// TODO get from DI factory later

	di.InitializeSingletons("prod")

	//db := &database.Database{}
	flightsSearchController := &FlightsSearchController{database: di.DatabaseInstance}
	travelSearchController := &TravelSearchController{}
	//pointsController := api.NewPointsController(db)

	router := gin.Default()

	// Add custom template functions
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	router.LoadHTMLGlob("templates/*")

	router.GET("/", HomeController)

	flightsGroup := router.Group("/flights")
	flightsGroup.GET("/search", func(c *gin.Context) { flightsSearchController.SearchForm(c) })
	flightsGroup.POST("/search", func(c *gin.Context) { flightsSearchController.SearchResult(c) })

	travelGroup := router.Group("/travel")
	travelGroup.GET("/search", func(c *gin.Context) { travelSearchController.SearchForm(c) })
	travelGroup.POST("/search", func(c *gin.Context) { travelSearchController.SearchResult(c) })

	apiGroup := router.Group("/api")
	apiGroup.GET("/points", func(c *gin.Context) { di.ApiPointsControllerInstance.GetAll(c) })
	apiGroup.GET("/points/bounds", func(c *gin.Context) { di.ApiPointsControllerInstance.GetBounds(c) })
	apiGroup.GET("/travels/bounds", func(c *gin.Context) { di.ApiTravelsControllerInstance.GetTimeBounds(c) })

	return router
}
