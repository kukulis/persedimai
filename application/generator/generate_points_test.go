package generator

import (
	"darbelis.eu/persedimai/tables"
	"testing"
)

func TestGeneratePoints(t *testing.T) {
	// TODO
	gf := GeneratorFactory{}

	idGenerator := SimpleIdGenerator{}
	g := gf.createGenerator(5, 1000, 0, idGenerator)

	g.GeneratePoints()

	generatedPoints := g.GeneratedPoints()
	//neighbourPairs := g.NeighbourPairs()

	// ID ir vardus atskirai
	expectedGeneratedPoints := []*tables.Point{
		&tables.Point{ID: 1, X: 0, Y: 0},
		&tables.Point{ID: 2, X: 2000, Y: 0},
		&tables.Point{ID: 3, X: 4000, Y: 0},
		&tables.Point{ID: 4, X: 0, Y: 2000},
		&tables.Point{ID: 5, X: 2000, Y: 2000},
		&tables.Point{ID: 6, X: 4000, Y: 2000},
		&tables.Point{ID: 7, X: 0, Y: 4000},
		&tables.Point{ID: 8, X: 2000, Y: 4000},
		&tables.Point{ID: 9, X: 4000, Y: 4000},
	}

	if len(generatedPoints) != len(expectedGeneratedPoints) {
		t.Errorf("GeneratePoints returned %d generated points, expected %d", len(generatedPoints), len(expectedGeneratedPoints))
	}

}
