package main

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	for _, size := range []int{-10, -1, 0, 1, 7, 8, 9, 1000} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			data := generateRandomElements(size)
			wantSize := size
			if wantSize < 0 {
				wantSize = 0
			}
			if len(data) != wantSize {
				t.Fatalf("length = %d, want %d", len(data), wantSize)
			}
			for i, value := range data {
				if value <= 0 {
					t.Fatalf("data[%d] = %d, want a positive integer", i, value)
				}
			}
		})
	}
}

func TestMaximumFunctions(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	tests := []struct {
		name string
		data []int
		want int
	}{
		{"nil", nil, 0},
		{"empty", []int{}, 0},
		{"single positive", []int{42}, 42},
		{"single negative", []int{-42}, -42},
		{"zero", []int{0}, 0},
		{"zeros", []int{0, 0, 0}, 0},
		{"equal values", []int{5, 5, 5, 5}, 5},
		{"maximum first", []int{10, 3, 2, 1}, 10},
		{"maximum middle", []int{1, 2, 10, 3, 4}, 10},
		{"maximum last", []int{1, 2, 3, 10}, 10},
		{"eight elements", []int{1, 2, 3, 8, 5, 6, 7, 4}, 8},
		{"remainder", []int{1, 2, 3, 4, 5, 6, 7, 8, 99}, 99},
		{"negative remainder", []int{-9, -8, -7, -6, -5, -4, -3, -2, -1}, -1},
		{"negative values", []int{-9, -2, -5}, -2},
		{"mixed values", []int{-2, 0, 3, -1}, 3},
		{"integer limits", []int{minInt, 0, maxInt}, maxInt},
		{"minimum integer", []int{minInt}, minInt},
	}
	functions := []struct {
		name string
		find func([]int) int
	}{
		{"maximum", maximum},
		{"maxChunks", maxChunks},
	}

	for _, function := range functions {
		t.Run(function.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					before := make([]int, len(tt.data))
					copy(before, tt.data)
					if got := function.find(tt.data); got != tt.want {
						t.Errorf("maximum = %d, want %d", got, tt.want)
					}
					if len(tt.data) > 0 && !reflect.DeepEqual(tt.data, before) {
						t.Error("input slice was modified")
					}
				})
			}
		})
	}
}

func TestMaxChunksMatchesMaximum(t *testing.T) {
	// Cover all remainders and sizes below, equal to, and above CHUNKS.
	for size := 0; size <= CHUNKS*4+1; size++ {
		data := make([]int, size)
		for i := range data {
			data[i] = -100
		}
		for position := range data {
			data[position] = -1
			if got, want := maxChunks(data), maximum(data); got != want {
				t.Fatalf("size %d, position %d: got %d, want %d", size, position, got, want)
			}
			data[position] = -100
		}
	}

	data := generateRandomElements(100_003)
	if got, want := maxChunks(data), maximum(data); got != want {
		t.Fatalf("random data: got %d, want %d", got, want)
	}
}
