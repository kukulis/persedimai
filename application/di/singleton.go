package di

import (
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/web/api"
	"log"
)

var DatabaseInstance *database.Database = nil
var DatabasesContainerInstance database.DatabasesContainer = nil
var ApiPointsControllerInstance *api.PointsController = nil
var ApiTravelsControllerInstance *api.TravelsController = nil

func InitializeSingletons(defaultEnv string) {
	DatabasesContainerInstance = NewDatabasesMapContainer()
	ApiPointsControllerInstance = api.NewPointsController(DatabasesContainerInstance)
	ApiTravelsControllerInstance = api.NewTravelsController(DatabasesContainerInstance)

	var err error
	DatabaseInstance, err = NewDatabase(defaultEnv)
	if err != nil {
		log.Fatal(err)
	}
}
