package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/oklog/run"
	"github.com/pkg/errors"
)

func main() {
	var g run.Group
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			fmt.Println("task 01 start ")
			select {
			case <-time.After(time.Second):
				fmt.Printf("-----------------The first actor had its time elapsed\n")
				//	return nil
				return errors.New("immediate teardown")
			case <-cancel:
				time.Sleep(5 * time.Second)
				fmt.Printf("The first actor was canceled\n")
				return nil

			}
		}, func(err error) {
			time.Sleep(3 * time.Second)
			fmt.Printf(" 0000000 The first actor was interrupted with: %v\n", err)

			close(cancel)
		})
	}
	{
		g.Add(func() error {
			fmt.Println("task 02 start ")
			fmt.Printf("----------- The second actor is returning immediately\n")
			return errors.New("immediate teardown")
		}, func(err error) {
			// Note that this interrupt function is called, even though the
			// corresponding execute function has already returned.
			fmt.Printf(" ************ The second actor was interrupted with: %v\n", err)
		})
	}
	for i := 0; i < 10; i++ {
		var j int
		j = i
		g.Add(func() error {
			id := "task" + strconv.Itoa(j) + "start"
			time.Sleep(1 * time.Second)
			fmt.Println(id)
			fmt.Printf("----------- The second actor is returning immediately\n")
			// return errors.New("immediate teardown")
			return nil
		}, func(err error) {
			// Note that this interrupt function is called, even though the
			// corresponding execute function has already returned.
			//	fmt.Printf(" ************ The second actor was interrupted with: %v\n", err)
			fmt.Println("done" + strconv.Itoa(j))
		})

	}
	result := g.Run()
	fmt.Printf("The group was terminated with: %v\n", result)
}
