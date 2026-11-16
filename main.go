package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements returns size positive random integers, or an empty slice
// if size is not positive.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	data := make([]int, size)
	maxInt := int(^uint(0) >> 1)
	for i := range data {
		data[i] = rand.IntN(maxInt) + 1
	}
	return data
}

// maximum returns the largest element, or zero for an empty slice.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	maxValue := data[0]
	for _, value := range data[1:] {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

// maxChunks searches CHUNKS parts concurrently and returns their maximum.
// Each goroutine writes to its own result element; results are read after Wait.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	maxima := make([]int, CHUNKS)
	chunkSize := len(data) / CHUNKS
	var workers sync.WaitGroup
	workers.Add(CHUNKS)

	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == CHUNKS-1 {
			// Include the remainder, also when there are fewer than CHUNKS elements.
			end = len(data)
		}

		go func(index int, chunk []int) {
			defer workers.Done()
			// Empty parts must not introduce zero when all input values are negative.
			maxima[index] = data[0]
			if len(chunk) > 0 {
				maxima[index] = maximum(chunk)
			}
		}(i, data[start:end])
	}

	workers.Wait()
	return maximum(maxima)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	started := time.Now()
	maxValue := maximum(data)
	elapsed := time.Since(started).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)

	fmt.Printf("Ищем максимальное значение в %d горутинах\n", CHUNKS)
	started = time.Now()
	maxValue = maxChunks(data)
	elapsed = time.Since(started).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)
}
