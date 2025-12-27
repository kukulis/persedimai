package aviation_edge

// ScheduleConsumer interface for processing schedule data as it's collected
type ScheduleConsumer interface {
	// Consume processes a batch of schedule responses
	// Returns error if processing fails
	Consume(schedules []ScheduleResponse) error
}
