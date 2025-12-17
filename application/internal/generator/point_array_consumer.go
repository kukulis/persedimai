package generator

import "darbelis.eu/persedimai/internal/tables"

type PointArrayConsumer struct {
	Points []*tables.Point
}

func NewPointArrayConsumer() *PointArrayConsumer {
	return &PointArrayConsumer{Points: []*tables.Point{}}
}

func (consumer *PointArrayConsumer) Consume(point *tables.Point) error {
	consumer.Points = append(consumer.Points, point)

	return nil
}
