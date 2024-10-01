package sortpkg

import (
	"math"
	"sync"
)

/*
start
                   m
 		    m             1
		m      3        1    2
	  m  4    3  5     1 6  2 7

*/

var ComparisonCount int

func IncrementComparisonCount() {
	ComparisonCount++
}

func get_max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getMax(arr []int) int {
	if len(arr) == 0 {
		panic("Cannot find max of an empty slice")
	}

	max := arr[0]
	for _, value := range arr {
		if value > max {
			max = value
		}
	}
	return max
}

func Basic_Parallel_Mergesort(I []int, S []int, depth, max_depth int) {
	if len(I) <= 1 {
		copy(S, I)
		return
	}

	mid := len(I) / 2
	left, right := I[:mid], I[mid:]
	S_left, S_right := S[:mid], S[mid:]

	if depth < max_depth {
		var wg sync.WaitGroup
		wg.Add(1)

		// Goroutine for the left half
		go func() {
			defer wg.Done()
			Basic_Parallel_Mergesort(left, S_left, depth+1, max_depth)
		}()
		Basic_Parallel_Mergesort(right, S_right, depth+1, max_depth)
		wg.Wait()
	} else {
		Basic_Parallel_Mergesort(left, S_left, depth+1, max_depth)
		Basic_Parallel_Mergesort(right, S_right, depth+1, max_depth)
	}

	if depth%2 == 0 {
		// On even depths, merge frSm I to S
		Sequential_merge(left, right, S)
	} else {
		// On odd depths, mergS from S to I
		Sequential_merge(S_left, S_right, I)
	}
}

func Basic_Parallel_Mergesort_comparisons(I []int, S []int, depth, max_depth int, ch chan int) {
	if len(I) <= 1 {
		copy(S, I)
		ch <- 1
		return
	}

	mid := len(I) / 2
	left, right := I[:mid], I[mid:]
	S_left, S_right := S[:mid], S[mid:]

	resultCh1 := make(chan int, 1)
	resultCh2 := make(chan int, 1)
	operations := 0
	if depth < max_depth {
		var wg sync.WaitGroup
		wg.Add(1)

		// Goroutine for the left half
		go func() {
			defer wg.Done()
			Basic_Parallel_Mergesort_comparisons(left, S_left, depth+1, max_depth, resultCh1)
		}()
		Basic_Parallel_Mergesort_comparisons(right, S_right, depth+1, max_depth, resultCh2)
		wg.Wait()
		leftmax := <-resultCh1
		rightmax := <-resultCh2
		operations = get_max(leftmax, rightmax)
	} else {
		Basic_Parallel_Mergesort_comparisons(left, S_left, depth+1, max_depth, resultCh1)
		Basic_Parallel_Mergesort_comparisons(right, S_right, depth+1, max_depth, resultCh2)
		leftmax := <-resultCh1
		rightmax := <-resultCh2
		operations = leftmax + rightmax
	}

	close(resultCh1)
	close(resultCh2)

	merge_compares := 0
	if depth%2 == 0 {
		// On even depths, merge frSm I to S
		merge_compares = Sequential_merge_count_comparisons(left, right, S)
	} else {
		// On odd depths, mergS from S to I
		merge_compares = Sequential_merge_count_comparisons(S_left, S_right, I)
	}
	ch <- operations + merge_compares
}

func Parallel_Mergesort(I []int, S []int, depth, max_depth int) {
	if len(I) <= 1 {
		copy(S, I)
		return
	}
	mid := len(I) / 2
	I_left, I_right := I[:mid], I[mid:]
	S_left, S_right := S[:mid], S[mid:]

	if depth < max_depth {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			Parallel_Mergesort(I_left, S_left, depth+1, max_depth)
		}()
		Parallel_Mergesort(I_right, S_right, depth+1, max_depth)
		wg.Wait()
	} else {
		Parallel_Mergesort(I_left, S_left, depth+1, max_depth)
		Parallel_Mergesort(I_right, S_right, depth+1, max_depth)
	}

	if depth%2 == 0 {
		// On even depths, merge from I to S
		Parallel_merge(I_left, I_right, S, 20)
	} else {
		// On odd depths, merge from S to I
		Parallel_merge(S_left, S_right, I, 20)
	}
}

