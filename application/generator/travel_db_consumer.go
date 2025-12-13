package generator

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/tables"
)

type TravelDbConsumer struct {
	travelDao     *dao.TravelDao
	travelsBuffer []*tables.Travel
	bufferSize    int
}

func NewTravelConsumer(travelDao *dao.TravelDao, bufferSize int) *TravelDbConsumer {
	return &TravelDbConsumer{
		travelDao:     travelDao,
		bufferSize:    bufferSize,
		travelsBuffer: []*tables.Travel{},
	}
}

func (consumer *TravelDbConsumer) Consume(travel *tables.Travel) error {
	consumer.travelsBuffer = append(consumer.travelsBuffer, travel)
	if len(consumer.travelsBuffer) >= consumer.bufferSize {
		return consumer.Flush()
	}

	return nil
}

func (consumer *TravelDbConsumer) Flush() error {
	err := consumer.travelDao.InsertMany(consumer.travelsBuffer)
	consumer.travelsBuffer = []*tables.Travel{}

	return err
}
