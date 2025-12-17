package generator

import "darbelis.eu/persedimai/internal/tables"

type PointConsumerInterface interface {
	Consume(point *tables.Point) error
}
