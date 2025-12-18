package api

import (
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TravelsController struct {
	databasesContainer database.DatabasesContainer
}

func NewTravelsController(db database.DatabasesContainer) *TravelsController {
	return &TravelsController{databasesContainer: db}
}

// GetTimeBounds returns the min/max time boundaries for all travels
func (controller *TravelsController) GetTimeBounds(c *gin.Context) {
	env := "test"
	if databaseParam := c.Query("database"); databaseParam != "" {
		env = databaseParam
	}

	db, err := controller.databasesContainer.GetDatabase(env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database: " + err.Error(),
		})
		return
	}

	travelDao := dao.NewTravelDao(db)

	bounds, err := travelDao.GetTimeBounds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bounds": gin.H{
			"minDeparture": bounds.MinDeparture.Format(time.DateTime),
			"maxDeparture": bounds.MaxDeparture.Format(time.DateTime),
			"minArrival":   bounds.MinArrival.Format(time.DateTime),
			"maxArrival":   bounds.MaxArrival.Format(time.DateTime),
		},
		"database": env,
	})
}
