package data

import "darbelis.eu/persedimai/tables"

type PointGetter interface {
	GetPoint(id string) *tables.Point
}

type MapPointGetter struct {
	points map[string]*tables.Point
}

func (m MapPointGetter) GetPoint(id string) *tables.Point {
	return m.points[id]
}

func NewMapPointGetter(points map[string]*tables.Point) *MapPointGetter {
	return &MapPointGetter{
		points: points,
	}
}
