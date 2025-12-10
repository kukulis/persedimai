package generator

import (
	"darbelis.eu/persedimai/tables"
	"fmt"
	"testing"
	"time"
)

func TestGenerateTravels2Points(t *testing.T) {
	gf := &GeneratorFactory{}
	g := gf.createGenerator(5, 1000, 0, &SimpleIdGenerator{})

	pointA := tables.Point{ID: "1", X: 0, Y: 0}
	pointB := tables.Point{ID: "2", X: 10000, Y: 0}
	from := GetTime("2027-01-01 00:00:00")
	till := GetTime("2027-01-02 00:00:00")
	g.GenerateTravelsForTwoPoints(pointA, pointB, from, till, 1000, 2)

	expectedTravels := []*tables.Travel{
		{From: 1, To: 2, Departure: GetTime("2027-01-01 00:00:00"), Arrival: GetTime("2027-01-01 10:00:00")},
		{From: 2, To: 1, Departure: GetTime("2027-01-01 12:00:00"), Arrival: GetTime("2027-01-01 22:00:00")},
	}

	if len(expectedTravels) != len(g.Travels()) {
		t.Errorf("expectedTravels size: %d, actualTravels size: %d", len(expectedTravels), len(g.Travels()))
	}
}

func GetTime(timeStr string) time.Time {
	t, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		_ = fmt.Errorf("error parsing time: %v", err)
	}
	return t
}
