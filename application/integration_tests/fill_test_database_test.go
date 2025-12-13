package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"testing"
)

func TestFillTestDatabase(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	err = FillTestDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	// Test passes if no error occurred
	// Detailed logs will be printed to console showing counts
}
