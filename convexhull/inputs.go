package convexhull

import (
	"math"
	"math/rand"
)

// Point represents a point in 2D space
type Point struct {
	X float64
	Y float64
}

func Generate_square_inputs(n int) []Point {
	// this function generates n uniformly distributed points in a square
	// with corners at (0,0) and (1,1)
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i].X = rand.Float64()
		points[i].Y = rand.Float64()
	}

	return points
}

func Generate_circle_inputs(n int) []Point {
	// this function generates n uniformly distributed points in a circle
	// with radius 1 and center at (0,0)
	points := make([]Point, n)
	for i := 0; i < n; i++ {

		// generate random angle and radius
		theta := rand.Float64() * 2 * math.Pi
		r := rand.Float64()

		points[i].X = r * math.Cos(theta)
		points[i].Y = r * math.Sin(theta)
	}

	return points
}

// This functions generates input that are on the curve Y = -X^2
func Generate_polynomial_inputs(n int) []Point {
	// this function generates n points on the curve y = -x^2
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		x := rand.Float64()
		points[i].X = x
		points[i].Y = -x * x
	}

	return points
}
