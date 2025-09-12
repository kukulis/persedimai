package generator

type IdGenerator interface {
	NextId() int
}

type SimpleIdGenerator struct {
	CurrentId int
}

func (idg *SimpleIdGenerator) NextId() int {
	idg.CurrentId++

	return idg.CurrentId
}
