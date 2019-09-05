package main

import (
	"github.com/tsingson/experiments/interfacce"
)

func main() {
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i
	}

	interface_test.NumbersToStringsBadInstrumented(numbers)
}
