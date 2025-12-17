package tables

import "time"

type Transfer struct {
	ID        string
	From      string
	To        string
	Departure time.Time
	Arrival   time.Time
}
