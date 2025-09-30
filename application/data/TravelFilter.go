package data

import "time"

type TravelFilter struct {
	Source          int
	Destination     int
	ArrivalTimeFrom time.Time
	ArrivalTimeTo   time.Time
	TravelCount     int
}
