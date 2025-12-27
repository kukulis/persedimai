package aviation_edge

import (
	"fmt"
	"sync"
)

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
