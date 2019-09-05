package interfacce

import (
	"fmt"
	"strconv"
)

type getter interface {
	get() int
}

type zero struct{}

//go:noinline
func (z zero) get() int {
	return 0
}

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func numbersToStringsBad(numbers []int) []string {
	vals := []string{}
	for _, n := range numbers {
		vals = append(vals, strconv.Itoa(n))
	}
	return vals
}

func numbersToStringsBetter(numbers []int) []string {
	vals := make([]string, 0, len(numbers))
	for _, n := range numbers {
		vals = append(vals, strconv.Itoa(n))
	}
	return vals
}

func NumbersToStringsBadInstrumented(numbers []int) []string {
	vals := []string{}
	oldCapacity := cap(vals)
	for _, n := range numbers {
		vals = append(vals, strconv.Itoa(n))
		if capacity := cap(vals); capacity != oldCapacity {
			fmt.Printf("len(vals)=%d, cap(vals)=%d (was %d)\n", len(vals), capacity, oldCapacity)
			oldCapacity = capacity
		}
	}
	return vals
}
