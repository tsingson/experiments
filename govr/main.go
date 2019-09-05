package main

import (
	"context"
	"fmt"
	"time"

	"github.com/orbs-network/govnr"
)

type stdoutErrorer struct{}

func (s *stdoutErrorer) Error(err error) {
	fmt.Println(err.Error())
}

func main() {
	errorHandler := &stdoutErrorer{}
	ctx, cancel := context.WithCancel(context.Background())

	data := make(chan int)
	handle := govnr.Forever(ctx, "an example process", errorHandler, func() {
		for {
			select {
			case i := <-data:
				fmt.Printf("goroutine got data: %d\n", i)
			case <-ctx.Done():
				return
			}
		}
	})

	supervisor := &govnr.TreeSupervisor{}
	supervisor.Supervise(handle)

	data <- 3
	data <- 2
	data <- 1
	cancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	supervisor.WaitUntilShutdown(shutdownCtx)
}
