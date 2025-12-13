package generator

import "testing"

func TestIdGenerator(t *testing.T) {
	idGenerator := SimpleIdGenerator{}

	id := idGenerator.NextId()
	AssertIdEquals(t, "1", id)

	id = idGenerator.NextId()
	AssertIdEquals(t, "2", id)

	id = idGenerator.NextId()
	AssertIdEquals(t, "3", id)

}

func AssertIdEquals(t *testing.T, want string, got string) {
	if want != got {
		t.Errorf("Id generator fail expected %s, actual %s", want, got)
	}
}
