package generator

import "strconv"

type IdGenerator interface {
	NextId() string
}

type SimpleIdGenerator struct {
	CurrentId int
}

func (idg *SimpleIdGenerator) NextId() string {
	idg.CurrentId++

	return strconv.Itoa(idg.CurrentId)
}
