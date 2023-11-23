package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
    "runtime"
)

// Function to generate a slice of size n with random numbers
func generateRandomSlice(size int) []int {
    slice := make([]int, size)
    for i := range slice {
        slice[i] = rand.Intn(100) // Random integers between 0 and 99
    }
    return slice
}

// Concurrent Quick Sort function
func concurrentQuickSort(slice []int, wg *sync.WaitGroup) {
    defer wg.Done()

    if len(slice) < 2 {
        return
    }

    // Partitioning the slice
    pivot := slice[len(slice)/2]
    left := make([]int, 0)
    right := make([]int, 0)

    for _, value := range slice {
        if value < pivot {
            left = append(left, value)
        } else if value > pivot {
            right = append(right, value)
        }
    }

    // Recursively sorting the partitions in parallel
    var leftWg, rightWg sync.WaitGroup
    leftWg.Add(1)
    go func() {
        concurrentQuickSort(left, &leftWg)
    }()
    rightWg.Add(1)
    go func() {
        concurrentQuickSort(right, &rightWg)
    }()
    leftWg.Wait()
    rightWg.Wait()

    // Merging results
    copy(slice[:len(left)], left)
    slice[len(left)] = pivot
    copy(slice[len(left)+1:], right)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Generate a slice of 1000 random numbers
    slice := generateRandomSlice(1000)

    fmt.Println("Original Slice:", slice)

    // Set the number of CPU cores to be used
    runtime.GOMAXPROCS(15)

    var wg sync.WaitGroup
    wg.Add(1)
    start := time.Now()
    concurrentQuickSort(slice, &wg)
    wg.Wait()
    duration := time.Since(start)

    fmt.Println("Sorted Slice:", slice)
    fmt.Println("Execution Time:", duration)
}
