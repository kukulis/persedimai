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

	//generatedPoints []*tables.Point
	//neighbourPairs []*data.PointPair
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
	// 1) calculate distance between two points
	// 2) calculate time to travel between two points
	// 3) increase time moment by the travel length and by the rest time, then calculate back travel
}

func (g *Generator) GenerateSingleTravel(point1 tables.Point, point2 tables.Point, fromDate time.Time, speed float64) tables.Travel {
	travel := tables.Travel{
		ID:        "",
		From:      "",
		To:        "",
		Departure: time.Time{},
		Arrival:   time.Time{},
	}
	// TODO

	return travel

}
