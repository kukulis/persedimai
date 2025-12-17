package generator

import "darbelis.eu/persedimai/internal/tables"

type TravelConsumerInterface interface {
	Consume(point *tables.Transfer) error
}
