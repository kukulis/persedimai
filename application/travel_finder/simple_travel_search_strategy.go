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
	var transfers []*tables.Transfer
	var err error

	switch filter.TravelCount {
	case 1:
		transfers, err = s.travelDao.FindPathSimple1(filter)
	case 2:
		transfers, err = s.travelDao.FindPathSimple2(filter)
	default:
		if filter.TravelCount > 2 {
			return nil, errors.New("unimplemented: TravelCount > 2 not supported")
		}
		return nil, errors.New("invalid TravelCount: must be 1 or 2")
	}

	if err != nil {
		return nil, err
	}

	if transfers == nil || len(transfers) == 0 {
		return nil, nil
	}

	// Build TravelPath from transfers
	travelPath := &TravelPath{
		Travels:       transfers,
		TransferCount: len(transfers) - 1,
	}

	// Calculate total duration and distance
	if len(transfers) > 0 {
		firstTransfer := transfers[0]
		lastTransfer := transfers[len(transfers)-1]
		travelPath.TotalDuration = lastTransfer.Arrival.Sub(firstTransfer.Departure)
	}

	return travelPath, nil
}

// GetName returns the strategy name
func (s *SimpleTravelSearchStrategy) GetName() string {
	return "Simple"
}
