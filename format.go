package main

import (
	"fmt"
	"time"
)

func formatClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Round(time.Second) / time.Second)
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

func formatSeconds(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	return fmt.Sprintf("%d", int(d.Round(time.Second)/time.Second))
}
