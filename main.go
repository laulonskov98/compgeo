package main

import (
	"compgeo/sortpkg"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func take_time_basic_sort_threads() {
	rand.Seed(42)
	recorded_times := make([][]int64, 10)

	for max_depth := 0; max_depth < 10; max_depth++ {
		sub_recorded_times := make([]int64, 10)
		for i := 0; i < 10; i++ {
			n := 500000
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rand.Intn(n)
			}

			// Scratch space
			scratch := make([]int, len(arr))

			// Call the parallel merge sort
			now := time.Now()
			sortpkg.Basic_Parallel_Mergesort(arr, scratch, 0, max_depth)
			sub_recorded_times[i] = time.Since(now).Nanoseconds()
		}
		recorded_times[max_depth] = sub_recorded_times
	}

	// convert recorded times to csv file
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println("(", i, ",", recorded_times[i][j], ")")
		}
	}

	fmt.Println("Recorded times:", recorded_times)
}

func take_time_basicsort_increasing_n() {
	rand.Seed(42)
	max_itr := 70
	multiple := 20000

	recorded_times := make([][]float64, max_itr)
	for itr := 1; itr < max_itr+1; itr++ {
		sub_recorded_times := make([]float64, 10)

		for i := 0; i < 10; i++ {
			n := itr * multiple
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rand.Intn(n)
			}

			fmt.Println(n)
			// Scratch space
			scratch := make([]int, len(arr))

			// Call the parallel merge sort
			now := time.Now()
			sortpkg.Basic_Parallel_Mergesort(arr, scratch, 0, 3)
			sub_recorded_times[i] = float64(time.Since(now).Nanoseconds()) / (float64(n) * float64(math.Log2(float64(n))))
		}
		recorded_times[itr-1] = sub_recorded_times
	}

	// convert recorded times to csv file
	for i := 1; i < max_itr+1; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println("(", i*multiple, ",", recorded_times[i-1][j], ")")
		}
	}

}

func take_time_sequantial_merge(array1, array2 []int) {

}
func take_time_assym_merge(array1, array2 []int) {}

func generate_variable_overlap_arrays(size int, overlap float64) ([]int, []int) {
	half_size := size / 2
	array1 := make([]int, half_size)
	array2 := make([]int, half_size)

	overlap_size := int(float64(half_size)*(1+overlap)) - half_size
	rand_max := half_size + overlap_size
	for i := 0; i < half_size; i++ {
		array1[i] = rand.Intn(rand_max)
		array2[i] = rand.Intn(rand_max) + half_size - overlap_size
	}

	//sort both arrays
	sort.Ints(array1)
	sort.Ints(array2)

	return array1, array2
}

func generate_variable_overlap_array_sorting(size int, overlap float64) []int {
	half_size := size / 2
	array := make([]int, size)

	overlap_size := int(float64(half_size)*(1+overlap)) - half_size
	rand_max := half_size + overlap_size
	for i := 0; i < half_size; i++ {
		array[i] = rand.Intn(rand_max)
		array[i+half_size] = rand.Intn(rand_max) + half_size - overlap_size
	}

	return array
}

func generate_slight_overlap_arrays(size int) ([]int, []int) {
	half_size := size / 2
	array1 := make([]int, half_size)
	array2 := make([]int, half_size)

	overlap_size := int(float64(half_size)*1.1) - half_size
	rand_max := half_size + overlap_size
	for i := 0; i < half_size; i++ {
		array1[i] = rand.Intn(rand_max)
		array2[i] = rand.Intn(rand_max) + half_size - overlap_size
	}

	//sort both arrays
	sort.Ints(array1)
	sort.Ints(array2)

	return array1, array2
}

func generate_fully_overlap_arrays(size int) ([]int, []int) {
	half_size := size / 2
	array1 := make([]int, half_size)
	array2 := make([]int, half_size)

	rand_max := half_size
	for i := 0; i < half_size; i++ {
		array1[i] = rand.Intn(rand_max)
		array2[i] = rand.Intn(rand_max)
	}

	//sort both arrays
	sort.Ints(array1)
	sort.Ints(array2)

	return array1, array2
}

func generate_non_overlapping_arrays(size int) ([]int, []int) {
	half_size := size / 2
	array1 := make([]int, half_size)
	array2 := make([]int, half_size)

	rand_max := half_size
	for i := 0; i < half_size; i++ {
		array1[i] = rand.Intn(rand_max)
		array2[i] = rand.Intn(rand_max) + half_size + 1
	}

	//sort both arrays
	sort.Ints(array1)
	sort.Ints(array2)

	return array1, array2
}

func take_time_merge_increasing_overlap() {
	rand.Seed(42)
	overlapper_steps := 20
	overlap_pct_step := 1 / float64(overlapper_steps)
	recorded_times := make([][]int64, overlapper_steps+1)
	for itr := 0; itr < overlapper_steps+1; itr++ {
		sub_recorded_times := make([]int64, 10)

		for i := 0; i < 10; i++ {
			n := 500000
			overlap_pct := float64(itr) * overlap_pct_step
			res := make([]int, 2*n)
			array1, array2 := generate_variable_overlap_arrays(n, overlap_pct)
			// Call the parallel merge sort
			now := time.Now()
			sortpkg.Parallel_merge(array1, array2, res, 8)
			sub_recorded_times[i] = time.Since(now).Nanoseconds()
		}
		recorded_times[itr] = sub_recorded_times
	}

	// convert recorded times to csv file
	for i := 0; i < overlapper_steps+1; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println("(", float64(i)*1/float64(overlapper_steps), ",", recorded_times[i][j], ")")
		}
	}
}

func take_time_sort_increasing_overlap() {
	rand.Seed(42)
	overlapper_steps := 20
	overlap_pct_step := 1 / float64(overlapper_steps)
	recorded_times := make([][]int64, overlapper_steps)
	for itr := 1; itr < overlapper_steps+1; itr++ {
		sub_recorded_times := make([]int64, 10)

		for i := 0; i < 10; i++ {
			n := 80000
			overlap_pct := float64(itr) * overlap_pct_step
			array := generate_variable_overlap_array_sorting(n, overlap_pct)
			// Call the parallel merge sort
			scratch := make([]int, len(array))
			resultCh := make(chan int, 1)
			sortpkg.Parallel_Mergesort_comparisons(array, scratch, 0, 1000, resultCh, 10000)
			res := <-resultCh
			sub_recorded_times[i] = int64(res)
		}
		recorded_times[itr-1] = sub_recorded_times
	}

	// convert recorded times to csv file
	for i := 1; i < overlapper_steps+1; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println("(", float64(i-1)*1/float64(overlapper_steps), ",", recorded_times[i-1][j], ")")
		}
	}
}

func take_time_merge_functions() {
	max_itr := 50
	recorded_times := make([][]int64, max_itr)
	for itr := 1; itr < max_itr+1; itr++ {
		sub_recorded_times := make([]int64, 10)

		for i := 0; i < 10; i++ {
			n := itr * 10000

			res := make([]int, 2*n)
			array1, array2 := generate_fully_overlap_arrays(n)
			// Call the parallel merge sort
			now := time.Now()
			sortpkg.Parallel_merge(array1, array2, res, 8)
			sub_recorded_times[i] = time.Since(now).Nanoseconds()
		}
		recorded_times[itr-1] = sub_recorded_times
	}

	// convert recorded times to csv file
	for i := 1; i < max_itr+1; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println("(", i*10000, ",", recorded_times[i-1][j], ")")
		}
	}

}
func main() {
	take_time_sort_increasing_overlap()

}
