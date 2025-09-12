package factory

import (
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/generator"
)

type MainFactory struct {
	environment string

	dbConfig    *database.DBConfig
	database    *database.Database
	dbGenerator *generator.Generator
}

func (m MainFactory) Environment() string {
	return m.environment
}

func (m MainFactory) DbConfig() *database.DBConfig {
	return m.dbConfig
}

func (m MainFactory) Database() *database.Database {
	return m.database
}

func (m MainFactory) DbGenerator() *generator.Generator {
	return m.dbGenerator
}

func NewMainFactory(environment string) *MainFactory {
	return &MainFactory{environment: environment}
}
