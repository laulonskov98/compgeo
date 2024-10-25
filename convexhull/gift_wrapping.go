package convexhull

import (
	"math"
)

// GiftWrappingUpperHull computes the upper hull of a set of points using the Gift Wrapping (Jarvis March) algorithm.
func GIFT_CH(points []Point) []Point {
	if len(points) < 2 {
		return points // Not enough points to form a hull.
	}

	comparisons_count := 0
	// Step 1: Find the leftmost point as the starting point.
	start := 0
	for i := 1; i < len(points); i++ {
		comparisons_count++
		if points[i].X < points[start].X || (points[i].X == points[start].X && points[i].Y < points[start].Y) {
			start = i
		}
	}

	// Initialize the hull with the starting point.
	hull := []Point{points[start]}
	existing_points := make(map[Point]bool)
	p := start
	last_round_orient := math.MaxFloat64
	// Step 2: Iteratively find the next point that makes the most counter-clockwise turn.
	for {
		next := -1
		final_orient := math.MaxFloat64
		for i := 0; i < len(points); i++ {
			if i == p {
				continue
			}
			comparisons_count++
			if next == -1 {
				next = i // First candidate
			} else {
				comparisons_count++
				orient := orientation(points[p], points[next], points[i])
				if orient > 0 || (orient == 0 && points[i].X > points[next].X) {
					next = i
					final_orient = orient
				}
			}
		}
		// Stop if we've come back to the starting point
		if points[next].X < hull[len(hull)-1].X || existing_points[points[next]] || (last_round_orient == 0 && final_orient == 0) {
			break
		}

		// Append the next point to the hull
		hull = append(hull, points[next])
		existing_points[points[next]] = true
		p = next
		last_round_orient = final_orient
	}

	return hull
}

// GiftWrappingUpperHull computes the upper hull of a set of points using the Gift Wrapping (Jarvis March) algorithm.
func GIFT_CH_comparison(points []Point) ([]Point, int) {
	if len(points) < 2 {
		return points, 0 // Not enough points to form a hull.
	}

	comparisons_count := 0
	// Step 1: Find the leftmost point as the starting point.
	start := 0
	for i := 1; i < len(points); i++ {
		comparisons_count++
		if points[i].X < points[start].X || (points[i].X == points[start].X && points[i].Y < points[start].Y) {
			start = i
		}
	}

	// Initialize the hull with the starting point.
	hull := []Point{points[start]}
	existing_points := make(map[Point]bool)
	p := start
	last_round_orient := math.MaxFloat64
	// Step 2: Iteratively find the next point that makes the most counter-clockwise turn.
	for {
		next := -1
		final_orient := math.MaxFloat64
		for i := 0; i < len(points); i++ {
			if i == p {
				continue
			}
			comparisons_count++
			if next == -1 {
				next = i // First candidate
			} else {
				comparisons_count++
				orient := orientation(points[p], points[next], points[i])
				if orient > 0 || (orient == 0 && points[i].X > points[next].X) {
					next = i
					final_orient = orient
				}
			}
		}
		// Stop if we've come back to the starting point
		if points[next].X < hull[len(hull)-1].X || existing_points[points[next]] || (last_round_orient == 0 && final_orient == 0) {
			break
		}

		// Append the next point to the hull
		hull = append(hull, points[next])
		existing_points[points[next]] = true
		p = next
		last_round_orient = final_orient
	}

	return hull, comparisons_count
}
