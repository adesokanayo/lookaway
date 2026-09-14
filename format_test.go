package main

import (
	"testing"
	"time"
)

func TestFormatClock(t *testing.T) {
	if got := formatClock(20 * time.Minute); got != "20:00" {
		t.Fatalf("got %q", got)
	}
	if got := formatClock(19*time.Minute + 4*time.Second); got != "19:04" {
		t.Fatalf("got %q", got)
	}
	if got := formatSeconds(20 * time.Second); got != "20" {
		t.Fatalf("got %q", got)
	}
	if got := warningSecond(5 * time.Second); got != 5 {
		t.Fatalf("warn 5s = %d", got)
	}
	if got := warningSecond(4900 * time.Millisecond); got != 5 {
		t.Fatalf("warn 4.9s = %d", got)
	}
	if got := warningSecond(1500 * time.Millisecond); got != 2 {
		t.Fatalf("warn 1.5s = %d", got)
	}
	if got := warningSecond(6 * time.Second); got != 0 {
		t.Fatalf("warn 6s = %d", got)
	}
}
