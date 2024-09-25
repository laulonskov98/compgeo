package sortpkg

import (
	"fmt"
	"sync"
)

func Basic_Parallel_mergesort(I []int, S []int, wg *sync.WaitGroup, depth, max_depth int) {
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
	if depth < max_depth {
		childWg.Add(1)

		go func() {
			Basic_Parallel_mergesort(I_left, S_left, &childWg, depth+1, max_depth)
		}()

		childWg.Wait()
	} else {
		Basic_Parallel_mergesort(I_left, S_left, &childWg, depth+1, max_depth)
		Basic_Parallel_mergesort(I_right, S_right, &childWg, depth+1, max_depth)
	}

	merged := sequential_merge(S_left, S_right, I)
	copy(S, merged)
}

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
	merged := Parallel_merge(S_left, S_right, I, 4)
	copy(S, merged)
}

func sequential_merge(left []int, right []int, merged []int) []int {
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

func AssymMerge(a_boundary_cutoffs []int, a_start, a_end, chunkSize int, A, B, result []int) {
	fmt.Println("")
	fmt.Println("AssymMerge", A, B, result)
	fmt.Println("chunk", chunkSize)
	a_chunk_size := a_end - a_start - 1
	k := int(a_chunk_size / chunkSize)
	if a_chunk_size%chunkSize != 0 {
		k++
	}

	b_boundary_cutoffs := make([]int, k+1)
	b_boundary_cutoffs[0] = 0
	for j := 1; j < k+1; j++ {
		// merge subarrays of A and B
		end_val_of_chunk_a := A[a_end-1]
		if j != k {
			end_val_of_chunk_a = A[j*chunkSize-1]
		}
		fmt.Println("end_val_of_chunk_a", end_val_of_chunk_a)
		idx := binary_search_index(B, end_val_of_chunk_a)
		b_boundary_cutoffs[j] = idx
	}
	fmt.Println(b_boundary_cutoffs)

	total_used_length := 0
	for j := 1; j < k+1; j++ {
		a_start_idx := (j - 1) * chunkSize
		a_end_idx := (j) * chunkSize

		b_start_idx := 0
		b_end_idx := b_boundary_cutoffs[j]
		if j > 0 {
			b_start_idx = b_boundary_cutoffs[j-1]
		}
		if j == k {
			b_end_idx = len(B)
			a_end_idx = len(A)
		}

		total_length := a_end_idx - a_start_idx + b_end_idx - b_start_idx

		fmt.Println("RES BEFORE", result)
		fmt.Println("a_end_idx", a_end_idx, "b_end_idx", b_end_idx, "total_length", total_length, "total_used_length", total_used_length)
		sequential_merge(A[a_start_idx:a_end_idx], B[b_start_idx:b_end_idx], result[total_used_length:total_used_length+total_length])
		fmt.Println("RES AFTER", result)

		total_used_length += total_length
	}
	fmt.Println("AssymMerge result", result)
}

// todo: tilføj et array som der merges ind i som parameter dertil.
func Parallel_merge(A, B []int, result []int, p int) []int {
	n := len(B)
	if len(B) != n {
		panic("Arrays A and B must be of the same length")
	}

	// Create subarrays for parallel processing
	var wg sync.WaitGroup
	chunkSize := n / p

	//binary search for the end of each chunk in B in A
	if n%p != 0 {
		chunkSize++
	}
	a_boundary_cutoffs := make([]int, p)
	for i := 0; i < p; i++ {
		end_val_of_chunk := B[n-1]
		if i < p-1 {
			end_val_of_chunk = B[(i+1)*chunkSize-1]
		}
		idx := binary_search_index(A, end_val_of_chunk)
		a_boundary_cutoffs[i] = idx

		if (i+1)*chunkSize-1+idx > len(result)-1 {
			result[len(result)-1] = end_val_of_chunk
		} else {
			result[(i+1)*chunkSize-1+idx] = end_val_of_chunk
		}
	}

	fmt.Println("PRE RES", result)
	for i := 0; i < p; i++ {
		wg.Add(1)

		b_start := (i) * chunkSize
		b_end := (i + 1) * chunkSize

		a_start := 0
		if i > 0 {
			a_start = a_boundary_cutoffs[i-1]
		}
		a_end := a_boundary_cutoffs[i]

		start := 0
		end := start + chunkSize + a_end
		if i > 0 {
			start = (i)*chunkSize + a_boundary_cutoffs[i-1]
			end = start + chunkSize + (a_end - a_start)
		}

		fmt.Println("start", start, "end", end, "chunk size", chunkSize, "a_boundary_cutoffs", a_boundary_cutoffs[i])
		if i == p-1 {
			end = len(result)
			b_end = len(B) - 1
		}
		fmt.Println(start, end, a_boundary_cutoffs, n)

		go func(i, start, end int) {
			defer wg.Done()
			result_cutout := result[start:end]
			fmt.Println(start, end, len(result_cutout))

			a_chunk_size := a_boundary_cutoffs[i]
			if i > 0 {
				a_chunk_size -= a_boundary_cutoffs[i-1]
			}
			if a_chunk_size > chunkSize {
				fmt.Println("res", result_cutout, result)
				fmt.Println("a chunk size", a_chunk_size, chunkSize)
				AssymMerge(a_boundary_cutoffs, a_start, a_end, chunkSize, A[a_start:a_end], B[b_start:b_end], result_cutout)
			} else {
				sequential_merge(A[a_start:a_end], B[b_start:b_end], result_cutout)
			}
		}(i, start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return result
}
