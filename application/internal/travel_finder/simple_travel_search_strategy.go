package travel_finder

import (
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/tables"
	"darbelis.eu/persedimai/internal/util"
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

// FindPaths finds a sequence of travels from source to destination based on the filter criteria
func (s *SimpleTravelSearchStrategy) FindPaths(filter *data.TravelFilter) ([]*TravelPath, error) {
	var sequences []*tables.TransferSequence
	var err error

	switch filter.TravelCount {
	case 1:
		sequences, err = s.travelDao.FindPathSimple1(filter)
	case 2:
		sequences, err = s.travelDao.FindPathSimple2(filter)
	case 3:
		sequences, err = s.travelDao.FindPathSimple3(filter)
	default:
		if filter.TravelCount > 3 {
			return nil, errors.New("unimplemented: TravelCount > 3 not supported")
		}
		return nil, errors.New("invalid TravelCount: must be 1, 2, or 3")
	}

	if err != nil {
		return nil, err
	}

	if sequences == nil || len(sequences) == 0 {
		return nil, nil
	}

	travelPaths := util.ArrayMap(sequences, func(seq *tables.TransferSequence) *TravelPath {
		return MakeTravelPathOfTransferSequence(seq)
	})

	return travelPaths, nil
}

// GetName returns the strategy name
func (s *SimpleTravelSearchStrategy) GetName() string {
	return "Simple"
}
