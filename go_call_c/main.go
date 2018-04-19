package main

import "fmt"

/*
// Note the -std=gnu99. Using -std=c99 will not work.
#cgo CFLAGS: -std=gnu99
#include <stdint.h>

void cMultiply(int len, uint32_t *input, uint32_t *output) {
    for (int i = 0; i < len; i++) {
        output[i] = input[i] * 2;
    }
}
*/
import "C"

func multiply(input []uint32) []uint32 {
	output := make([]uint32, len(input))
	C.cMultiply(C.int(len(input)),
		(*C.uint32_t)(&input[0]),
		(*C.uint32_t)(&output[0]))
	return output
}

func main() {
	list := []uint32{23, 42, 17}
	list = multiply(list)
	for idx, val := range list {
		fmt.Printf("index %d: value %d\n", idx, val)
	}
}
