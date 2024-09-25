package main

import (
	"compgeo/sortpkg"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Hello, World!")
	arr := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Println("Unsorted:", arr)

	// Scratch space
	scratch := make([]int, len(arr))

	// WaitGroup to wait for the top-level merge sort to complete
	var wg sync.WaitGroup
	wg.Add(1)

	// Call the parallel merge sort
	sortpkg.Basic_Parallel_mergesort(arr, scratch, &wg, 0, 2)

	// Wait for the sorting to finish
	wg.Wait()

	fmt.Println("Sorted:", arr)
}
