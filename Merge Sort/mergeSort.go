package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
    "runtime"
)

func generateRandomSlice(size int) []int {
    slice := make([]int, size)
    for i := 0; i < size; i++ {
        slice[i] = rand.Intn(100) // Random integers between 0 and 99
    }
    return slice
}


func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    for len(left) > 0 || len(right) > 0 {
        if len(left) == 0 {
            return append(result, right...)
        }
        if len(right) == 0 {
            return append(result, left...)
        }
        if left[0] <= right[0] {
            result = append(result, left[0])
            left = left[1:]
        } else {
            result = append(result, right[0])
            right = right[1:]
        }
    }
    return result
}

func concurrentMergeSort(slice []int, wg *sync.WaitGroup) []int {
    defer wg.Done()
    if len(slice) < 2 {
        return slice
    }
    mid := len(slice) / 2
    var left, right []int
    var leftWg, rightWg sync.WaitGroup
    leftWg.Add(1)
    go func() {
        left = concurrentMergeSort(slice[:mid], &leftWg)
    }()
    rightWg.Add(1)
    go func() {
        right = concurrentMergeSort(slice[mid:], &rightWg)
    }()
    leftWg.Wait()
    rightWg.Wait()
    return merge(left, right)
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
    sortedSlice := concurrentMergeSort(slice, &wg)
    wg.Wait()
    duration := time.Since(start)

    fmt.Println("Sorted Slice:", sortedSlice)
    fmt.Println("Execution Time:", duration)
}