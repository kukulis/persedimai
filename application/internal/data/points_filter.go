package data

type PointsFilter struct {
	Limit    int
	X        *float64
	Y        *float64
	NamePart string
	IdPart   string
}

func NewPointsFilter() *PointsFilter {
	return &PointsFilter{
		Limit: 10, // default limit
	}
}
