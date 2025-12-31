package data

import "time"

type TravelFilter struct {
	Source          string
	Destination     string
	ArrivalTimeFrom time.Time
	ArrivalTimeTo   time.Time
	TravelCount     int
	Limit           int // default 10
	// @deprecated
	MaxWaitHoursBetweenTransits int // default 24
	MinConnectionTimeMinutes    int // default 30, minimum time between transfers for comfortable walking
	MaxConnectionTimeHours      int // default 32, maximum time allowed between connections
}

// NewTravelFilter creates a new TravelFilter with default values for Limit, MaxWaitHoursBetweenTransits, MinConnectionTimeMinutes, and MaxConnectionTimeHours
func NewTravelFilter(source, destination string, arrivalTimeFrom, arrivalTimeTo time.Time, travelCount int) *TravelFilter {
	return &TravelFilter{
		Source:          source,
		Destination:     destination,
		ArrivalTimeFrom: arrivalTimeFrom,
		ArrivalTimeTo:   arrivalTimeTo,
		TravelCount:     travelCount,
		Limit:           100,
		// @deprecated
		MaxWaitHoursBetweenTransits: 24,
		MinConnectionTimeMinutes:    30,
		MaxConnectionTimeHours:      32,
	}
}

// ValidateMaxConnectionTime validates that MaxConnectionTimeHours is one of the allowed values
// Returns true if valid, false otherwise
func (tf *TravelFilter) ValidateMaxConnectionTime(allowedValues []int) bool {
	for _, validTime := range allowedValues {
		if tf.MaxConnectionTimeHours == validTime {
			return true
		}
	}
	return false
}
