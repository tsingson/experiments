package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

func main() {
	// A huge allocation to give the GC work to do
	lotsOf := make([]*int, 15e8)
	fmt.Println("Background GC work generated")
	// Force a GC to set a baseline we can see if we set GODEBUG=gctrace=1
	runtime.GC()

	// Use up all the CPU doing work that causes allocations that could be cleaned up by the GC.
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work()
		}()
	}

	wg.Wait()

	// Make sure that this memory isn't optimised away
	runtime.KeepAlive(lotsOf)
}

func work() {
	for {
		work := make([]*int, 1e6)
		if f := factorial(20); f != 2432902008176640000 {
			fmt.Println(f)
			os.Exit(1)
		}
		runtime.KeepAlive(work)
	}
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * factorial(n-1)
}
