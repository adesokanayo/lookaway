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
}
