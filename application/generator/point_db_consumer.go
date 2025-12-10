package generator

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/tables"
)

type PointDbConsumer struct {
	pointDao     *dao.PointDao
	pointsBuffer []*tables.Point
	bufferSize   int
}

func NewPointConsumer(pointDao *dao.PointDao, bufferSize int) *PointDbConsumer {
	return &PointDbConsumer{
		pointDao:     pointDao,
		bufferSize:   bufferSize,
		pointsBuffer: []*tables.Point{},
	}
}

func (consumer *PointDbConsumer) Consume(point *tables.Point) error {
	consumer.pointsBuffer = append(consumer.pointsBuffer, point)
	if len(consumer.pointsBuffer) >= consumer.bufferSize {
		return consumer.Flush()
	}

	return nil
}

func (consumer *PointDbConsumer) Flush() error {
	err := consumer.pointDao.InsertMany(consumer.pointsBuffer)
	consumer.pointsBuffer = []*tables.Point{}

	return err
}
