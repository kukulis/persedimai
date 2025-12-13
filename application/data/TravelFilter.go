package data

import "time"

type TravelFilter struct {
	Source          string
	Destination     string
	ArrivalTimeFrom time.Time
	ArrivalTimeTo   time.Time
	TravelCount     int
}
