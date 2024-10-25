package convexhull

import (
	"fmt"
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

	SortPoints(points)

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

	min_point1_idx, min_point2_idx, min_point_j := -1, -1, Point{}
	i := 0
	for i < p {
		min_rotation := -math.MaxFloat64
		min_index := -1
		for j := i + 1; j < p; j++ {
			// Find the tangent with the smallest rotation
			point1_idx, point2_idx := findUpperTangentBinary(upperHulls[i], upperHulls[j])
			point1, point2 := upperHulls[i][point1_idx], upperHulls[j][point2_idx]
			angle := computeAngle(point1, point2)

			if angle >= min_rotation {
				min_rotation = angle
				min_index = j
				min_point1_idx = point1_idx
				min_point2_idx = point2_idx
				min_point_j = point2
			}
		}

		if min_index == -1 {
			mergedHull = append(mergedHull, upperHulls[i][min_point2_idx+1:]...)
			break
		}

		if i == 0 {
			mergedHull = append(mergedHull, upperHulls[i][:min_point1_idx+1]...)
		}

		min_point_incrementer := 0
		points_to_append := make([]Point, 0)

		// Collect points in reverse order
		for upperHulls[i][min_point1_idx-min_point_incrementer] != mergedHull[len(mergedHull)-1] &&
			upperHulls[i][min_point1_idx-min_point_incrementer].X < min_point_j.X {
			points_to_append = append(points_to_append, upperHulls[i][min_point1_idx-min_point_incrementer])
			min_point_incrementer++
		}

		// Reverse the collected points
		for j := len(points_to_append) - 1; j >= 0; j-- {
			mergedHull = append(mergedHull, points_to_append[j])
		}

		mergedHull = append(mergedHull, min_point_j)
		i = min_index
	}

	return mergedHull
}

type HullResult struct {
	Hull []Point // The upper hull result
	Scan int     // Number of comparisons during the scan
	Sort int     // Number of comparisons during the sort
}

