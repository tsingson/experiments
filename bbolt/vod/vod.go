package vod

import "time"

// Config type
type Config struct {
	ID     int `storm:"id,increment"`
	Height float64
}

// Entry type
type Entry struct {
	EntryID  string `json:"movie_id" storm:"id"`
	Date     time.Time
	Calories int
	Food     string
	Config   []*Config
}
