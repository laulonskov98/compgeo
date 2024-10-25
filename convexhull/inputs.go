package convexhull

import (
	"math"
	"math/rand"
	"time"
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
		points[i].Y = -(x * x)
	}

	return points
}

func Generate_n_side_polygon(sides, n int) []Point {
	var points []Point
	radius := rand.Float64()
	angleIncrement := 2 * math.Pi / float64(sides)
	for i := 0; i < sides; i++ {
		angle := angleIncrement * float64(i)
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		points = append(points, Point{X: x, Y: y})
	}
	return generateRandomPointsInPolygon(points, n)
}

// Generate random points inside the polygon
func generateRandomPointsInPolygon(vertices []Point, n int) []Point {
	var points []Point
	sides := len(vertices)
	if sides < 3 {
		return points // Not a polygon
	}

	// Calculate the center of the polygon
	var centerX, centerY float64
	for _, v := range vertices {
		centerX += v.X
		centerY += v.Y
	}
	centerX /= float64(sides)
	centerY /= float64(sides)
	center := Point{X: centerX, Y: centerY}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		// Randomly choose a triangle (uniformly)
		j := rand.Intn(sides)
		nextJ := (j + 1) % sides
		triangle := [3]Point{vertices[j], vertices[nextJ], center}

		// Generate a random point inside the triangle
		r1 := rand.Float64()
		r2 := rand.Float64()

		// Adjust r1 and r2 to ensure the point lies inside the triangle
		if r1+r2 > 1 {
			r1 = 1 - r1
			r2 = 1 - r2
		}

		// Compute the point using barycentric coordinates
		x := r1*triangle[0].X + r2*triangle[1].X + (1-r1-r2)*triangle[2].X
		y := r1*triangle[0].Y + r2*triangle[1].Y + (1-r1-r2)*triangle[2].Y

		points = append(points, Point{X: x, Y: y})
	}

	return points
}
