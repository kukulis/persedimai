package data

import "time"

type TravelFilter struct {
	Source                      string
	Destination                 string
	ArrivalTimeFrom             time.Time
	ArrivalTimeTo               time.Time
	TravelCount                 int
	Limit                       int // default 10
	MaxWaitHoursBetweenTransits int // default 24
	MinConnectionTimeMinutes    int // default 30, minimum time between transfers for comfortable walking
}

// NewTravelFilter creates a new TravelFilter with default values for Limit, MaxWaitHoursBetweenTransits, and MinConnectionTimeMinutes
func NewTravelFilter(source, destination string, arrivalTimeFrom, arrivalTimeTo time.Time, travelCount int) *TravelFilter {
	return &TravelFilter{
		Source:                      source,
		Destination:                 destination,
		ArrivalTimeFrom:             arrivalTimeFrom,
		ArrivalTimeTo:               arrivalTimeTo,
		TravelCount:                 travelCount,
		Limit:                       10,
		MaxWaitHoursBetweenTransits: 24,
		MinConnectionTimeMinutes:    30,
	}
}
