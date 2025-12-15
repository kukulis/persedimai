package tables

import (
	"fmt"
	"math"
)

type Point struct {
	ID   string
	X, Y float64
	Name string
}

func (p1 Point) CalculateDistance(p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (p1 Point) BuildLocationKey() string {
	return fmt.Sprintf("%.5f_%.5f", p1.X, p1.Y)
}

func (p1 Point) BuildYLocationKey() string {
	return fmt.Sprintf("%.5f", p1.Y)
}

func BuildYLocationKey(y float64) string {
	return fmt.Sprintf("%.5f", y)
}
