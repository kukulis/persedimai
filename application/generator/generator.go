package generator

import (
	"darbelis.eu/persedimai/tables"
	"time"
)

// Generator generates points one to a sub-square where square size and amount of squares are
// set in to this generator properties.
type Generator struct {
	n           int
	squareSize  float64
	randFactor  float64
	idGenerator IdGenerator

	travels []*tables.Travel
}

func (g *Generator) Travels() []*tables.Travel {
	return g.travels
}

//func (g *Generator) GeneratedPoints() []*tables.Point {
//	return g.generatedPoints
//}

//func (g *Generator) NeighbourPairs() []*data.PointPair {
//	return g.neighbourPairs
//}

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
			err := pointConsumer.Consume(&tables.Point{ID: id, X: x, Y: y})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) GenerateTravels() {
	// TODO
}

// GenerateTravelsForTwoPoints generates multiple travels between two points
func (g *Generator) GenerateTravelsForTwoPoints(point1 tables.Point, point2 tables.Point, fromDate time.Time, toDate time.Time, speed float64, restHours int) {
	currentDeparture := fromDate
	currentFrom := point1
	currentTo := point2

	for {
		// 1) Use GenerateSingleTravel to generate travel
		travel := g.GenerateSingleTravel(currentFrom, currentTo, currentDeparture, speed)

		// Check if arrival time is after toDate
		if travel.Arrival.After(toDate) {
			break
		}

		g.travels = append(g.travels, &travel)

		// Calculate next departure time by adding resting time to the arrival date
		nextDeparture := travel.Arrival.Add(time.Duration(restHours) * time.Hour)

		// Do these steps until 'toDate' is reached
		if toDate.Before(nextDeparture) {
			break
		}

		// 2) Repeat step with a reversed direction
		currentFrom, currentTo = currentTo, currentFrom
		currentDeparture = nextDeparture
	}
}

func (g *Generator) GenerateSingleTravel(point1 tables.Point, point2 tables.Point, fromDate time.Time, speed float64) tables.Travel {
	// Calculate distance between two points
	distance := point1.CalculateDistance(point2)

	// Calculate travel time in hours
	travelTimeHours := distance / speed

	// Calculate arrival time
	arrivalTime := fromDate.Add(time.Duration(travelTimeHours * float64(time.Hour)))

	travel := tables.Travel{
		ID:        g.idGenerator.NextId(),
		From:      point1.ID,
		To:        point2.ID,
		Departure: fromDate,
		Arrival:   arrivalTime,
	}

	return travel
}
