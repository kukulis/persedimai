package generator

import (
	"darbelis.eu/persedimai/tables"
	"fmt"
	"math/rand"
	"time"
)

// Word lists for generating random location names
var nameAdjectives = []string{
	"New", "Old", "Upper", "Lower", "North", "South", "East", "West",
	"Great", "Little", "High", "Deep", "Green", "Blue", "Red", "White",
	"Dark", "Bright", "Golden", "Silver", "Crystal", "Royal", "Ancient",
}

var nameNouns = []string{
	"Point", "Town", "City", "Village", "Harbor", "Port", "Bay", "Valley",
	"Hill", "Mountain", "Lake", "River", "Forest", "Field", "Meadow",
	"Bridge", "Gate", "Cross", "Square", "Haven", "Springs", "Falls",
}

// Generator generates points one to a sub-square where square size and amount of squares are
// set in to this generator properties.
type Generator struct {
	n           int
	squareSize  float64
	randFactor  float64
	idGenerator IdGenerator
}

// generateRandomName creates a random location name by combining an adjective and a noun
func (g *Generator) generateRandomName() string {
	adjective := nameAdjectives[rand.Intn(len(nameAdjectives))]
	noun := nameNouns[rand.Intn(len(nameNouns))]
	return fmt.Sprintf("%s %s", adjective, noun)
}

func (g *Generator) GeneratePoints(pointConsumer PointConsumerInterface) error {
	// let it generate objects and we will insert them using dao classes
	for i := 0; i < g.n; i++ {
		if i%2 == 1 {
			continue
		}

		for j := 0; j < g.n; j++ {
			if j%2 == 1 {
				continue
			}
			x := g.squareSize * float64(j)
			y := g.squareSize * float64(i)
			id := g.idGenerator.NextId()
			name := g.generateRandomName()
			err := pointConsumer.Consume(&tables.Point{ID: id, X: x, Y: y, Name: name})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) GenerateTravels(
	point []*tables.Point,
	fromDate time.Time,
	toDate time.Time,
	speed float64,
	restHours int,
	travelConsumer TravelConsumerInterface,
) error {
	// 0) create points pairs ids map
	pointsPairsMap := make(map[string]bool)

	// 1) put all points to a map, using Point.BuildLocationKey
	pointsMap := make(map[string]*tables.Point)
	for _, p := range point {
		key := p.BuildLocationKey()
		pointsMap[key] = p
	}

	// 2) for each point try to find 8 neighbour points:
	for _, currentPoint := range point {
		// 3) calculate neighbour point by adding to the current point (dX, dY)
		// All combinations of (dX,dY) pairs dX: [-g.squareSize*2, 0, g.squareSize*2],
		// dY: [-g.squareSize*2, 0, g.squareSize*2], except (0,0).
		deltas := []float64{-g.squareSize * 2, 0, g.squareSize * 2}

		for _, dX := range deltas {
			for _, dY := range deltas {
				// Skip (0,0)
				if dX == 0 && dY == 0 {
					continue
				}

				// Calculate neighbor point coordinates
				neighborX := currentPoint.X + dX
				neighborY := currentPoint.Y + dY

				// Build location key for the neighbor
				neighborPoint := tables.Point{X: neighborX, Y: neighborY}
				neighborKey := neighborPoint.BuildLocationKey()

				// 4) then try to find the point in the points map (made in step 1), proceed if the point is found
				foundNeighbor, exists := pointsMap[neighborKey]
				if !exists {
					continue
				}

				// 5) check the pair of points in to points pairs ids map
				// (current point and the found point) (both in reverse order too)
				pairKey1 := currentPoint.ID + "_" + foundNeighbor.ID
				pairKey2 := foundNeighbor.ID + "_" + currentPoint.ID

				// if exists skip the next step
				if pointsPairsMap[pairKey1] || pointsPairsMap[pairKey2] {
					continue
				}

				// if does not exist then put to the points pair map and proceed to the next step
				pointsPairsMap[pairKey1] = true

				// 6) call GenerateTravelsForTwoPoints
				err := g.GenerateTravelsForTwoPoints(*currentPoint, *foundNeighbor, fromDate, toDate, speed, restHours, travelConsumer)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GenerateTravelsForTwoPoints generates multiple travels between two points
func (g *Generator) GenerateTravelsForTwoPoints(point1 tables.Point, point2 tables.Point, fromDate time.Time, toDate time.Time, speed float64, restHours int, travelConsumer TravelConsumerInterface) error {
	currentDeparture := fromDate
	currentFrom := point1
	currentTo := point2

	for {
		// Apply random factor to speed for this travel
		actualSpeed := g.applyRandomFactor(speed)

		// 1) Use GenerateSingleTravel to generate travel
		travel := g.GenerateSingleTravel(currentFrom, currentTo, currentDeparture, actualSpeed)

		// Check if arrival time is after toDate
		if travel.Arrival.After(toDate) {
			break
		}

		err := travelConsumer.Consume(&travel)
		if err != nil {
			return err
		}

		// Apply random factor to rest hours for this rest period
		actualRestHours := g.applyRandomFactor(float64(restHours))

		// Calculate next departure time by adding resting time to the arrival date
		nextDeparture := travel.Arrival.Add(time.Duration(actualRestHours * float64(time.Hour)))

		// Do these steps until 'toDate' is reached
		if toDate.Before(nextDeparture) {
			break
		}

		// 2) Repeat step with a reversed direction
		currentFrom, currentTo = currentTo, currentFrom
		currentDeparture = nextDeparture
	}

	return nil
}

func (g *Generator) GenerateSingleTravel(point1 tables.Point, point2 tables.Point, fromDate time.Time, speed float64) tables.Transfer {
	// Calculate distance between two points
	distance := point1.CalculateDistance(point2)

	// Calculate travel time in hours
	travelTimeHours := distance / speed

	// Calculate arrival time
	arrivalTime := fromDate.Add(time.Duration(travelTimeHours * float64(time.Hour)))

	travel := tables.Transfer{
		ID:        g.idGenerator.NextId(),
		From:      point1.ID,
		To:        point2.ID,
		Departure: fromDate,
		Arrival:   arrivalTime,
	}

	return travel
}

// applyRandomFactor applies random variation to a value based on randFactor
// If randFactor is 0, returns the original value unchanged
// Otherwise returns value * (1 + random variation), where variation is between -randFactor and +randFactor
func (g *Generator) applyRandomFactor(value float64) float64 {
	if g.randFactor == 0 {
		return value
	}

	// Generate random variation between -randFactor and +randFactor
	variation := (rand.Float64()*2 - 1) * g.randFactor
	return value * (1 + variation)
}
