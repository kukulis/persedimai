package tables

import "math"

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
