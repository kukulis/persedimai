package generator

import "darbelis.eu/persedimai/tables"

type PointConsumerInterface interface {
	Consume(point *tables.Point) error
}
