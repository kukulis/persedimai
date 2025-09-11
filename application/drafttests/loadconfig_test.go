package drafttests

import (
	"darbelis.eu/persedimai/database"
	"testing"
)

func TestLoadDbConfig(t *testing.T) {
	t.Run("Bare loading config", func(t *testing.T) {
		dbconfig := database.DBConfig{}

		dbconfig.LoadFromEnv("../.env")

		got := dbconfig.Dbname()
		want := "persedimai"

		if got != want {
			t.Errorf("Loaded dbname is wrong, got %q want %q", got, want)
		}
	})
}