func Parallel_Mergesort_comparisons(I []int, S []int, depth, max_depth int, ch chan int, merge_processors int) {
	if len(I) <= 1 {
		copy(S, I)
		ch <- 1
		return
	}
	mid := len(I) / 2
	I_left, I_right := I[:mid], I[mid:]
	S_left, S_right := S[:mid], S[mid:]

	resultCh1 := make(chan int, 1)
	resultCh2 := make(chan int, 1)
	operations := 0
	if depth < max_depth {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			Parallel_Mergesort_comparisons(I_left, S_left, depth+1, max_depth, resultCh1, merge_processors)
		}()
		Parallel_Mergesort_comparisons(I_right, S_right, depth+1, max_depth, resultCh2, merge_processors)

		wg.Wait()
		max_1 := <-resultCh1
		max_2 := <-resultCh2
		operations = get_max(max_1, max_2)
	} else {
		Parallel_Mergesort_comparisons(I_left, S_left, depth+1, max_depth, resultCh1, merge_processors)
		Parallel_Mergesort_comparisons(I_right, S_right, depth+1, max_depth, resultCh2, merge_processors)
		max_1 := <-resultCh1
		max_2 := <-resultCh2
		operations = max_1 + max_2
	}

	close(resultCh1)
	close(resultCh2)

	comparisons := 0
	if depth%2 == 0 {
		// On even depths, merge from I to S
		comparisons = Parallel_merge_comparison_counter(I_left, I_right, S, int(merge_processors))
	} else {
		// On odd depths, merge from S to I
		comparisons = Parallel_merge_comparison_counter(S_left, S_right, I, int(merge_processors))
	}
	ch <- operations + comparisons
}

func Sequential_merge(left []int, right []int, merged []int) []int {
	compares := 0
	if len(left) == 1 && len(right) == 1 {
		compares++
		if left[0] > right[0] {
			merged[1] = left[0]
			merged[0] = right[0]
		} else {
			merged[1] = right[0]
			merged[0] = left[0]
		}
		return merged
	}

	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		compares++
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
		compares++
		merged[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		compares++
		merged[k] = right[j]
		j++
		k++
	}
	return merged
}

func Sequential_merge_count_comparisons(left []int, right []int, merged []int) int {
	compares := 0
	if len(left) == 1 && len(right) == 1 {
		compares++
		if left[0] > right[0] {
			merged[1] = left[0]
			merged[0] = right[0]
		} else {
			merged[1] = right[0]
			merged[0] = left[0]
		}
		return compares
	}

	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		compares++
		if left[i] <= right[j] {
			merged[k] = left[i]
			i++
		} else {
			merged[k] = right[j]
			j++
		}
		k++
	}
	copy(merged[k:], left[i:])
	copy(merged[k+len(left)-i:], right[j:])
	return compares
}

func Sequential_merge_2(left, right, result []int) []int {
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result[k] = left[i]
			i++
		} else {
			result[k] = right[j]
			j++
		}
		k++
	}
	copy(result[k:], left[i:])
	copy(result[k+len(left)-i:], right[j:])
	return result
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

