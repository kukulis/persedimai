package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"testing"
)

func TestFillTestDatabase(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	err = integration_tests.FillTestDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	// Test passes if no error occurred
	// Detailed logs will be printed to console showing counts
}
