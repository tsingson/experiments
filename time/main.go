package main

import (
	"fmt"
	"time"
)

func main() {

	dtime := "2017-05-12 18:00"

	timestampp := StringToUnix(dtime)
	fmt.Println(timestampp)

}
func StringToTime(timeIn string) time.Time {
	//timeIn  := "2016-07-14 14:24"
	timestamp1, _ := time.Parse("2006-01-02 15:04", timeIn)
	return timestamp1
}

func StringToUnix(timeIn string) int64 {
	//timeIn  := "2016-07-14 14:24"
	timestamp1, _ := time.Parse("2006-01-02 15:04", timeIn)
	return timestamp1.Unix()
}
