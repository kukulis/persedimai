package generator

import "darbelis.eu/persedimai/tables"

type TravelConsumerInterface interface {
	Consume(point *tables.Travel) error
}
