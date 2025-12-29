package di

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/web/api"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

var DatabaseInstance *database.Database = nil
var DatabasesContainerInstance database.DatabasesContainer = nil
var ApiPointsControllerInstance *api.PointsController = nil
var ApiTravelsControllerInstance *api.TravelsController = nil
var ApiKey string

var instances = map[string]interface{}{}

func Wrap[T any](loader func() T) T {
	name := reflect.TypeOf(loader).String()

	if strings.HasPrefix(name, "func() ") {
		name = name[7:]
	}

	return WrapNamed(name, loader)
}

func WrapNamed[T any](name string, loader func() T) T {

	fmt.Printf("Loader name [%s] \n", name)

	instance := instances[name]
	if instance == nil {
		fmt.Printf("Creating new instance[%s] \n", name)
		instances[name] = loader()
	} else {
		fmt.Printf("Found existing instance[%s] \n", name)
	}

	return instances[name].(T)
}

func GetInstance(name string) interface{} {
	return instances[name]
}

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
}

func GetAviationEdgeClient() *aviation_edge.AviationEdgeApiClient {
	return aviation_edge.NewAviationEdgeApiClient(ApiKey)
}

func GetDataCollector() *aviation_edge.DataCollector {
	return aviation_edge.NewDataCollector(Wrap(GetAviationEdgeClient), Wrap(GetScheduleConsumer))
}

func GetScheduleConsumer() aviation_edge.ScheduleConsumer {
	return dao.NewDatabaseScheduleConsumer(Wrap(GetFlightSchedulesDao))
}

func GetFlightSchedulesDao() *dao.AviationEdgeFlightSchedulesDao {
	return dao.NewAviationEdgeFlightSchedulesDao(DatabaseInstance)
}

func GetAirportsDao() *dao.AirportsDao {
	return dao.NewAirportsDao(DatabaseInstance)
}
