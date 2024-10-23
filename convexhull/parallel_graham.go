package convexhull

import (
	"math"
	"sort"
	"sync"
)

// Compute the upper hull for a set of points using Graham's Scan (in parallel).
func PAR_GS(points []Point, p int) []Point {
	n := len(points)
	if n < 3 {
		return points // No convex hull possible with less than 3 points.
	}

	// sort points by x coordinate
	sort.Slice(points, func(i, j int) bool {
		if points[i].X == points[j].X {
			return points[i].Y < points[j].Y
		}
		return points[i].X < points[j].X
	})

	// Step 1: Split S into p equal-sized arrays, S1, S2, ..., Sp
	chunkSize := (n + p - 1) / p // Calculate the size of each chunk
	subsets := make([][]Point, 0, p)

	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		subsets = append(subsets, points[i:end])
	}

	// Step 2: Compute the upper hull of each subset in parallel
	var wg sync.WaitGroup
	hullCh := make(chan []Point, p)
	for _, subset := range subsets {
		wg.Add(1)
		go func(subset []Point) {
			defer wg.Done()
			hullCh <- INC_CH(subset)
		}(subset)
	}

	// Wait for all goroutines to finish and close the hullCh.
	go func() {
		wg.Wait()
		close(hullCh)
	}()
	// sort the upper hulls by the x coordinate of the last point

	// Collect the upper hulls.
	var upperHulls [][]Point
	for hull := range hullCh {
		upperHulls = append(upperHulls, hull)
	}

	sortRowsByFirstPointX(upperHulls)

	// Step 3: Merge the upper hulls using tangents
	mergedHull := make([]Point, 0, n)
	i := 0
	for i < p {
		min_rotation := math.MaxFloat32
		min_index := -1
		min_point_i, min_point_j := Point{X: math.MaxFloat32, Y: math.MaxFloat32}, Point{X: math.MaxFloat32, Y: math.MaxFloat32}
		for j := i + 1; j < p; j++ {
			// Find the tangent with the smallest rotation
			point1, point2 := findUpperTangent(upperHulls[i], upperHulls[j])
			angle := computeAngle(point1, point2)
			if angle <= min_rotation {
				min_rotation = angle
				min_index = j
				min_point_i = point1
				min_point_j = point2
			}
		}

		if min_index == -1 {
			break
		}

		if i == 0 {
			mergedHull = append(mergedHull, min_point_i)
		}
		mergedHull = append(mergedHull, min_point_j)
		i = min_index
	}

	return mergedHull
}

// findUpperTangent finds the upper tangent between two convex hulls Ui and Uj.
// Each hull is represented as a slice of Points, ordered from left to right.
func findUpperTangent(Ui, Uj []Point) (Point, Point) {
	i := len(Ui) - 1 // Start with the rightmost point of Ui
	j := 0           // Start with the leftmost point of Uj

	changed := true
	for changed {
		changed = false
		// Adjust Ui
		for {
			nextI := (i - 1 + len(Ui)) % len(Ui)
			if orientation(Ui[nextI], Uj[j], Ui[i]) > 0 || (orientation(Ui[nextI], Uj[j], Ui[i]) == 0 && Ui[i].X > Uj[j].X) {
				i = nextI
				changed = true
			} else {
				break
			}
		}
		// Adjust U
		for {
			nextJ := (j + 1) % len(Uj)

			if orientation(Ui[i], Uj[nextJ], Uj[j]) > 0 || (orientation(Ui[i], Uj[nextJ], Uj[j]) == 0 && Uj[j].X > Uj[nextJ].X) {
				j = nextJ
				changed = true
			} else {
				break
			}
		}
	}
	// The tangent is between Ui[i] and Uj[j]
	return Ui[i], Uj[j]
}

func sortRowsByFirstPointX(points [][]Point) {
	sort.Slice(points, func(i, j int) bool {
		// Check if the rows are non-empty
		if len(points[i]) == 0 || len(points[j]) == 0 {
			return false
		}
		return points[i][len(points[i])-1].X < points[j][len(points[j])-1].X
	})
}
