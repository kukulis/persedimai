package aviation_edge

import (
	"fmt"
	"log"
	"time"
)

// ScheduleConsumer interface for processing schedule data as it's collected
type ScheduleConsumer interface {
	// Consume processes a batch of schedule responses
	// Returns error if processing fails
	Consume(schedules []ScheduleResponse) error
}

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

// CollectSchedulesParams contains parameters for schedule collection
type CollectSchedulesParams struct {
	CountryCode       string
	StartDate         string // Format: "2025-12-27"
	EndDate           string // Format: "2025-12-30"
	IncludeDepartures bool
	IncludeArrivals   bool
	Consumer          ScheduleConsumer
	RateLimitDelay    time.Duration // Delay between API calls to avoid rate limiting
}

// CollectSchedules collects flight schedules for all airports in a country for a given period
// It retrieves airport data and then collects schedules for each airport/date combination
func (dc *DataCollector) CollectSchedules(params CollectSchedulesParams) error {
	// Validate parameters
	if params.CountryCode == "" {
		return fmt.Errorf("country code is required")
	}
	if params.Consumer == nil {
		return fmt.Errorf("schedule consumer is required")
	}
	if !params.IncludeDepartures && !params.IncludeArrivals {
		return fmt.Errorf("must include at least departures or arrivals")
	}

	// Set default rate limit delay if not specified
	if params.RateLimitDelay == 0 {
		params.RateLimitDelay = 1 * time.Second
	}

	// Step 1: Get all airports in the country
	log.Printf("Fetching airports for country: %s", params.CountryCode)
	airports, err := dc.apiClient.GetAirportsByCountry(params.CountryCode)
	if err != nil {
		return fmt.Errorf("failed to get airports for country %s: %w", params.CountryCode, err)
	}
	log.Printf("Found %d airports in %s", len(airports), params.CountryCode)

	// Step 2: Generate date range
	dates, err := generateDateRange(params.StartDate, params.EndDate)
	if err != nil {
		return fmt.Errorf("failed to generate date range: %w", err)
	}
	log.Printf("Collecting schedules for %d dates (%s to %s)", len(dates), params.StartDate, params.EndDate)

	// Step 3: Iterate through airports and dates
	totalSchedules := 0
	for airportIdx, airport := range airports {
		// Skip airports without IATA code
		if airport.CodeIataAirport == "" {
			log.Printf("Skipping airport without IATA code: %s", airport.NameAirport)
			continue
		}

		log.Printf("[%d/%d] Processing airport: %s (%s)",
			airportIdx+1, len(airports), airport.NameAirport, airport.CodeIataAirport)

		for _, date := range dates {
			// Collect departure schedules
			if params.IncludeDepartures {
				if err := dc.collectAndConsume(airport.CodeIataAirport, date, "departure", params.Consumer); err != nil {
					log.Printf("Warning: Failed to collect departures for %s on %s: %v",
						airport.CodeIataAirport, date, err)
				} else {
					totalSchedules++
				}

				// Rate limiting
				time.Sleep(params.RateLimitDelay)
			}

			// Collect arrival schedules
			if params.IncludeArrivals {
				if err := dc.collectAndConsume(airport.CodeIataAirport, date, "arrival", params.Consumer); err != nil {
					log.Printf("Warning: Failed to collect arrivals for %s on %s: %v",
						airport.CodeIataAirport, date, err)
				} else {
					totalSchedules++
				}

				// Rate limiting
				time.Sleep(params.RateLimitDelay)
			}
		}
	}

	log.Printf("Schedule collection completed. Total API calls: %d", totalSchedules)
	return nil
}

// collectAndConsume is a helper method that collects schedules and passes them to the consumer
func (dc *DataCollector) collectAndConsume(airportIata, date, scheduleType string, consumer ScheduleConsumer) error {
	schedules, err := dc.apiClient.GetHistoricalSchedules(map[string]string{
		"code":      airportIata,
		"type":      scheduleType,
		"date_from": date,
		"date_to":   date, // Same date for single day
	})
	if err != nil {
		return err
	}

	// Only consume if there are schedules
	if len(schedules) > 0 {
		if err := consumer.Consume(schedules); err != nil {
			return fmt.Errorf("consumer failed to process schedules: %w", err)
		}
		log.Printf("  Collected and consumed %d %s schedules for %s on %s",
			len(schedules), scheduleType, airportIata, date)
	}

	return nil
}

// generateDateRange generates a slice of date strings between startDate and endDate (inclusive)
// Date format: "2025-12-27"
func generateDateRange(startDate, endDate string) ([]string, error) {
	if startDate == "" || endDate == "" {
		return []string{time.Now().Format(time.DateOnly)}, nil
	}

	start, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	end, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(time.DateOnly))
	}

	return dates, nil
}
