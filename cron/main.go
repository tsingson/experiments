package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
)

type TestJob struct {
}

func (this TestJob) Run() {
	fmt.Println("testJob1...", time.Now().String())
}

type Test2Job struct {
}

func (this Test2Job) Run() {
	fmt.Println("testJob2...", time.Now().String())
}

// 启动多个任务
func main() {
	i := 0
	c := cron.New()

	// AddFunc
	spec := "0 30 * * * *"

	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})

	// AddJob方法
	c.AddJob(spec, TestJob{})
	c.AddJob(spec, Test2Job{})

	// 启动计划任务
	c.Start()

	// 关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