func binary_search_index_comparisons(B []int, end_val_of_chunk int) (int, int) {
	start := 0
	end := len(B) - 1
	comparisons := 0
	for start <= end {
		comparisons++
		mid := start + (end-start)/2
		if B[mid] == end_val_of_chunk {
			return mid, comparisons
		} else if B[mid] < end_val_of_chunk {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return start, comparisons
}

func AssymMerge(chunk_size int, A, B, result []int) {
	a_chunk_size := len(A)
	k := int(a_chunk_size / chunk_size)
	if a_chunk_size%chunk_size != 0 {
		k++
	}

	b_boundary_cutoffs := make([]int, k+1)
	b_boundary_cutoffs[0] = 0
	for j := 1; j < k+1; j++ {
		// merge subarrays of A and B
		end_val_of_chunk_a := A[a_chunk_size-1]
		if j != k {
			end_val_of_chunk_a = A[j*chunk_size-1]
		}
		idx := binary_search_index(B, end_val_of_chunk_a)
		b_boundary_cutoffs[j] = idx
	}

	var wg sync.WaitGroup
	total_used_length := 0
	for j := 1; j < k+1; j++ {
		a_start_idx := (j - 1) * chunk_size
		a_end_idx := (j) * chunk_size

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
		wg.Add(1)

		go func(A, B, result []int) {
			defer wg.Done()
			Sequential_merge(A, B, result)
		}(A[a_start_idx:a_end_idx], B[b_start_idx:b_end_idx], result[total_used_length:total_used_length+total_length])
		total_used_length += total_length
	}
	wg.Wait()

}

func AssymMerge_comparison_counter(chunk_size int, A, B, result []int) int {
	a_chunk_size := len(A)
	k := int(a_chunk_size / chunk_size)
	if a_chunk_size%chunk_size != 0 {
		k++
	}

	b_boundary_cutoffs := make([]int, k+1)
	b_boundary_cutoffs[0] = 0
	max_bin_comp := 0
	for j := 1; j < k+1; j++ {
		// merge subarrays of A and B
		end_val_of_chunk_a := A[a_chunk_size-1]
		if j != k {
			end_val_of_chunk_a = A[j*chunk_size-1]
		}
		idx, bin_comps := binary_search_index_comparisons(B, end_val_of_chunk_a)
		if bin_comps > max_bin_comp {
			max_bin_comp = bin_comps
		}

		b_boundary_cutoffs[j] = idx
	}

	var wg sync.WaitGroup
	compare_array := make([]int, k)
	total_used_length := 0
	for j := 1; j < k+1; j++ {
		a_start_idx := (j - 1) * chunk_size
		a_end_idx := (j) * chunk_size

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
		wg.Add(1)

		go func(A, B, result []int, idx int) {
			defer wg.Done()
			compares := Sequential_merge_count_comparisons(A, B, result)
			compare_array[idx] = compares
		}(A[a_start_idx:a_end_idx], B[b_start_idx:b_end_idx], result[total_used_length:total_used_length+total_length], j-1)
		total_used_length += total_length
	}
	wg.Wait()

	return getMax(compare_array) + max_bin_comp

}

func get_chunk_variables(n, p int) (chunk_size int, no_of_chunks int, processors int) {
	if p > 100 {
		chunk_size = int(math.Log(float64(n)))
		if chunk_size == 0 {
			chunk_size = 1
		}
		no_of_chunks = n/chunk_size + 1

		processors = no_of_chunks
		return chunk_size, no_of_chunks, processors
	}
	if n <= p {
		p = 1
	}
	chunk_size = n / p

	if n%p != 0 {
		chunk_size++
	}

	// divide and round up
	no_of_chunks = n / chunk_size

	if n%no_of_chunks != 0 {
		no_of_chunks++
	}

	return chunk_size, no_of_chunks, p
}

// todo: tilføj et array som der merges ind i som parameter dertil.
func Parallel_merge(A, B []int, result []int, p int) []int {
	if len(A) == 1 && len(B) == 1 {
		if A[0] > B[0] {
			result[0] = B[0]
			result[1] = A[0]
		} else {
			result[0] = A[0]
			result[1] = B[0]
		}
		return result
	}

	n := len(B)
	if len(B) != n {
		panic("Arrays A and B must be of the same length")
	}

	chunk_size, no_of_chunks, p_val := get_chunk_variables(n, p)

	// Create subarrays for parallel processing
	var wg sync.WaitGroup

	a_boundary_cutoffs := make([]int, no_of_chunks)

	//binary search for the end of each chunk in B in A
	for i := 0; i < no_of_chunks; i++ {
		end_val_of_chunk := B[n-1]
		if i < no_of_chunks-1 {
			end_val_of_chunk = B[(i+1)*chunk_size-1]
		}

		idx := binary_search_index(A, end_val_of_chunk)
		a_boundary_cutoffs[i] = idx

		if (i+1)*chunk_size-1+idx > len(result)-1 {
			result[len(result)-1] = end_val_of_chunk
		} else {
			result[(i+1)*chunk_size-1+idx] = end_val_of_chunk
		}
	}

	for i := 0; i < no_of_chunks; i++ {

		b_start := (i) * chunk_size
		b_end := (i + 1) * chunk_size

		a_start := 0
		if i > 0 {
			a_start = a_boundary_cutoffs[i-1]
		}
		a_end := a_boundary_cutoffs[i]
		if a_end == 0 && a_start == 0 && i == p_val-1 {
			a_end = len(A) - 1
		}

		start := 0
		end := start + chunk_size + a_end
		if i > 0 {
			start = (i)*chunk_size + a_boundary_cutoffs[i-1]
			end = start + chunk_size + (a_end - a_start)
		}

		if i == no_of_chunks-1 {
			a_end = len(A)
			b_end = len(B)
			end = len(result)
		}
		wg.Add(1)

		go func(i, start, end int) {
			defer wg.Done()
			result_cutout := result[start:end]

			a_chunk_size := a_boundary_cutoffs[i]
			if i > 0 {
				a_chunk_size -= a_boundary_cutoffs[i-1]
			}
			if a_chunk_size > chunk_size {
				AssymMerge(chunk_size, A[a_start:a_end], B[b_start:b_end], result_cutout)
			} else {
				Sequential_merge(A[a_start:a_end], B[b_start:b_end], result_cutout)
			}
		}(i, start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return result
}

func Parallel_merge_comparison_counter(A, B []int, result []int, p int) int {

	local_compares := 0
	if len(A) == 1 && len(B) == 1 {
		if A[0] > B[0] {
			result[0] = B[0]
			result[1] = A[0]
		} else {
			result[0] = A[0]
			result[1] = B[0]
		}
		return 1
	}

	n := len(B)
	if len(B) != n {
		panic("Arrays A and B must be of the same length")
	}

	chunk_size, no_of_chunks, p_val := get_chunk_variables(n, p)

	// Create subarrays for parallel processing
	var wg sync.WaitGroup

	a_boundary_cutoffs := make([]int, no_of_chunks)

	//binary search for the end of each chunk in B in A
	for i := 0; i < no_of_chunks; i++ {
		local_compares++
		end_val_of_chunk := B[n-1]
		if i < no_of_chunks-1 {
			end_val_of_chunk = B[(i+1)*chunk_size-1]
		}

		idx := binary_search_index(A, end_val_of_chunk)
		a_boundary_cutoffs[i] = idx

		if (i+1)*chunk_size-1+idx > len(result)-1 {
			result[len(result)-1] = end_val_of_chunk
		} else {
			result[(i+1)*chunk_size-1+idx] = end_val_of_chunk
		}
	}

	compare_array := make([]int, no_of_chunks)
	for i := 0; i < no_of_chunks; i++ {
		local_compares++

		b_start := (i) * chunk_size
		b_end := (i + 1) * chunk_size

		a_start := 0
		if i > 0 {
			a_start = a_boundary_cutoffs[i-1]
		}
		a_end := a_boundary_cutoffs[i]
		if a_end == 0 && a_start == 0 && i == p_val-1 {
			a_end = len(A) - 1
		}

		start := 0
		end := start + chunk_size + a_end
		if i > 0 {
			start = (i)*chunk_size + a_boundary_cutoffs[i-1]
			end = start + chunk_size + (a_end - a_start)
		}

		if i == no_of_chunks-1 {
			a_end = len(A)
			b_end = len(B)
			end = len(result)
		}
		wg.Add(1)

		go func(i, start, end int) {
			defer wg.Done()
			result_cutout := result[start:end]

			a_chunk_size := a_boundary_cutoffs[i]
			if i > 0 {
				a_chunk_size -= a_boundary_cutoffs[i-1]
			}
			compares := 0
			if a_chunk_size > chunk_size {
				compares = AssymMerge_comparison_counter(chunk_size, A[a_start:a_end], B[b_start:b_end], result_cutout)
			} else {
				compares = Sequential_merge_count_comparisons(A[a_start:a_end], B[b_start:b_end], result_cutout)
			}
			compare_array[i] = compares
		}(i, start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return getMax(compare_array) + local_compares
}

// todo: tilføj et array som der merges ind i som parameter dertil.
func Parallel_merge_only_seq(A, B []int, result []int, p int) []int {
	if len(A) == 1 && len(B) == 1 {
		if A[0] > B[0] {
			result[0] = B[0]
			result[1] = A[0]
		} else {
			result[0] = A[0]
			result[1] = B[0]
		}
		return result
	}

	n := len(B)
	if len(B) != n {
		panic("Arrays A and B must be of the same length")
	}

	chunk_size, no_of_chunks, p_val := get_chunk_variables(n, p)

	// Create subarrays for parallel processing
	var wg sync.WaitGroup

	a_boundary_cutoffs := make([]int, no_of_chunks)

	//binary search for the end of each chunk in B in A
	for i := 0; i < no_of_chunks; i++ {
		end_val_of_chunk := B[n-1]
		if i < no_of_chunks-1 {
			end_val_of_chunk = B[(i+1)*chunk_size-1]
		}

		idx := binary_search_index(A, end_val_of_chunk)
		a_boundary_cutoffs[i] = idx

		if (i+1)*chunk_size-1+idx > len(result)-1 {
			result[len(result)-1] = end_val_of_chunk
		} else {
			result[(i+1)*chunk_size-1+idx] = end_val_of_chunk
		}
	}

	for i := 0; i < no_of_chunks; i++ {

		b_start := (i) * chunk_size
		b_end := (i + 1) * chunk_size

		a_start := 0
		if i > 0 {
			a_start = a_boundary_cutoffs[i-1]
		}
		a_end := a_boundary_cutoffs[i]
		if a_end == 0 && a_start == 0 && i == p_val-1 {
			a_end = len(A) - 1
		}

		start := 0
		end := start + chunk_size + a_end
		if i > 0 {
			start = (i)*chunk_size + a_boundary_cutoffs[i-1]
			end = start + chunk_size + (a_end - a_start)
		}

		if i == no_of_chunks-1 {
			a_end = len(A)
			b_end = len(B)
			end = len(result)
		}
		wg.Add(1)

		go func(i, start, end int) {
			defer wg.Done()
			result_cutout := result[start:end]

			a_chunk_size := a_boundary_cutoffs[i]
			if i > 0 {
				a_chunk_size -= a_boundary_cutoffs[i-1]
			}
			Sequential_merge(A[a_start:a_end], B[b_start:b_end], result_cutout)
		}(i, start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return result
}
