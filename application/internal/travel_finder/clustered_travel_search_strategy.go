package travel_finder

import (
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
	"errors"
	"fmt"
	"time"
)

// ClusteredTravelSearchStrategy implements a clustered travel search strategy using time-clustered data
type ClusteredTravelSearchStrategy struct {
	travelDao *dao.TravelDao
}

// NewClusteredTravelSearchStrategy creates a new ClusteredTravelSearchStrategy
func NewClusteredTravelSearchStrategy(travelDao *dao.TravelDao) *ClusteredTravelSearchStrategy {
	return &ClusteredTravelSearchStrategy{
		travelDao: travelDao,
	}
}

// FindPath finds a sequence of travels from source to destination based on the filter criteria
// Uses clustered data tables for improved performance
func (s *ClusteredTravelSearchStrategy) FindPath(filter *data.TravelFilter) ([]*TravelPath, error) {
	var sequences []*tables.TransferSequence
	var err error

	fmt.Printf("ClusteredTravelSearchStrategy FindPath called, travel filter: %v\n", filter)

	switch filter.TravelCount {
	case 1:
		sequences, err = s.travelDao.FindPathSimple1(filter)
	case 2:
		if filter.MaxConnectionTimeHours <= 8 {
			sequences, err = s.travelDao.FindPathClustered2(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		} else {
			sequences, err = s.travelDao.FindPath8Clustered2(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		}
	case 3:
		if filter.MaxConnectionTimeHours <= 8 {
			sequences, err = s.travelDao.FindPathClustered3(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		} else {
			sequences, err = s.travelDao.FindPath8Clustered3(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		}
	case 4:
		if filter.MaxConnectionTimeHours <= 8 {
			sequences, err = s.travelDao.FindPathClustered4(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		} else {
			sequences, err = s.travelDao.FindPath8Clustered4(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		}
	case 5:
		if filter.MaxConnectionTimeHours <= 8 {
			sequences, err = s.travelDao.FindPathClustered5(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		} else {
			sequences, err = s.travelDao.FindPath8Clustered5(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo, filter.MaxConnectionTimeHours, filter.Limit)
		}
	default:
		if filter.TravelCount > 5 {
			return nil, errors.New("unimplemented: TravelCount > 5 not supported")
		}
		return nil, errors.New("invalid TravelCount: must be 1, 2, 3, 4, or 5")
	}

	if err != nil {
		return nil, err
	}

	if sequences == nil || len(sequences) == 0 {
		return nil, nil
	}

	// Reload actual transfers from database to get precise timestamps
	err = s.reloadActualTransfers(sequences)
	if err != nil {
		return nil, err
	}

	// Filter sequences by location connectivity and minimum connection time
	minConnectionTime := time.Duration(filter.MinConnectionTimeMinutes) * time.Minute
	filteredSequences := util.ArrayFilter(sequences, func(sequence *tables.TransferSequence) bool {
		return sequence.AreLocationsConnected() && sequence.ValidateMinConnectionTime(minConnectionTime)
		// TODO check the arrival time range too
	})

	travelPaths := util.ArrayMap(filteredSequences, func(seq *tables.TransferSequence) *TravelPath {
		return MakeTravelPathOfTransferSequence(seq)
	})

	return travelPaths, nil
}

// reloadActualTransfers loads actual transfer data from the database and replaces
// the cluster-based transfers in sequences with precise timestamp data
func (s *ClusteredTravelSearchStrategy) reloadActualTransfers(sequences []*tables.TransferSequence) error {
	// Collect all transfer IDs from all sequences
	transferIDsMap := make(map[string]bool)
	for _, seq := range sequences {
		for _, transfer := range seq.Transfers {
			transferIDsMap[transfer.ID] = true
		}
	}

	// Convert map keys to slice
	transferIDs := make([]string, 0, len(transferIDsMap))
	for id := range transferIDsMap {
		transferIDs = append(transferIDs, id)
	}

	// Load actual transfers from database
	actualTransfers, err := s.travelDao.FindByIDs(transferIDs)
	if err != nil {
		return err
	}

	// Create map of transfers by ID
	transfersMap := make(map[string]*tables.Transfer)
	for _, transfer := range actualTransfers {
		transfersMap[transfer.ID] = transfer
	}

	// Replace transfers in each sequence with actual transfers
	for _, seq := range sequences {
		for i, transfer := range seq.Transfers {
			if actualTransfer, exists := transfersMap[transfer.ID]; exists {
				seq.Transfers[i] = actualTransfer
			}
		}
	}

	return nil
}

// GetName returns the strategy name
func (s *ClusteredTravelSearchStrategy) GetName() string {
	return "Clustered"
}
