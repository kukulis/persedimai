package generator

import (
	"darbelis.eu/persedimai/tables"
	"reflect"
	"testing"
)

func TestGeneratePoints(t *testing.T) {
	gf := GeneratorFactory{}

	var idGenerator IdGenerator
	idGenerator = &SimpleIdGenerator{}
	g := gf.createGenerator(5, 1000, 0, idGenerator)

	g.GeneratePoints()

	generatedPoints := g.GeneratedPoints()
	//neighbourPairs := g.NeighbourPairs()

	// ID ir vardus atskirai
	expectedGeneratedPoints := []*tables.Point{
		{ID: 1, X: 0, Y: 0},
		{ID: 2, X: 2000, Y: 0},
		{ID: 3, X: 4000, Y: 0},
		{ID: 4, X: 0, Y: 2000},
		{ID: 5, X: 2000, Y: 2000},
		{ID: 6, X: 4000, Y: 2000},
		{ID: 7, X: 0, Y: 4000},
		{ID: 8, X: 2000, Y: 4000},
		{ID: 9, X: 4000, Y: 4000},
	}

	if len(generatedPoints) != len(expectedGeneratedPoints) {
		t.Errorf("GeneratePoints returned %d generated points, expected %d", len(generatedPoints), len(expectedGeneratedPoints))
	}

	for i := 0; i < len(generatedPoints); i++ {
		generatedPoint := generatedPoints[i]
		expectedGeneratedPoint := expectedGeneratedPoints[i]

		if !reflect.DeepEqual(expectedGeneratedPoint, generatedPoint) {
			t.Errorf("GeneratePoints result at index %d,  %v not equal expected %v", i, generatedPoint, expectedGeneratedPoint)
		}
	}

	// TODO neighbours
}
