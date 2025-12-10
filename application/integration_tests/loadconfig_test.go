package integration_tests

import (
	"darbelis.eu/persedimai/di"
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
