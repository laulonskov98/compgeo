package sortpkg

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

// TestSequentialMerge tests the sequential merge function
func TestSequentialMerge(t *testing.T) {
	// get n and m random numbers
	n := 10
	m := 10
	left := make([]int, n)
	right := make([]int, m)

	for i := 0; i < n; i++ {
		left[i] = i
	}
	for i := 0; i < m; i++ {
		right[i] = i + n
	}

	// merge the two arrays
	merged := make([]int, n+m)
	merged = sequential_merge(left, right, merged)

	// test if all numbers from left and right occur in merged
	for i := 0; i < n; i++ {
		found := false
		for j := 0; j < n+m; j++ {
			if left[i] == merged[j] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Number %d from left array not found in merged array", left[i])
		}
	}
	for i := 0; i < m; i++ {
		found := false
		for j := 0; j < n+m; j++ {
			if right[i] == merged[j] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Number %d from right array not found in merged array", right[i])
		}
	}

	// test if merged is always increasing
	for i := 1; i < n+m; i++ {
		if merged[i-1] > merged[i] {
			t.Errorf("Merged array not sorted")
		}
	}

}

func TestParallelMergeSort(t *testing.T) {
	// get n random numbers
	n := 1000
	I := make([]int, n)
	for i := 0; i < n; i++ {
		I[i] = rand.Intn(n)
	}

	// create scratch space
	S := make([]int, n)

	var wg sync.WaitGroup
	wg.Add(1)
	// call the parallel merge sort
	Parallelmergesort(I, S, &wg)

	// test if I is sorted
	for i := 1; i < n; i++ {
		if I[i-1] >= I[i] {
			t.Errorf("Array not sorted")
		}
	}
}

func TestParallelMerge(t *testing.T) {
	A := []int{1, 3, 5, 7, 9, 11, 13, 15}
	B := []int{2, 4, 6, 8, 10, 12, 14, 16}

	// Number of processors (goroutines) to use
	p := 4

	// Perform the parallel merge
	result := Parallel_merge(A, B, p)

	// Print the result
	fmt.Println("Merged array:", result)
}

func TestBinarySearch(t *testing.T) {
	A := []int{1, 3, 5, 7, 9, 11, 13, 15}
	val := 4

	idx := binary_search_index(A, val)
	if idx != 2 {
		t.Errorf("Expected index 2, got %d", idx)
	}

	A = []int{1, 3, 5, 7, 9, 11, 13, 15}
	val = 10

	idx = binary_search_index(A, val)
	if idx != 5 {
		t.Errorf("Expected index 2, got %d", idx)
	}
}
