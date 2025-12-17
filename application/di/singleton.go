package di

import (
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/web/api"
	"log"
)

var DatabaseInstance *database.Database = nil
var DatabasesContainerInstance database.DatabasesContainer = nil
var ApiPointsControllerInstance *api.PointsController = nil

func InitializeSingletons(defaultEnv string) {
	DatabasesContainerInstance = NewDatabasesMapContainer()
	ApiPointsControllerInstance = api.NewPointsController(DatabasesContainerInstance)

	var err error
	DatabaseInstance, err = NewDatabase(defaultEnv)
	if err != nil {
		log.Fatal(err)
	}
}
