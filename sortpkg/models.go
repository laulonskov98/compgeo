package sortpkg

import (
	"fmt"
	"sync"
)

func Parallelmergesort(I []int, S []int, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(I) <= 1 {
		copy(S, I)
		return
	}
	mid := len(I) / 2
	I_left := I[:mid]
	I_right := I[mid:]
	S_left := S[:mid]
	S_right := S[mid:]

	var childWg sync.WaitGroup
	childWg.Add(2)
	go func() {
		Parallelmergesort(I_left, S_left, &childWg)
	}()
	go func() {
		Parallelmergesort(I_right, S_right, &childWg)
	}()

	childWg.Wait()
	merged := sequential_merge(S_left, S_right, I)
	copy(S, merged)
}

func sequential_merge(left []int, right []int, merged []int) []int {

	fmt.Println("left", left)
	fmt.Println("right", right)
	fmt.Println("merged", merged)
	fmt.Println("")
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			merged[k] = left[i]
			i++
		} else {
			merged[k] = right[j]
			j++
		}
		k++
	}
	for i < len(left) {
		merged[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		merged[k] = right[j]
		j++
		k++
	}
	return merged
}

func binary_search_index(B []int, end_val_of_chunk int) int {
	start := 0
	end := len(B) - 1
	for start <= end {
		mid := start + (end-start)/2
		if B[mid] == end_val_of_chunk {
			return mid
		} else if B[mid] < end_val_of_chunk {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return start
}

func AssymMerge(a_boundary_cutoffs []int, i, chunkSize int, A, B, result []int) {
	a_chunk_size := a_boundary_cutoffs[i] - a_boundary_cutoffs[i-1]
	k := int(a_chunk_size / chunkSize)

	b_boundary_cutoffs := make([]int, k+1)
	b_boundary_cutoffs[0] = 0
	for j := 1; j < k; j++ {
		// merge subarrays of A and B
		end_val_of_chunk_a := A[a_boundary_cutoffs[i-1]+j*chunkSize]
		idx := binary_search_index(B, end_val_of_chunk_a)
		b_boundary_cutoffs[j] = idx
	}

	total_used_lenght := 0
	for j := 0; j < k; j++ {
		a_start := a_boundary_cutoffs[i-1] + j*chunkSize
		a_end := a_start + chunkSize
		if j == k {
			a_end = a_boundary_cutoffs[i]
		}
		b_start := b_boundary_cutoffs[j-1]
		b_end := b_boundary_cutoffs[j]
		total_length := a_end - a_start + b_end - b_start
		sequential_merge(A[a_start:a_end], B[b_start:b_end], result[total_used_lenght:total_length])
	}
}

func Parallel_merge(A, B []int, p int) []int {
	n := len(B)
	if len(B) != n {
		panic("Arrays A and B must be of the same length")
	}

	// Result array to store the merged array
	result := make([]int, 2*n)

	// Create subarrays for parallel processing
	var wg sync.WaitGroup
	chunkSize := n / p

	//binary search for the end of each chunk in B in A
	a_boundary_cutoffs := make([]int, p)
	for i := 0; i < p; i++ {
		end_val_of_chunk := B[(i+1)*chunkSize-1]
		idx := binary_search_index(A, end_val_of_chunk)
		a_boundary_cutoffs[i] = idx
		result[(i+1)*chunkSize-1+idx] = end_val_of_chunk
	}

	fmt.Println("PRE RES", result)
	for i := 1; i < p; i++ {
		wg.Add(1)

		start := (i-1)*chunkSize + a_boundary_cutoffs[i-1]
		end := start + chunkSize + a_boundary_cutoffs[i]

		if i == p-1 {
			end = len(result) - 1
		}
		fmt.Println(start, end, a_boundary_cutoffs, n)

		b_start := (i - 1) * chunkSize
		b_end := i * chunkSize
		go func(i, start, end int) {
			defer wg.Done()
			result_cutout := result[start:end]
			fmt.Println(start, end, len(result_cutout))

			a_chunk_size := a_boundary_cutoffs[i] - a_boundary_cutoffs[i-1]
			fmt.Println("a chunk size", a_chunk_size, chunkSize)
			if a_chunk_size > chunkSize {
				AssymMerge(a_boundary_cutoffs, i, chunkSize, A, B, result_cutout)
			} else {
				sequential_merge(A[a_boundary_cutoffs[i-1]:a_boundary_cutoffs[i]], B[b_start:b_end], result_cutout)
			}
		}(i, start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return result
}
