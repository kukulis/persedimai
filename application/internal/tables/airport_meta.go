package tables

import "time"

type AirportMeta struct {
	AirportCode  string
	ImportedFrom *time.Time
	ImportedTo   *time.Time
}
