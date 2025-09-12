package tables

import "time"

type Travel struct {
	ID        int
	From      int
	To        int
	Departure time.Time
	Arrival   time.Time
}
