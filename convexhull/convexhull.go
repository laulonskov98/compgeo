package convexhull

import "sort"

func INC_CH(points []Point) []Point {
	if len(points) < 3 {
		// Convex hull is not defined for fewer than 3 points.
		return points
	}

	// Step 1: Sort points by X-coordinate (and by Y-coordinate for ties).
	sort.Slice(points, func(i, j int) bool {
		if points[i].X == points[j].X {
			return points[i].Y < points[j].Y
		}
		return points[i].X < points[j].X
	})

	// Step 2: Initialize upper hull
	upperHull := []Point{points[0], points[1]}
	s := 2

	// Step 3: Build upper hull
	for i := 2; i < len(points); i++ {
		for s >= 2 && orientation(upperHull[s-2], upperHull[s-1], points[i]) <= 0 {
			upperHull = upperHull[:s-1] // Remove the last point
			s--
		}
		upperHull = append(upperHull, points[i])
		s++
	}

	return upperHull
}

// orientation returns:
// >0 if the sequence of points a->b->c is counter-clockwise,
// <0 if clockwise,
// =0 if colinear.
func orientation(a, b, c Point) float64 {
	return crossProduct(a, b, c)
}

// crossProduct computes the cross product of vectors AB and AC.
func crossProduct(a, b, c Point) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}
