package generator

import (
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/tables"
)

type Generator struct {
	n           int
	squareSize  float64
	randFactor  float64
	idGenerator *IdGenerator

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
}

func (g *Generator) GenerateTravels() {

}