func PAR_GS_comparison(points []Point, p int) ([]Point, int) {
	n := len(points)
	if n < 3 {
		return points, 0 // No convex hull possible with less than 3 points.
	}

	SortPoints(points)

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
	hullCh := make(chan HullResult, p)
	for _, subset := range subsets {
		wg.Add(1)
		go func(subset []Point) {
			defer wg.Done()
			points, scan, sort := INC_CH_comparison(subset)
			hullCh <- HullResult{points, scan, sort}
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
	highestcompare_inch := 0
	for hull := range hullCh {
		upperHulls = append(upperHulls, hull.Hull)
		if hull.Scan+hull.Sort > highestcompare_inch {
			highestcompare_inch = hull.Scan + hull.Sort
		}
	}

	fmt.Println("Highest number of comparisons in a single hull:", highestcompare_inch)

	sortRowsByFirstPointX(upperHulls)

	// Step 3: Merge the upper hulls using tangents
	mergedHull := make([]Point, 0, n)

	min_point1_idx, min_point2_idx, min_point_j := -1, -1, Point{}
	i := 0
	for i < p {
		min_rotation := -math.MaxFloat64
		min_index := -1
		for j := i + 1; j < p; j++ {
			highestcompare_inch++
			// Find the tangent with the smallest rotation
			point1_idx, point2_idx, binarycomps := findUpperTangentBinaryComp(upperHulls[i], upperHulls[j])
			highestcompare_inch += binarycomps

			point1, point2 := upperHulls[i][point1_idx], upperHulls[j][point2_idx]
			angle := computeAngle(point1, point2)

			if angle >= min_rotation {
				highestcompare_inch++
				min_rotation = angle
				min_index = j
				min_point1_idx = point1_idx
				min_point2_idx = point2_idx
				min_point_j = point2
			}
		}

		if min_index == -1 {
			mergedHull = append(mergedHull, upperHulls[i][min_point2_idx+1:]...)
			break
		}

		if i == 0 {
			mergedHull = append(mergedHull, upperHulls[i][:min_point1_idx+1]...)
		}

		min_point_incrementer := 0
		points_to_append := make([]Point, 0)

		// Collect points in reverse order
		for upperHulls[i][min_point1_idx-min_point_incrementer] != mergedHull[len(mergedHull)-1] &&
			upperHulls[i][min_point1_idx-min_point_incrementer].X < min_point_j.X {
			highestcompare_inch++
			points_to_append = append(points_to_append, upperHulls[i][min_point1_idx-min_point_incrementer])
			min_point_incrementer++
		}

		// Reverse the collected points
		for j := len(points_to_append) - 1; j >= 0; j-- {
			mergedHull = append(mergedHull, points_to_append[j])
		}

		mergedHull = append(mergedHull, min_point_j)
		i = min_index
	}

	return mergedHull, highestcompare_inch
}

// findUpperTangent finds the upper tangent between two convex hulls Ui and Uj.
// Each hull is represented as a slice of Points, ordered from left to right.
func findUpperTangent(Ui, Uj []Point) (int, int) {
	i := len(Ui) - 1 // Start with the rightmost point of Ui
	j := 0           // Start with the leftmost point of Uj

	changed := true
	for changed {
		changed = false
		// Adjust Ui
		for {
			nextI := (i - 1 + len(Ui)) % len(Ui)
			orient := orientation(Ui[i], Uj[j], Ui[nextI])
			if orient > 0 || (orient == 0 && Ui[i].X > Uj[j].X) {
				i = nextI
				changed = true
			} else {
				break
			}
		}
		// Adjust U
		for {
			nextJ := (j + 1) % len(Uj)
			orient := orientation(Ui[i], Uj[j], Uj[nextJ])

			if orient > 0 || (orientation(Ui[i], Uj[nextJ], Uj[j]) == 0 && Uj[j].X > Uj[nextJ].X) {
				j = nextJ
				changed = true
			} else {
				break
			}
		}
	}
	// The tangent is between Ui[i] and Uj[j]
	return i, j
}

// findTangentOnUi finds the index i on Ui such that the line from Ui[i] to Uj_j
// is the upper tangent to Ui.
func findTangentOnUi(Ui []Point, Uj_j Point) int {
	low := 0
	high := len(Ui) - 1

	for {
		mid := (low + high) / 2
		n := len(Ui)

		midPrev := mid - 1
		if midPrev < 0 {
			midPrev = 0
		}
		midNext := mid + 1
		if midNext >= n {
			midNext = n - 1
		}

		orientPrev := orientation(Ui[mid], Uj_j, Ui[midPrev])
		orientNext := orientation(Ui[mid], Uj_j, Ui[midNext])

		if orientPrev > 0 {
			// Left neighbor is above the line, move left
			high = mid - 1
		} else if orientNext > 0 {
			// Right neighbor is above the line, move right
			low = mid + 1
		} else {
			// Found the tangent point
			return mid
		}

		if low > high {
			// Converged to the best candidate
			return mid
		}
	}
}

// findTangentOnUj finds the index j on Uj such that the line from Ui_i to Uj[j]
// is the upper tangent to Uj.
func findTangentOnUj(Uj []Point, Ui_i Point) int {
	low := 0
	high := len(Uj) - 1

	for {
		mid := (low + high) / 2
		n := len(Uj)

		midPrev := mid - 1
		if midPrev < 0 {
			midPrev = 0
		}
		midNext := mid + 1
		if midNext >= n {
			midNext = n - 1
		}

		orientPrev := orientation(Uj[mid], Ui_i, Uj[midPrev])
		orientNext := orientation(Uj[mid], Ui_i, Uj[midNext])

		if orientPrev < 0 {
			// Left neighbor is above the line (since we invert the orientation), move left
			high = mid - 1
		} else if orientNext < 0 {
			// Right neighbor is above the line, move right
			low = mid + 1
		} else {
			// Found the tangent point
			return mid
		}

		if low > high {
			// Converged to the best candidate
			return mid
		}
	}
}

// findUpperTangent finds the upper tangent between two convex hulls Ui and Uj.
// Each hull is represented as a slice of Points, ordered from left to right.
func findUpperTangentBinary(Ui, Uj []Point) (int, int) {
	// Initial indices
	i := len(Ui) - 1 // Start with the rightmost point of Ui
	j := 0           // Start with the leftmost point of Uj

	for {
		prevI := i
		prevJ := j

		// Find the tangent on Uj for Ui[i]
		j = findTangentOnUj(Uj, Ui[i])

		// Find the tangent on Ui for Uj[j]
		i = findTangentOnUi(Ui, Uj[j])

		// Check if indices have stabilized
		if i == prevI && j == prevJ {
			break
		}
	}

	// The tangent is between Ui[i] and Uj[j]
	return i, j
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

func findTangentOnUiComp(Ui []Point, Uj_j Point) (int, int) {
	low := 0
	high := len(Ui) - 1
	compares := 0

	for {
		mid := (low + high) / 2
		n := len(Ui)

		midPrev := mid - 1
		if midPrev < 0 {
			midPrev = 0
		}
		midNext := mid + 1
		if midNext >= n {
			midNext = n - 1
		}

		orientPrev := orientation(Ui[mid], Uj_j, Ui[midPrev])
		orientNext := orientation(Ui[mid], Uj_j, Ui[midNext])
		compares += 1

		if orientPrev > 0 {
			// Left neighbor is above the line, move left
			high = mid - 1
		} else if orientNext > 0 {
			// Right neighbor is above the line, move right
			low = mid + 1
		} else {
			// Found the tangent point
			return mid, compares
		}

		if low > high {
			// Converged to the best candidate
			return mid, compares
		}
	}
}

// findTangentOnUj finds the index j on Uj such that the line from Ui_i to Uj[j]
// is the upper tangent to Uj.
func findTangentOnUjComp(Uj []Point, Ui_i Point) (int, int) {
	compares := 0
	low := 0
	high := len(Uj) - 1

	for {
		mid := (low + high) / 2
		n := len(Uj)

		midPrev := mid - 1
		if midPrev < 0 {
			midPrev = 0
		}
		midNext := mid + 1
		if midNext >= n {
			midNext = n - 1
		}
		compares += 1

		orientPrev := orientation(Uj[mid], Ui_i, Uj[midPrev])
		orientNext := orientation(Uj[mid], Ui_i, Uj[midNext])

		if orientPrev < 0 {
			// Left neighbor is above the line (since we invert the orientation), move left
			high = mid - 1
		} else if orientNext < 0 {
			// Right neighbor is above the line, move right
			low = mid + 1
		} else {
			// Found the tangent point
			return mid, compares
		}

		if low > high {
			// Converged to the best candidate
			return mid, compares
		}
	}
}

// findUpperTangent finds the upper tangent between two convex hulls Ui and Uj.
// Each hull is represented as a slice of Points, ordered from left to right.
func findUpperTangentBinaryComp(Ui, Uj []Point) (int, int, int) {
	compares := 0
	// Initial indices
	i := len(Ui) - 1 // Start with the rightmost point of Ui
	j := 0           // Start with the leftmost point of Uj
	icompare, jcompare := 0, 0
	for {
		prevI := i
		prevJ := j

		// Find the tangent on Uj for Ui[i]
		j, jcompare = findTangentOnUjComp(Uj, Ui[i])

		// Find the tangent on Ui for Uj[j]
		i, icompare = findTangentOnUiComp(Ui, Uj[j])
		compares += jcompare + icompare

		// Check if indices have stabilized
		if i == prevI && j == prevJ {
			break
		}
	}

	// The tangent is between Ui[i] and Uj[j]
	return i, j, compares
}
