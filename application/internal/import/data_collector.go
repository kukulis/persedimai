package _import

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"fmt"
	"log"
	"time"
)

// DataCollector handles data collection operations using the Aviation Edge API
type DataCollector struct {
	apiClient       *aviation_edge.AviationEdgeApiClient
	consumer        aviation_edge.ScheduleConsumer
	airportsDao     *dao.AirportsDao
	airportsMetaDao *dao.AirportsMetaDao
}

// NewDataCollector creates a new DataCollector with dependency injection
func NewDataCollector(
	apiClient *aviation_edge.AviationEdgeApiClient,
	consumer aviation_edge.ScheduleConsumer,
	airportsDao *dao.AirportsDao,
	airportsMetaDao *dao.AirportsMetaDao,
) *DataCollector {
	return &DataCollector{
		apiClient:       apiClient,
		consumer:        consumer,
		airportsDao:     airportsDao,
		airportsMetaDao: airportsMetaDao,
	}
}

func (dc *DataCollector) GetConsumer() aviation_edge.ScheduleConsumer {
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
	log.Printf("Collecting current schedules for airport: %s, day %s", airportCode, day)

	var allSchedules []aviation_edge.ScheduleResponse

	// Collect departure schedules
	log.Printf("Fetching departure schedules...")
	schedules, err := dc.apiClient.GetFutureSchedules(aviation_edge.FutureSchedulesParams{
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

	schedules = util.ArrayMap(schedules, func(s aviation_edge.ScheduleResponse) aviation_edge.ScheduleResponse {
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

// InitializeEuropeanAirportsMeta creates metadata records for all European airports
// It uses the European countries constant to find airports and initializes their metadata
func (dc *DataCollector) InitializeEuropeanAirportsMeta() error {
	log.Println("Initializing metadata for European airports...")

	// Get airports from European countries
	airports, err := dc.airportsDao.GetByCountries(util.EuropeanCountryCodes)
	if err != nil {
		return fmt.Errorf("failed to get European airports: %w", err)
	}

	log.Printf("Found %d European airports", len(airports))

	// Create metadata record for each airport
	for _, airport := range airports {
		meta := &tables.AirportMeta{
			AirportCode:  airport.CodeIataAirport,
			ImportedFrom: nil,
			ImportedTo:   nil,
		}

		err := dc.airportsMetaDao.Upsert(meta, false)
		if err != nil {
			return fmt.Errorf("failed to upsert metadata for airport %s: %w", airport.CodeIataAirport, err)
		}

		//log.Printf("Initialized metadata for airport: %s (%s, %s)",
		//	airport.CodeIataAirport, airport.NameAirport, airport.NameCountry)
	}

	log.Printf("Successfully initialized metadata for %d European airports", len(airports))

	return nil
}
