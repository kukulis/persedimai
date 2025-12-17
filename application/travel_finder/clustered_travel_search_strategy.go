package travel_finder

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/data"
	"darbelis.eu/persedimai/tables"
	"darbelis.eu/persedimai/util"
	"errors"
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

	switch filter.TravelCount {
	case 1:
		sequences, err = s.travelDao.FindPathSimple1(filter)
	case 2:
		sequences, err = s.travelDao.FindPathClustered2(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo)
	case 3:
		sequences, err = s.travelDao.FindPathClustered3(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo)
	case 4:
		sequences, err = s.travelDao.FindPathClustered4(filter.Source, filter.Destination, filter.ArrivalTimeFrom, filter.ArrivalTimeTo)
	default:
		if filter.TravelCount > 4 {
			return nil, errors.New("unimplemented: TravelCount > 4 not supported")
		}
		return nil, errors.New("invalid TravelCount: must be 2, 3, or 4")
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
func (s *ClusteredTravelSearchStrategy) GetName() string {
	return "Clustered"
}
