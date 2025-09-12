package generator

type IdGenerator interface {
	NextId() int
}

type SimpleIdGenerator struct {
	currentId int
}

func (idg SimpleIdGenerator) NextId() int {
	idg.currentId++

	return idg.currentId
}
