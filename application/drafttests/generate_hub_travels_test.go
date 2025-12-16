//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/integration_tests"
	"testing"
)

func TestFillTestDatabase_WithHubPoints(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}

	dbFiller := &integration_tests.DatabaseFiller{}
	err = dbFiller.FillTestDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	// Test passes if no error occurred
	// Detailed logs will be printed to console showing counts
}
