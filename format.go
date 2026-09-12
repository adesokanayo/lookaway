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

func formatTodayStats(completed, skipped int) string {
	word := "lookaways"
	if completed == 1 {
		word = "lookaway"
	}
	s := fmt.Sprintf("Today %d %s", completed, word)
	if skipped > 0 {
		s += fmt.Sprintf(" · %d skipped", skipped)
	}
	return s
}

func formatSeconds(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	return fmt.Sprintf("%d", int(d.Round(time.Second)/time.Second))
}
