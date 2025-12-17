package api

import (
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PointsController struct {
	databasesContainer database.DatabasesContainer
}

func NewPointsController(db database.DatabasesContainer) *PointsController {
	return &PointsController{databasesContainer: db}
}

// GetAll returns all points as JSON
func (controller *PointsController) GetAll(c *gin.Context) {
	// Create filter and populate from query parameters
	filter := data.NewPointsFilter()

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// Parse X coordinate
	if xStr := c.Query("x"); xStr != "" {
		if x, err := strconv.ParseFloat(xStr, 64); err == nil {
			filter.X = &x
		}
	}

	// Parse Y coordinate
	if yStr := c.Query("y"); yStr != "" {
		if y, err := strconv.ParseFloat(yStr, 64); err == nil {
			filter.Y = &y
		}
	}

	// Parse name part
	if namePart := c.Query("name_part"); namePart != "" {
		filter.NamePart = namePart
	}

	// Parse ID part
	if idPart := c.Query("id_part"); idPart != "" {
		filter.IdPart = idPart
	}

	// TODO take env from a query path later or in some other way
	env := "test"

	db, err := controller.databasesContainer.GetDatabase(env)

	pointDao := dao.NewPointDao(db)

	points, err := pointDao.SelectWithFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"points": points,
		"count":  len(points),
		"filter": filter,
	})
}
