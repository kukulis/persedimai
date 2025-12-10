package generator

import (
	"testing"
)

func TestGenerateTravels(t *testing.T) {
	gf := GeneratorFactory{}

	g := gf.CreateGenerator(5, 1000, 0, &SimpleIdGenerator{})

	pointArrayConsumer := NewPointArrayConsumer()

	g.GeneratePoints(pointArrayConsumer)

	g.GenerateTravels()

	//travels := g.Travels()
	//
	//expectedTravels := []tables.Travel {
	//	tables.Travel{}
	//}

	// TODO asserts

}
