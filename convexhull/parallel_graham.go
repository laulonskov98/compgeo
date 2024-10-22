package convexhull

import (
	"fmt"
	"sync"
)

// Compute the upper hull for a set of points using Graham's Scan (in parallel).
func PAR_GS(points []Point, p int) []Point {
	n := len(points)
	if n < 3 {
		return points // No convex hull possible with less than 3 points.
	}

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

	// Collect the upper hulls.
	var upperHulls [][]Point
	for hull := range hullCh {
		upperHulls = append(upperHulls, hull)
	}

	// Step 3: Merge the upper hulls using tangents
	mergedHull := upperHulls[0]
	for i := 1; i < len(upperHulls); i++ {
		mergedHull = mergeHulls(mergedHull, upperHulls[i])
	}

	return mergedHull
}

// mergeHulls merges two upper hulls by finding the tangent and connecting the two hulls.
func mergeHulls(hull1, hull2 []Point) []Point {
	// We need to find the left tangent of hull1 and the right tangent of hull2.
	// We will use a simple linear scan to find these tangents.
	i, j := len(hull1)-1, 0
	fmt.Println(hull1, hull2)
	// Find the tangents
	for {
		changed := false

		// Check if we need to adjust the i index (left tangent of hull1)
		if i > 0 && orientation(hull1[i-1], hull1[i], hull2[j]) <= 0 {
			i--
			changed = true
		}

		// Check if we need to adjust the j index (right tangent of hull2)
		if j < len(hull2)-1 && orientation(hull1[i], hull2[j], hull2[j+1]) <= 0 {
			j++
			changed = true
		}

		// If no changes were made, we have found the tangents.
		if !changed {
			break
		}
	}

	// Combine hull1 and hull2 along the tangent.
	mergedHull := append(hull1[:i+1], hull2[j:]...)
	return mergedHull
}
