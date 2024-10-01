package sortpkg

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
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
	merged = Sequential_merge(left, right, merged)

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
	n := 800000
	I := make([]int, n)
	for i := 0; i < n; i++ {
		I[i] = rand.Intn(n)
	}
	I_copy := make([]int, n)
	copy(I_copy, I)

	counts_of_expected_elements := make(map[int]int)
	for i := 0; i < n; i++ {
		counts_of_expected_elements[I[i]]++
	}

	// create scratch space
	S := make([]int, n)
	S_Copy := make([]int, n)

	start := time.Now()
	resultCh := make(chan int, 1)
	Basic_Parallel_Mergesort_comparisons(I_copy, S_Copy, 0, 4, resultCh)
	res := <-resultCh
	close(resultCh)

	fmt.Println("max compares basic", res)

	fmt.Println("Basic parallel sort took:", time.Since(start))

	// call the parallel merge sort
	para_start := time.Now()
	resultCh = make(chan int, 1)
	Parallel_Mergesort_comparisons(I, S, 0, 3, resultCh, 2)
	res = <-resultCh
	fmt.Println("max compares parallel", res)
	fmt.Println("Parallel sort took:", time.Since(para_start))
	// test if I is sorted
	for i := 1; i < n; i++ {
		if S[i-1] > S[i] {
			t.Errorf("Array not sorted {i-1: %d, i: %d}", I[i-1], I[i])
		}
	}

	counts_of_result_elements := make(map[int]int)
	for i := 0; i < n; i++ {
		counts_of_result_elements[S[i]]++
	}

	// iterate k,v pairs of the expected elements
	for k, v := range counts_of_expected_elements {
		// check if the element is in the result
		if counts_of_result_elements[k] != v {
			t.Errorf("Element %d not found in result", k)
		}
	}

}

func TestAssymMerge(t *testing.T) {
	A := []int{34, 39, 53}
	B := []int{37, 55}
	res_array := []int{34, 37, 39, 53, 55}

	// Number of processors (goroutines) to use

	result_array := make([]int, len(A)+len(B))

	// Perform the parallel merge
	AssymMerge(2, A, B, result_array)
	// Print the result
	fmt.Println("Merged array:", result_array)
	for i := 0; i < len(result_array); i++ {
		if res_array[i] != result_array[i] {
			t.Errorf("Array not sorted")
			fmt.Println(result_array[i-1], result_array[i])
		}
	}
}

func TestTimeTakingBasicParallelMergeSort(t *testing.T) {
	recorded_times := make([][]int64, 10)
	for max_depth := 1; max_depth < 10; max_depth++ {
		sub_recorded_times := make([]int64, 10)
		for i := 0; i < 10; i++ {
			// get n random numbers
			n := 10000
			I := make([]int, n)
			for i := 0; i < n; i++ {
				I[i] = rand.Intn(n)
			}

			// create scratch space
			S := make([]int, n)

			var wg sync.WaitGroup
			wg.Add(1)
			// call the parallel merge sort
			now := time.Now()
			resultCh := make(chan int)
			Basic_Parallel_Mergesort_comparisons(I, S, 0, max_depth, resultCh)
			sub_recorded_times[i-1] = time.Since(now).Nanoseconds()
			for i := 1; i < n; i++ {
				if I[i-1] > I[i] {
					t.Errorf("Array not sorted")
					fmt.Println(I[i-1], I[i])
				}
			}
		}
		recorded_times[max_depth-1] = sub_recorded_times

	}
	fmt.Println("Recorded times:", recorded_times)
}

func TestParallelMerge(t *testing.T) {
	A := []int{1, 3, 5, 7, 9, 11, 13, 15}
	B := []int{2, 4, 6, 8, 10, 12, 14, 16}

	// Number of processors (goroutines) to use
	p := 4

	result_array := make([]int, len(A)+len(B))

	// Perform the parallel merge
	result := Parallel_merge(A, B, result_array, p)
	// Print the result
	fmt.Println("Merged array:", result)
	for i := 0; i < len(result); i++ {
		if i+1 != result[i] {
			t.Errorf("Array not sorted")
		}
	}
}

func TestParallelMergeUneven(t *testing.T) {
	A := []int{1, 2, 3, 4, 5, 6, 7, 9, 11, 12, 13, 14, 15, 17}
	B := []int{0, 8, 10, 16, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

	// Number of processors (goroutines) to use
	p := 4

	result_array := make([]int, len(A)+len(B))

	// Perform the parallel merge
	result := Parallel_merge(A, B, result_array, p)

	// Print the result
	fmt.Println("Merged array:", result)
	for i := 0; i < len(result); i++ {
		if i != result[i] {
			t.Errorf("Array not sorted {i: %d, result[i]: %d}", i, result[i])

		}
	}
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
