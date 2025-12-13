package generator

import "testing"

func TestIdGenerator(t *testing.T) {
	t.Run("SimpleIdGenerator", func(t *testing.T) {
		idGenerator := SimpleIdGenerator{}

		id := idGenerator.NextId()
		AssertIdEquals(t, "1", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "2", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "3", id)
	})

	t.Run("UUIDGenerator", func(t *testing.T) {
		idGenerator := &UUIDGenerator{}

		id1 := idGenerator.NextId()
		id2 := idGenerator.NextId()
		id3 := idGenerator.NextId()

		// UUIDs should not be empty
		if id1 == "" || id2 == "" || id3 == "" {
			t.Errorf("UUID generator returned empty ID")
		}

		// UUIDs should be unique
		if id1 == id2 || id2 == id3 || id1 == id3 {
			t.Errorf("UUID generator returned duplicate IDs: %s, %s, %s", id1, id2, id3)
		}

		// UUIDs should be 36 characters long (with hyphens)
		if len(id1) != 36 {
			t.Errorf("UUID should be 36 characters long, got %d for %s", len(id1), id1)
		}
	})

	t.Run("SequenceGenerator", func(t *testing.T) {
		idSequence := []string{"alpha", "beta", "gamma", "delta"}
		idGenerator := NewSequenceGenerator(idSequence)

		id := idGenerator.NextId()
		AssertIdEquals(t, "alpha", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "beta", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "gamma", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "delta", id)

		// After exhausting the sequence, should return empty string
		id = idGenerator.NextId()
		AssertIdEquals(t, "", id)
	})

	t.Run("SequenceGeneratorEmpty", func(t *testing.T) {
		idSequence := []string{}
		idGenerator := NewSequenceGenerator(idSequence)

		id := idGenerator.NextId()
		AssertIdEquals(t, "", id)
	})

	t.Run("SequenceGeneratorSingleElement", func(t *testing.T) {
		idSequence := []string{"single"}
		idGenerator := NewSequenceGenerator(idSequence)

		id := idGenerator.NextId()
		AssertIdEquals(t, "single", id)

		id = idGenerator.NextId()
		AssertIdEquals(t, "", id)
	})
}

func AssertIdEquals(t *testing.T, want string, got string) {
	if want != got {
		t.Errorf("Id generator fail expected %s, actual %s", want, got)
	}
}
