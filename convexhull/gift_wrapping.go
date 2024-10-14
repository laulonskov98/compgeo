package convexhull

// GiftWrappingUpperHull computes the upper hull of a set of points using the Gift Wrapping (Jarvis March) algorithm.
func GIFT_CH(points []Point) []Point {
	if len(points) < 2 {
		// Not enough points to form a hull.
		return points
	}

	// Step 1: Find the leftmost point as the starting point.
	start := 0
	for i := 1; i < len(points); i++ {
		if points[i].X < points[start].X || (points[i].X == points[start].X && points[i].Y < points[start].Y) {
			start = i
		}
	}

	// Initialize the hull with the starting point.
	hull := []Point{points[start]}
	p := start

	// Step 2: Iteratively find the next point that makes the smallest angle (most counter-clockwise turn).
	for {
		next := -1
		for i := 0; i < len(points); i++ {
			if i == p {
				continue
			}
			if next == -1 || orientation(points[p], points[next], points[i]) > 0 || (orientation(points[p], points[next], points[i]) == 0 && points[i].X > points[next].X) {
				next = i
			}
		}
		// If we are going to the left, stop since we're only interested in the upper hull.
		if points[next].X < points[p].X {
			break
		}
		hull = append(hull, points[next])
		p = next

		// Break if we've wrapped around back to the starting point or if we reach the rightmost point.
		if points[p].X > points[start].X && points[p].Y > points[start].Y {
			break
		}
	}

	return hull
}
