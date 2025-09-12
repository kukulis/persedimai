package generator

import (
	"testing"
)

func TestGenerateTravels(t *testing.T) {
	gf := GeneratorFactory{}

	g := gf.createGenerator(5, 1000, 0, &SimpleIdGenerator{})

	g.GeneratePoints()

	g.GenerateTravels()

	//travels := g.Travels()
	//
	//expectedTravels := []tables.Travel {
	//	tables.Travel{}
	//}

	// TODO

}
