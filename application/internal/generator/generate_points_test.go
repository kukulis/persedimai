package generator

import (
	"darbelis.eu/persedimai/internal/tables"
	"testing"
)

func TestGeneratePoints(t *testing.T) {
	gf := GeneratorFactory{}

	var idGenerator IdGenerator
	idGenerator = &SimpleIdGenerator{}
	g := gf.CreateGenerator(5, 1000, 0, idGenerator)

	pointArrayConsumer := NewPointArrayConsumer()
	g.GeneratePoints(pointArrayConsumer)

	generatedPoints := pointArrayConsumer.Points
	//neighbourPairs := g.NeighbourPairs()

	// ID ir vardus atskirai
	expectedGeneratedPoints := []*tables.Point{
		{ID: "1", X: 0, Y: 0},
		{ID: "2", X: 2000, Y: 0},
		{ID: "3", X: 4000, Y: 0},
		{ID: "4", X: 0, Y: 2000},
		{ID: "5", X: 2000, Y: 2000},
		{ID: "6", X: 4000, Y: 2000},
		{ID: "7", X: 0, Y: 4000},
		{ID: "8", X: 2000, Y: 4000},
		{ID: "9", X: 4000, Y: 4000},
	}

	if len(generatedPoints) != len(expectedGeneratedPoints) {
		t.Errorf("GeneratePoints returned %d generated points, expected %d", len(generatedPoints), len(expectedGeneratedPoints))
	}

	for i := 0; i < len(generatedPoints); i++ {
		generatedPoint := generatedPoints[i]
		expectedGeneratedPoint := expectedGeneratedPoints[i]

		if generatedPoint.ID != expectedGeneratedPoint.ID {
			t.Errorf("GeneratePoints result at index %d, ID = %s, expected %s", i, generatedPoint.ID, expectedGeneratedPoint.ID)
		}
		if generatedPoint.X != expectedGeneratedPoint.X {
			t.Errorf("GeneratePoints result at index %d, X = %f, expected %f", i, generatedPoint.X, expectedGeneratedPoint.X)
		}
		if generatedPoint.Y != expectedGeneratedPoint.Y {
			t.Errorf("GeneratePoints result at index %d, Y = %f, expected %f", i, generatedPoint.Y, expectedGeneratedPoint.Y)
		}
		if generatedPoint.Name == "" {
			t.Errorf("GeneratePoints result at index %d has empty Name", i)
		}
	}

	// TODO neighbours
}
