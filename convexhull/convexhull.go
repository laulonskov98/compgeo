package convexhull

import (
	"math"
	"sort"
)

func INC_CH(points []Point) []Point {
	if len(points) < 3 {
		// Convex hull is not defined for fewer than 3 points.
		return points
	}

	sort.Slice(points, func(i, j int) bool {

		return points[i].X < points[j].X
	})

	// Step 2: Initialize upper hull
	upperHull := []Point{}

	// Step 3: Build upper hull
	for i := 0; i < len(points); i++ {
		// Keep building the upper hull: remove points if they make a left turn (orientation <= 0)
		// but keep clockwise and counterclockwise valid turns
		for len(upperHull) >= 2 && (orientation(upperHull[len(upperHull)-2], upperHull[len(upperHull)-1], points[i]) >= 0) {
			upperHull = upperHull[:len(upperHull)-1] // Remove the last point
		}
		upperHull = append(upperHull, points[i])
	}

	return upperHull
}

func INC_CH_comparison(points []Point) ([]Point, int, int) {
	if len(points) < 3 {
		// Convex hull is not defined for fewer than 3 points.
		return points, 0, 0
	}

	comparisons_sort := 0
	comparison_scan := 0
	sort.Slice(points, func(i, j int) bool {
		comparisons_sort++

		return points[i].X < points[j].X
	})

	// Step 2: Initialize upper hull
	upperHull := []Point{}

	// Step 3: Build upper hull
	for i := 0; i < len(points); i++ {
		comparison_scan++
		// Keep building the upper hull: remove points if they make a left turn (orientation <= 0)
		// but keep clockwise and counterclockwise valid turns
		for len(upperHull) >= 2 && (orientation(upperHull[len(upperHull)-2], upperHull[len(upperHull)-1], points[i]) >= 0) {
			comparison_scan++
			upperHull = upperHull[:len(upperHull)-1] // Remove the last point
		}
		upperHull = append(upperHull, points[i])
	}

	return upperHull, comparison_scan, comparisons_sort
}
func INC_CH_comparison2(points []Point) ([]Point, int, int) {
	if len(points) < 3 {
		// Convex hull is not defined for fewer than 3 points.
		return points, 0, 0
	}

	comparisons_sort := 0
	comparison_scan := 0

	// Sorting the points
	sort.Slice(points, func(i, j int) bool {
		comparisons_sort++
		return points[i].X < points[j].X
	})

	// Initialize upper hull
	upperHull := []Point{}

	// Build upper hull
	for i := 0; i < len(points); i++ {
		// Optionally increment counter for loop control comparison
		// comparison_scan++ // For i < len(points)

		// Inner loop to maintain the upper hull property
		for {
			// First comparison: len(upperHull) >= 2
			comparison_scan++
			if len(upperHull) < 2 {
				break
			}

			// Second comparison: orientation >= 0
			comparison_scan++
			orient := orientationCount(upperHull[len(upperHull)-2], upperHull[len(upperHull)-1], points[i], &comparison_scan)
			if orient < 0 {
				break
			}

			// Remove the last point
			upperHull = upperHull[:len(upperHull)-1]
		}

		// Append the current point
		upperHull = append(upperHull, points[i])
	}

	return upperHull, comparison_scan, comparisons_sort
}

// orientation returns:
// >0 if the sequence of points a->b->c is counter-clockwise,
// <0 if clockwise,
// =0 if colinear.
func orientation(a, b, c Point) float64 {
	result := (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)

	return result
}

func orientationCount(p, q, r Point, counter *int) float64 {
	*counter++
	// Compute the orientation value
	return (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
}

func computeAngle(p1, p2 Point) float64 {
	return math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
}

func SortPoints(points []Point) []Point {
	// sort points by x coordinate
	sort.Slice(points, func(i, j int) bool {
		if points[i].X == points[j].X {
			return points[i].Y > points[j].Y
		}
		return points[i].X < points[j].X
	})
	// step 1.5: make sure the first point is the point with smallest value in both x and y
	min := 0
	for i := 1; i < len(points); i++ {
		if points[i].X < points[min].X {
			min = i
		} else if points[i].X == points[min].X {
			if points[i].Y < points[min].Y {
				min = i
			}
		}
	}
	points[0], points[min] = points[min], points[0]
	return points
}
