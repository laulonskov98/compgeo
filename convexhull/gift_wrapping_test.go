package convexhull

import (
	"fmt"
	"testing"
)

func TestSimpleGIFT_CH(t *testing.T) {
	// Test case 1: A basic example with a few points
	points := []Point{
		{0, 3}, {2, 3}, {1, 1}, {2, 1}, {3, 0},
		{0, 0}, {3, 3},
	}
	expected_results := []Point{
		{0, 0}, {0, 3}, {3, 3}, {3, 0},
	}

	// Compute the upper hull
	hull := GIFT_CH(points)

	fmt.Println(hull, expected_results)

	for i, p := range hull {
		if p != expected_results[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", expected_results[i].X, expected_results[i].Y, p.X, p.Y)
		}
	}

	points_x := []Point{
		{0, 0}, // Bottom-left corner
		{2, 2}, // Interior point
		{4, 1}, // Middle point on the right
		{6, 0}, // Bottom-right corner
		{1, 3}, // Upper-left point
		{5, 4}, // Upper-right point
		{3, 5}, // Top-most point
		{3, 1}, // Interior point near the middle
	}

	expected_results_x := []Point{
		{0, 0}, // Starting from the leftmost point
		{1, 3}, // Move upwards to the upper-left
		{3, 5}, // Top-most point
		{5, 4}, // Upper-right point
		{6, 0}, // Rightmost point
	}

	hull_x := GIFT_CH(points_x)
	for i, p := range hull_x {
		if p != expected_results_x[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", expected_results_x[i].X, expected_results_x[i].Y, p.X, p.Y)
		}
	}

	fmt.Println("XXXX")
	fmt.Println(hull_x, expected_results_x)
	fmt.Println("")

	// Test case 2: A set of points forming a triangle
	points2 := []Point{
		{0, 0}, {1, 2}, {2, 0}, {1, 1},
	}
	expected_results2 := []Point{
		{0, 0}, {1, 2}, {2, 0},
	}
	// Compute the upper hull
	hull2 := GIFT_CH(points2)

	fmt.Println(hull2, expected_results2)
	for i, p := range hull2 {
		if p != expected_results2[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", expected_results2[i].X, expected_results2[i].Y, p.X, p.Y)
		}
	}

	// Test case 3: Collinear points
	points3 := []Point{
		{0, 0}, {1, 1}, {2, 2}, {3, 3},
	}
	expected_results3 := []Point{
		{0, 0}, {3, 3},
	}

	hull3 := GIFT_CH(points3)
	fmt.Println(hull3, expected_results3)
	for i, p := range hull3 {
		if p != expected_results3[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", expected_results3[i].X, expected_results3[i].Y, p.X, p.Y)
		}
	}

	// Test case 4: Single point
	points4 := []Point{
		{0, 0},
	}

	hull4 := GIFT_CH(points4)

	for i, p := range hull4 {
		if p != points4[i] {
			t.Errorf("Expected: (%.1f, %.1f), Got: (%.1f, %.1f)", points4[i].X, points4[i].Y, p.X, p.Y)
		}
	}

}
