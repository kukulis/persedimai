package di

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/web/api"
	"fmt"
	"log"
	"os"
)

var DatabaseInstance *database.Database = nil
var DatabasesContainerInstance database.DatabasesContainer = nil
var ApiPointsControllerInstance *api.PointsController = nil
var ApiTravelsControllerInstance *api.TravelsController = nil
var ApiKey string
var ApiClientLoader func() *aviation_edge.AviationEdgeApiClient
var DataCollectorLoader func() *aviation_edge.DataCollector
var ScheduleConsumerLoader func() aviation_edge.ScheduleConsumer

func InitializeSingletons(defaultEnv string) {
	DatabasesContainerInstance = NewDatabasesMapContainer()
	ApiPointsControllerInstance = api.NewPointsController(DatabasesContainerInstance)
	ApiTravelsControllerInstance = api.NewTravelsController(DatabasesContainerInstance)

	var err error
	DatabaseInstance, err = NewDatabase(defaultEnv)
	if err != nil {
		log.Fatal(err)
	}

	ApiKey = os.Getenv("AVIATION_EDGE_API_KEY")
	if ApiKey == "" {
		fmt.Println("Error: AVIATION_EDGE_API_KEY is not set")
		fmt.Println("Set it in .env file or export AVIATION_EDGE_API_KEY=your_key")
		os.Exit(1)
	}

	ApiClientLoader = func() *aviation_edge.AviationEdgeApiClient {
		return aviation_edge.NewAviationEdgeApiClient(ApiKey)
	}

	ScheduleConsumerLoader = func() aviation_edge.ScheduleConsumer {
		return &aviation_edge.PrintScheduleConsumer{}
	}

	DataCollectorLoader = func() *aviation_edge.DataCollector {
		return aviation_edge.NewDataCollector(ApiClientLoader(), ScheduleConsumerLoader())
	}

}
