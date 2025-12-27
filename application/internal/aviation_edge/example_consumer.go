package aviation_edge

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Example implementations of ScheduleConsumer interface

// PrintScheduleConsumer prints schedules to stdout
type PrintScheduleConsumer struct {
	TotalCount int
	mu         sync.Mutex
}

func (p *PrintScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, schedule := range schedules {
		fmt.Printf("Flight %s (%s): %s -> %s | Type: %s | Status: %s\n",
			schedule.Flight.IataNumber,
			schedule.Airline.Name,
			schedule.Departure.IataCode,
			schedule.Arrival.IataCode,
			schedule.Type,
			schedule.Status)
	}
	p.TotalCount += len(schedules)
	return nil
}

// SliceScheduleConsumer collects all schedules into a slice
type SliceScheduleConsumer struct {
	Schedules []ScheduleResponse
	mu        sync.Mutex
}

func (s *SliceScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Schedules = append(s.Schedules, schedules...)
	return nil
}

// FileScheduleConsumer writes schedules to a JSON file
type FileScheduleConsumer struct {
	FilePath   string
	Schedules  []ScheduleResponse
	mu         sync.Mutex
	autoFlush  bool
	flushCount int
}

func NewFileScheduleConsumer(filePath string, autoFlush bool) *FileScheduleConsumer {
	return &FileScheduleConsumer{
		FilePath:   filePath,
		Schedules:  make([]ScheduleResponse, 0),
		autoFlush:  autoFlush,
		flushCount: 1000, // Flush every 1000 schedules if autoFlush is true
	}
}

func (f *FileScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.Schedules = append(f.Schedules, schedules...)

	// Auto-flush if enabled and threshold reached
	if f.autoFlush && len(f.Schedules) >= f.flushCount {
		return f.flush()
	}

	return nil
}

func (f *FileScheduleConsumer) flush() error {
	if len(f.Schedules) == 0 {
		return nil
	}

	data, err := json.MarshalIndent(f.Schedules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schedules: %w", err)
	}

	if err := os.WriteFile(f.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Flushed %d schedules to %s", len(f.Schedules), f.FilePath)
	return nil
}

// Flush manually writes all collected schedules to file
func (f *FileScheduleConsumer) Flush() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.flush()
}

// DatabaseScheduleConsumer example (placeholder for database integration)
type DatabaseScheduleConsumer struct {
	// db *sql.DB or other database connection
	InsertedCount int
	mu            sync.Mutex
}

func (d *DatabaseScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// TODO: Implement database insert logic
	// For now, just count
	for _, schedule := range schedules {
		// Example: INSERT INTO schedules (flight_number, airline, ...) VALUES (?, ?, ...)
		_ = schedule
		d.InsertedCount++
	}

	log.Printf("Inserted %d schedules into database (total: %d)", len(schedules), d.InsertedCount)
	return nil
}

// FilteredScheduleConsumer filters schedules before passing to another consumer
type FilteredScheduleConsumer struct {
	FilterFunc func(ScheduleResponse) bool
	Next       ScheduleConsumer
}

func (f *FilteredScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	filtered := make([]ScheduleResponse, 0)
	for _, schedule := range schedules {
		if f.FilterFunc(schedule) {
			filtered = append(filtered, schedule)
		}
	}

	if len(filtered) > 0 && f.Next != nil {
		return f.Next.Consume(filtered)
	}

	return nil
}

// MultiScheduleConsumer broadcasts to multiple consumers
type MultiScheduleConsumer struct {
	Consumers []ScheduleConsumer
}

func (m *MultiScheduleConsumer) Consume(schedules []ScheduleResponse) error {
	for _, consumer := range m.Consumers {
		if err := consumer.Consume(schedules); err != nil {
			log.Printf("Warning: Consumer failed: %v", err)
		}
	}
	return nil
}

// Example usage function
func ExampleDataCollectorUsage() {
	// Initialize API client
	apiClient := NewAviationEdgeApiClient("your-api-key")

	// Create data collector with dependency injection
	collector := NewDataCollector(apiClient)

	// Example 1: Print schedules to stdout
	printConsumer := &PrintScheduleConsumer{}
	err := collector.CollectSchedules(CollectSchedulesParams{
		CountryCode:       "US",
		StartDate:         "2025-12-20",
		EndDate:           "2025-12-21",
		IncludeDepartures: true,
		IncludeArrivals:   false,
		Consumer:          printConsumer,
		RateLimitDelay:    2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total schedules printed: %d", printConsumer.TotalCount)

	// Example 2: Collect schedules into a slice
	sliceConsumer := &SliceScheduleConsumer{}
	err = collector.CollectSchedules(CollectSchedulesParams{
		CountryCode:       "GB",
		StartDate:         "2025-12-27",
		EndDate:           "2025-12-27",
		IncludeDepartures: true,
		IncludeArrivals:   true,
		Consumer:          sliceConsumer,
		RateLimitDelay:    1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Collected %d schedules into slice", len(sliceConsumer.Schedules))

	// Example 3: Save schedules to file
	fileConsumer := NewFileScheduleConsumer("schedules.json", false)
	err = collector.CollectSchedules(CollectSchedulesParams{
		CountryCode:       "FR",
		StartDate:         "2025-12-25",
		EndDate:           "2025-12-26",
		IncludeDepartures: true,
		IncludeArrivals:   true,
		Consumer:          fileConsumer,
		RateLimitDelay:    1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	fileConsumer.Flush() // Don't forget to flush!

	// Example 4: Filter schedules (only active flights)
	filteredConsumer := &FilteredScheduleConsumer{
		FilterFunc: func(s ScheduleResponse) bool {
			return s.Status == "active" || s.Status == "scheduled"
		},
		Next: &PrintScheduleConsumer{},
	}
	collector.CollectSchedules(CollectSchedulesParams{
		CountryCode:       "DE",
		StartDate:         "2025-12-27",
		EndDate:           "2025-12-27",
		IncludeDepartures: true,
		IncludeArrivals:   false,
		Consumer:          filteredConsumer,
	})

	// Example 5: Multiple consumers (print and save)
	multiConsumer := &MultiScheduleConsumer{
		Consumers: []ScheduleConsumer{
			&PrintScheduleConsumer{},
			NewFileScheduleConsumer("backup.json", false),
			&DatabaseScheduleConsumer{},
		},
	}
	collector.CollectSchedules(CollectSchedulesParams{
		CountryCode:       "IT",
		StartDate:         "2025-12-27",
		EndDate:           "2025-12-27",
		IncludeDepartures: true,
		IncludeArrivals:   true,
		Consumer:          multiConsumer,
	})
}
