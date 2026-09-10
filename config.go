package main

import (
	"os"
	"time"
)

func loadIntervals() (work, rest time.Duration) {
	work = 20 * time.Minute
	rest = 20 * time.Second
	if v := os.Getenv("LOOKAWAY_WORK"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			work = d
		}
	}
	if v := os.Getenv("LOOKAWAY_BREAK"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			rest = d
		}
	}
	return work, rest
}
