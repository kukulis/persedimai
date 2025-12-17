package di

import (
	"darbelis.eu/persedimai/internal/database"
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func NewDatabase(environment string) (*database.Database, error) {
	dbConfig, err := NewDbConfig(environment)

	return database.NewDatabase(dbConfig), err
}

func NewDbConfig(environment string) (*database.DBConfig, error) {
	envFile := GetEnvFile(environment)

	envMap := map[string]string{}
	var err error

	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		// let's try the directory up
		envFile = "../" + envFile
	}

	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		envFile = ""
	}

	if envFile != "" {
		envMap, err = godotenv.Read(envFile)
		if err != nil {
			return nil, err
		}
	}
	dbConfig := &database.DBConfig{}

	// nothing was loaded trying to get the config data directly from env
	if len(envMap) == 0 {
		for _, paramName := range dbConfig.GetRequiredParamsNames() {
			envMap[paramName] = os.Getenv(paramName)
		}
	}

	err = dbConfig.InitializeFromEnvMap(envMap)

	return dbConfig, err
}

func GetEnvFile(environment string) string {
	if environment == "prod" {
		return ".env"
	}

	return ".env." + environment
}
