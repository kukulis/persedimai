package dao

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"log"
	"sync"
)

// DatabaseScheduleConsumer saves schedules to the database using AviationEdgeFlightSchedulesDao
type DatabaseScheduleConsumer struct {
	dao        *AviationEdgeFlightSchedulesDao
	TotalCount int
	mu         sync.Mutex
}

// TODO call from di package.

// NewDatabaseScheduleConsumer creates a new DatabaseScheduleConsumer with the given DAO
func NewDatabaseScheduleConsumer(dao *AviationEdgeFlightSchedulesDao) *DatabaseScheduleConsumer {
	return &DatabaseScheduleConsumer{
		dao:        dao,
		TotalCount: 0,
	}
}

// Consume saves the schedules to the database using the DAO's UpsertFlightSchedules method
func (d *DatabaseScheduleConsumer) Consume(schedules []aviation_edge.ScheduleResponse) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(schedules) == 0 {
		return nil
	}

	// Convert to pointers for the DAO
	schedulePtrs := make([]*aviation_edge.ScheduleResponse, len(schedules))
	for i := range schedules {
		schedulePtrs[i] = &schedules[i]
	}

	// Upsert schedules to database
	err := d.dao.UpsertFlightSchedules(schedulePtrs)
	if err != nil {
		return err
	}

	d.TotalCount += len(schedules)
	log.Printf("Inserted/updated %d schedules to database (total: %d)", len(schedules), d.TotalCount)

	return nil
}
