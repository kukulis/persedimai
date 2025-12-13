package generator

import (
	"darbelis.eu/persedimai/dao"
	"darbelis.eu/persedimai/tables"
	"log"
)

type TravelDbConsumer struct {
	travelDao     *dao.TravelDao
	travelsBuffer []*tables.Transfer
	bufferSize    int
	totalCount    int
}

func NewTravelConsumer(travelDao *dao.TravelDao, bufferSize int) *TravelDbConsumer {
	return &TravelDbConsumer{
		travelDao:     travelDao,
		bufferSize:    bufferSize,
		travelsBuffer: []*tables.Transfer{},
	}
}

func (consumer *TravelDbConsumer) Consume(travel *tables.Transfer) error {
	consumer.travelsBuffer = append(consumer.travelsBuffer, travel)
	consumer.totalCount++
	if (consumer.totalCount % 50000) == 0 {
		log.Println("Total count: ", consumer.totalCount)
	}
	if len(consumer.travelsBuffer) >= consumer.bufferSize {
		return consumer.Flush()
	}

	return nil
}

func (consumer *TravelDbConsumer) Flush() error {
	err := consumer.travelDao.InsertMany(consumer.travelsBuffer)
	consumer.travelsBuffer = []*tables.Transfer{}

	return err
}
