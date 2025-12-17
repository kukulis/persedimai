package tables

import "testing"

func TestPoint(t *testing.T) {
	t.Run("BuildLocationKey", func(t *testing.T) {
		point := Point{X: 123.456789, Y: 987.654321}
		key := point.BuildLocationKey()
		expected := "123.45679_987.65432"

		if key != expected {
			t.Errorf("BuildLocationKey() = %s, expected %s", key, expected)
		}
	})

	t.Run("BuildLocationKeyWithIntegers", func(t *testing.T) {
		point := Point{X: 100.0, Y: 200.0}
		key := point.BuildLocationKey()
		expected := "100.00000_200.00000"

		if key != expected {
			t.Errorf("BuildLocationKey() = %s, expected %s", key, expected)
		}
	})

	t.Run("BuildLocationKeyWithNegatives", func(t *testing.T) {
		point := Point{X: -123.456, Y: -987.654}
		key := point.BuildLocationKey()
		expected := "-123.45600_-987.65400"

		if key != expected {
			t.Errorf("BuildLocationKey() = %s, expected %s", key, expected)
		}
	})

	t.Run("BuildLocationKeyWithZero", func(t *testing.T) {
		point := Point{X: 0.0, Y: 0.0}
		key := point.BuildLocationKey()
		expected := "0.00000_0.00000"

		if key != expected {
			t.Errorf("BuildLocationKey() = %s, expected %s", key, expected)
		}
	})

	t.Run("CalculateDistance", func(t *testing.T) {
		p1 := Point{X: 0, Y: 0}
		p2 := Point{X: 3, Y: 4}
		distance := p1.CalculateDistance(p2)
		expected := 5.0

		if distance != expected {
			t.Errorf("CalculateDistance() = %f, expected %f", distance, expected)
		}
	})

	t.Run("CalculateDistanceZero", func(t *testing.T) {
		p1 := Point{X: 10, Y: 20}
		p2 := Point{X: 10, Y: 20}
		distance := p1.CalculateDistance(p2)
		expected := 0.0

		if distance != expected {
			t.Errorf("CalculateDistance() = %f, expected %f", distance, expected)
		}
	})

	t.Run("CalculateDistanceDiagonal", func(t *testing.T) {
		p1 := Point{X: 0, Y: 0}
		p2 := Point{X: 3000, Y: 4000}
		distance := p1.CalculateDistance(p2)
		expected := 5000.0

		if distance != expected {
			t.Errorf("CalculateDistance() = %f, expected %f", distance, expected)
		}
	})
}
