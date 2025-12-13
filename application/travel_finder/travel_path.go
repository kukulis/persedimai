package travel_finder

import (
	"darbelis.eu/persedimai/tables"
	"time"
)

// TravelPath represents a found path as a sequence of travels from source to destination
type TravelPath struct {
	Travels       []*tables.Transfer
	TotalDuration time.Duration
	TotalDistance float64
	TransferCount int
}
