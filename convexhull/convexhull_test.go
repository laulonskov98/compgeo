package convexhull

import (
	"testing"
)

func TestGrahamScan(t *testing.T) {
	// Test case 1: A basic example with a few points
	points := []Point{
		{0, 3}, {2, 3}, {1, 1}, {2, 1}, {3, 0},
		{0, 0}, {3, 3},
	}
	expected_results := []Point{
		{0, 0}, {3, 0}, {3, 3},
	}

	// Compute the upper hull
	hull := INC_CH(points)

	for i, p := range hull {
		if p != expected_results[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", expected_results[i].X, expected_results[i].Y, p.X, p.Y)
		}
	}

	// Test case 2: A set of points forming a triangle
	points2 := []Point{
		{0, 0}, {1, 2}, {2, 0}, {1, 1},
	}
	expected_results2 := []Point{
		{0, 0}, {2, 0}, {1, 2},
	}
	// Compute the upper hull
	hull2 := INC_CH(points2)

	for i, p := range hull2 {
		if p != expected_results2[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", points2[i].X, points2[i].Y, p.X, p.Y)
		}
	}

	// Test case 3: Collinear points
	points3 := []Point{
		{0, 0}, {1, 1}, {2, 2}, {3, 3},
	}
	expected_results3 := []Point{
		{0, 0}, {3, 3},
	}

	hull3 := INC_CH(points3)

	for i, p := range hull3 {
		if p != expected_results3[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", points3[i].X, points3[i].Y, p.X, p.Y)
		}
	}

	// Test case 4: Single point
	points4 := []Point{
		{0, 0},
	}

	hull4 := INC_CH(points4)

	for i, p := range hull4 {
		if p != points4[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", points4[i].X, points4[i].Y, p.X, p.Y)
		}
	}

}
