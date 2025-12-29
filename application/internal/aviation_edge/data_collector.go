package aviation_edge

import (
	"darbelis.eu/persedimai/internal/util"
	"fmt"
	"log"
	"time"
)

// DataCollector handles data collection operations using the Aviation Edge API
type DataCollector struct {
	apiClient *AviationEdgeApiClient
	consumer  ScheduleConsumer
}

// NewDataCollector creates a new DataCollector with dependency injection
func NewDataCollector(apiClient *AviationEdgeApiClient, consumer ScheduleConsumer) *DataCollector {
	return &DataCollector{
		apiClient: apiClient,
		consumer:  consumer,
	}
}

func (dc *DataCollector) GetConsumer() ScheduleConsumer {
	return dc.consumer
}

func (dc *DataCollector) CollectDepartureSchedules(airportCode string, dateFrom, dateTo string) error {
	dateRange, err := util.GenerateDateRange(dateFrom, dateTo)

	if err != nil {
		return err
	}
	for _, date := range dateRange {
		err = dc.CollectDepartureSchedulesForOneDay(airportCode, date)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *DataCollector) CollectDepartureSchedulesForOneDay(airportCode string, day string) error {
	log.Printf("Collecting current schedules for airport: %s", airportCode)

	var allSchedules []ScheduleResponse

	// Collect departure schedules
	log.Printf("Fetching departure schedules...")
	schedules, err := dc.apiClient.GetFutureSchedules(FutureSchedulesParams{
		IataCode: airportCode,
		Type:     "departure",
		Date:     day,
	})

	nexDay := day

	timeDay, err := time.Parse(time.DateOnly, day)
	if err == nil {
		nextTimeDay := timeDay.AddDate(0, 0, 1)
		nexDay = nextTimeDay.Format(time.DateOnly)
	}

	schedules = util.ArrayMap(schedules, func(s ScheduleResponse) ScheduleResponse {
		s.Type = "departure"
		s.Status = "future"
		s.Departure.ScheduledTime = day + " " + s.Departure.ScheduledTime

		arrivalDay := day

		if s.Arrival.ScheduledTime < s.Departure.ScheduledTime {
			arrivalDay = nexDay
		}

		s.Arrival.ScheduledTime = arrivalDay + " " + s.Arrival.ScheduledTime

		if s.Airline.Name == "" {
			s.Airline.Name = "-"
		}

		if s.Airline.IataCode == "" {
			s.Airline.IataCode = "-"
		}

		if s.Flight.IataNumber == "" {
			s.Flight.IataNumber = "-"
		}
		return s
	})

	if err != nil {
		return fmt.Errorf("failed to get departure schedules: %w", err)
	}
	log.Printf("Found %d departure schedules", len(schedules))
	allSchedules = append(allSchedules, schedules...)

	// Consume all collected schedules
	if len(allSchedules) > 0 {
		if err := dc.consumer.Consume(allSchedules); err != nil {
			return fmt.Errorf("consumer failed: %w", err)
		}
		log.Printf("Total schedules collected: %d", len(allSchedules))
	} else {
		log.Println("No schedules found")
	}

	return nil
}
