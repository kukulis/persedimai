package web

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/travel_finder"
	"darbelis.eu/persedimai/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TravelSearchController struct {
	// Dependencies will be injected later
}

type SearchFormData struct {
	Strategies []StrategyOption
	Databases  []DatabaseOption
	Points     []*tables.Point
}

type StrategyOption struct {
	Name  string
	Value string
}

type DatabaseOption struct {
	Name  string
	Value string
}

type SearchResultData struct {
	Strategy      string
	Database      string
	Source        string
	Destination   string
	ArrivalFrom   string
	ArrivalTo     string
	TravelCount   int
	Paths         []*TravelPath
	ExecutionTime string
	Error         string
}

type TravelPath struct {
	Transfers     []*TransferDisplay
	TotalDuration string
	TransferCount int
}

type TransferDisplay struct {
	From      string
	To        string
	Departure string
	Arrival   string
	Duration  string
}

func (controller *TravelSearchController) SearchForm(c *gin.Context) {
	// Get available databases from .env files
	databases := getAvailableDatabases()

	// Get strategies
	strategies := []StrategyOption{
		{Name: "Clustered Strategy", Value: "clustered"},
		{Name: "Simple Strategy", Value: "simple"},
	}

	// For now, we'll load points from the first available database
	// In the real implementation, this could be done via AJAX when database is selected
	var points []*tables.Point
	if len(databases) > 0 {
		db, err := di.NewDatabase(databases[0].Value)
		if err == nil {
			pointDao := dao.NewPointDao(db)
			points, _ = pointDao.SelectAll()
		}
	}

	formData := SearchFormData{
		Strategies: strategies,
		Databases:  databases,
		Points:     points,
	}

	c.HTML(http.StatusOK, "travel-search-form.html", gin.H{
		"data": formData,
	})
}

func (controller *TravelSearchController) SearchResult(c *gin.Context) {
	startTime := time.Now()

	// Get form parameters
	strategyType := c.PostForm("strategy")
	dbEnv := c.PostForm("database")
	source := c.PostForm("source")
	destination := c.PostForm("destination")
	arrivalFrom := c.PostForm("arrival_from")
	arrivalTo := c.PostForm("arrival_to")
	travelCountStr := c.PostForm("travel_count")

	travelCount, err := strconv.Atoi(travelCountStr)
	if err != nil {
		travelCount = 3 // default
	}

	// Parse time
	arrivalTimeFrom, err := util.TryToParseDate(arrivalFrom, []string{"2006-01-02 15:04", "2006-01-02"})
	if err != nil {
		c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
			"data": SearchResultData{Error: "Invalid arrival from time: " + err.Error()},
		})
		return
	}

	arrivalTimeTo, err := util.TryToParseDate(arrivalTo, []string{"2006-01-02 15:04", "2006-01-02"})
	if err != nil {
		c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
			"data": SearchResultData{Error: "Invalid arrival to time: " + err.Error()},
		})
		return
	}

	// Connect to database
	db, err := di.NewDatabase(dbEnv)
	if err != nil {
		c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
			"data": SearchResultData{Error: "Database connection error: " + err.Error()},
		})
		return
	}

	// Create DAO
	travelDao := dao.NewTravelDao(db)

	// Create strategy
	var strategy travel_finder.TravelSearchStrategy
	switch strategyType {
	case "simple":
		strategy = travel_finder.NewSimpleTravelSearchStrategy(travelDao)
	case "clustered":
		strategy = travel_finder.NewClusteredTravelSearchStrategy(travelDao)
	default:
		c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
			"data": SearchResultData{Error: "Unknown strategy: " + strategyType},
		})
		return
	}

	// Create filter
	filter := data.NewTravelFilter(source, destination, arrivalTimeFrom, arrivalTimeTo, travelCount)

	// Execute search
	paths, err := strategy.FindPath(filter)
	if err != nil {
		c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
			"data": SearchResultData{Error: "Search error: " + err.Error()},
		})
		return
	}

	// Get point data for display
	pointDao := dao.NewPointDao(db)
	pointsData, _ := pointDao.SelectAll()
	pointMap := make(map[string]*tables.Point)
	for _, p := range pointsData {
		pointMap[p.ID] = p
	}

	// Convert to display format
	displayPaths := make([]*TravelPath, len(paths))
	for i, path := range paths {
		transfers := make([]*TransferDisplay, len(path.Transfers))
		for j, transfer := range path.Transfers {
			fromName := transfer.From
			toName := transfer.To
			if p, ok := pointMap[transfer.From]; ok {
				fromName = fmt.Sprintf("%s (%s)", p.Name, p.BuildLocationKey())
			}
			if p, ok := pointMap[transfer.To]; ok {
				toName = fmt.Sprintf("%s (%s)", p.Name, p.BuildLocationKey())
			}

			transfers[j] = &TransferDisplay{
				From:      fromName,
				To:        toName,
				Departure: transfer.Departure.Format(time.DateTime),
				Arrival:   transfer.Arrival.Format(time.DateTime),
				Duration:  transfer.Arrival.Sub(transfer.Departure).String(),
			}
		}
		displayPaths[i] = &TravelPath{
			Transfers:     transfers,
			TotalDuration: path.TotalDuration.String(),
			TransferCount: path.TransferCount,
		}
	}

	executionTime := time.Since(startTime)

	resultData := SearchResultData{
		Strategy:      strategyType,
		Database:      dbEnv,
		Source:        source,
		Destination:   destination,
		ArrivalFrom:   arrivalFrom,
		ArrivalTo:     arrivalTo,
		TravelCount:   travelCount,
		Paths:         displayPaths,
		ExecutionTime: executionTime.String(),
	}

	c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
		"data": resultData,
	})
}

func getAvailableDatabases() []DatabaseOption {
	var databases []DatabaseOption

	// Get all .env.* files
	files, err := filepath.Glob(".env.*")
	if err != nil {
		return databases
	}

	// Try parent directory too
	parentFiles, err := filepath.Glob("../.env.*")
	if err == nil {
		files = append(files, parentFiles...)
	}

	for _, file := range files {
		// Skip files containing "root" in the name
		if strings.Contains(file, "root") {
			continue
		}

		// Extract environment name from filename
		baseName := filepath.Base(file)
		envName := strings.TrimPrefix(baseName, ".env.")

		if envName != "" && envName != baseName {
			databases = append(databases, DatabaseOption{
				Name:  envName,
				Value: envName,
			})
		}
	}

	return databases
}

// Helper to check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
