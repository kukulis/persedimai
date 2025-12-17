package travel_finder

import (
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/tables"
	"fmt"
	"strings"
	"time"
)

// TravelPath represents a found path as a sequence of travels from source to destination
type TravelPath struct {
	Transfers     []*tables.Transfer
	TotalDuration time.Duration
	TotalDistance float64
	TransferCount int
}

func MakeTravelPathOfTransferSequence(sequence *tables.TransferSequence) *TravelPath {
	return &TravelPath{
		Transfers:     sequence.Transfers,
		TransferCount: sequence.TransferCount(),
		TotalDuration: sequence.TotalDuration(),
	}
}

// ToString returns a formatted string representation of the travel path
func (tp *TravelPath) ToString(pointGetter data.PointGetter) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Travel Path (%d transfer(s), Duration: %v)\n", tp.TransferCount, tp.TotalDuration))

	for i, travel := range tp.Transfers {
		fromPoint := pointGetter.GetPoint(travel.From)
		toPoint := pointGetter.GetPoint(travel.To)

		fromKey := travel.From
		toKey := travel.To

		if fromPoint != nil {
			fromKey = fromPoint.BuildLocationKey()
		}
		if toPoint != nil {
			toKey = toPoint.BuildLocationKey()
		}

		sb.WriteString(fmt.Sprintf("  %d. %s → %s\n", i+1, fromKey, toKey))
		sb.WriteString(fmt.Sprintf("     Depart: %s\n", travel.Departure.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("     Arrive: %s\n", travel.Arrival.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("     Duration: %v\n", travel.Arrival.Sub(travel.Departure)))
	}

	return sb.String()
}
