package generator

import "darbelis.eu/persedimai/tables"

type TravelArrayConsumer struct {
	Travels []*tables.Transfer
}

func NewTravelArrayConsumer() *TravelArrayConsumer {
	return &TravelArrayConsumer{Travels: []*tables.Transfer{}}
}

func (consumer *TravelArrayConsumer) Consume(travel *tables.Transfer) error {
	consumer.Travels = append(consumer.Travels, travel)

	return nil
}
