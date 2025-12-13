package tables

import "time"

type Travel struct {
	ID        string
	From      string
	To        string
	Departure time.Time
	Arrival   time.Time
}
