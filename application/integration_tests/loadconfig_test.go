package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestLoadDbConfig(t *testing.T) {
	t.Run("Bare loading config", func(t *testing.T) {
		dbconfig, err := di.NewDbConfig("test")
		if err != nil {
			t.Fatal(err)
		}

		got := dbconfig.Dbname
		want := "test"

		if got != want {
			t.Errorf("Loaded dbname is wrong, got %q want %q", got, want)
		}
	})
}

func TestLoadConfigWithoutFile(t *testing.T) {

	// lets say this was preset befre actual .env file reading
	err := godotenv.Load("./data/preset_test_env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbconfig, err := di.NewDbConfig("labadiena")
	if err != nil {
		t.Fatal(err)
	}

	got := dbconfig.Dbname
	want := "test"

	if got != want {
		t.Errorf("Loaded dbname is wrong, got %q want %q", got, want)
	}
}
