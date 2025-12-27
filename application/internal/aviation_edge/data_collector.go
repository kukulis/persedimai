package aviation_edge

import (
	"darbelis.eu/persedimai/internal/util"
	"fmt"
	"log"
)

// DataCollector handles data collection operations using the Aviation Edge API
type DataCollector struct {
	apiClient *AviationEdgeApiClient
}

// NewDataCollector creates a new DataCollector with dependency injection
func NewDataCollector(apiClient *AviationEdgeApiClient) *DataCollector {
	return &DataCollector{
		apiClient: apiClient,
	}
}

func (dc *DataCollector) CollectDepartureSchedules(airportCode string, dateFrom, dateTo string, consumer ScheduleConsumer) error {
	dateRange, err := util.GenerateDateRange(dateFrom, dateTo)

	if err != nil {
		return err
	}
	for _, date := range dateRange {
		err = dc.CollectDepartureSchedulesForOneDay(airportCode, date, consumer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *DataCollector) CollectDepartureSchedulesForOneDay(airportCode string, day string, consumer ScheduleConsumer) error {
	log.Printf("Collecting current schedules for airport: %s", airportCode)

	var allSchedules []ScheduleResponse

	// Collect departure schedules
	log.Printf("Fetching departure schedules...")
	schedules, err := dc.apiClient.GetFutureSchedules(FutureSchedulesParams{
		IataCode: airportCode,
		Type:     "departure",
		Date:     day,
	})
	if err != nil {
		return fmt.Errorf("failed to get departure schedules: %w", err)
	}
	log.Printf("Found %d departure schedules", len(schedules))
	allSchedules = append(allSchedules, schedules...)

	// Consume all collected schedules
	if len(allSchedules) > 0 {
		if err := consumer.Consume(allSchedules); err != nil {
			return fmt.Errorf("consumer failed: %w", err)
		}
		log.Printf("Total schedules collected: %d", len(allSchedules))
	} else {
		log.Println("No schedules found")
	}

	return nil
}
