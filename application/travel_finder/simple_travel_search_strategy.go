package travel_finder

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/tables"
	"errors"
)

// SimpleTravelSearchStrategy implements a basic travel search strategy
type SimpleTravelSearchStrategy struct {
	travelDao *dao.TravelDao
}

// NewSimpleTravelSearchStrategy creates a new SimpleTravelSearchStrategy
func NewSimpleTravelSearchStrategy(travelDao *dao.TravelDao) *SimpleTravelSearchStrategy {
	return &SimpleTravelSearchStrategy{
		travelDao: travelDao,
	}
}

// FindPath finds a sequence of travels from source to destination based on the filter criteria
func (s *SimpleTravelSearchStrategy) FindPath(filter *data.TravelFilter) (*TravelPath, error) {
	var sequences []*tables.TransferSequence
	var err error

	switch filter.TravelCount {
	case 1:
		sequences, err = s.travelDao.FindPathSimple1(filter)
	case 2:
		sequences, err = s.travelDao.FindPathSimple2(filter)
	default:
		if filter.TravelCount > 2 {
			return nil, errors.New("unimplemented: TravelCount > 2 not supported")
		}
		return nil, errors.New("invalid TravelCount: must be 1 or 2")
	}

	if err != nil {
		return nil, err
	}

	if sequences == nil || len(sequences) == 0 {
		return nil, nil
	}

	// Select the best sequence (earliest arrival, already sorted by SQL)
	bestSequence := sequences[0]

	// Build TravelPath from the best sequence
	travelPath := &TravelPath{
		Travels:       bestSequence.Transfers,
		TransferCount: bestSequence.TransferCount() - 1,
		TotalDuration: bestSequence.TotalDuration(),
	}

	return travelPath, nil
}

// GetName returns the strategy name
func (s *SimpleTravelSearchStrategy) GetName() string {
	return "Simple"
}
