package api

import (
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PointsController struct {
	database *database.Database
}

func NewPointsController(db *database.Database) *PointsController {
	return &PointsController{database: db}
}

// GetAll returns all points as JSON
func (controller *PointsController) GetAll(c *gin.Context) {
	pointDao := dao.NewPointDao(controller.database)

	points, err := pointDao.SelectAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"points": points,
		"count":  len(points),
	})
}
