package generator

import "darbelis.eu/persedimai/tables"

type TravelArrayConsumer struct {
	Travels []*tables.Travel
}

func NewTravelArrayConsumer() *TravelArrayConsumer {
	return &TravelArrayConsumer{Travels: []*tables.Travel{}}
}

func (consumer *TravelArrayConsumer) Consume(travel *tables.Travel) error {
	consumer.Travels = append(consumer.Travels, travel)

	return nil
}
