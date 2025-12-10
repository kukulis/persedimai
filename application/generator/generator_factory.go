package generator

type GeneratorFactory struct {
}

func (g *GeneratorFactory) CreateGenerator(
	n int,
	squareSize float64,
	randFactor float64,
	idGenerator IdGenerator,
) *Generator {
	return &Generator{
		n:           n,
		squareSize:  squareSize,
		randFactor:  randFactor,
		idGenerator: idGenerator,
	}
}
