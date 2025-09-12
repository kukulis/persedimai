package generator

import (
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/tables"
)

type Generator struct {
	n           int
	squareSize  float64
	randFactor  float64
	idGenerator IdGenerator

	generatedPoints []*tables.Point
	neighbourPairs  []*data.PointPair
}

func (g *Generator) GeneratedPoints() []*tables.Point {
	return g.generatedPoints
}

func (g *Generator) NeighbourPairs() []*data.PointPair {
	return g.neighbourPairs
}

func (g *Generator) GeneratePoints() {
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
			g.generatedPoints = append(g.generatedPoints, &tables.Point{ID: id, X: x, Y: y})

		}
	}
}

func (g *Generator) GenerateTravels() {

}
