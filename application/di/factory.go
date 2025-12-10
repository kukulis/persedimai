package di

import (
	"darbelis.eu/persedimai/database"
	"github.com/joho/godotenv"
)

func NewDatabase(environment string) (*database.Database, error) {
	dbConfig, err := NewDbConfig(environment)

	return database.NewDatabase(dbConfig), err
}

func NewDbConfig(environment string) (*database.DBConfig, error) {
	envFile := GetEnvFile(environment)

	envMap, err := godotenv.Read(envFile)

	if err != nil {
		return nil, err
	}

	dbConfig := &database.DBConfig{}
	err = dbConfig.InitializeFromEnvMap(envMap)

	return dbConfig, err
}

func GetEnvFile(environment string) string {
	if environment == "prod" {
		return ".env"
	}

	return ".env." + environment
}
