package main

import (
	"path/filepath"
	"testing"
)

func TestTodayStatsAndWelcome(t *testing.T) {
	dir := t.TempDir()
	statsPath = filepath.Join(dir, "stats.json")
	t.Cleanup(func() { statsPath = "" })

	if !needsWelcome() {
		t.Fatal("expected first-run welcome")
	}
	markWelcome()
	if needsWelcome() {
		t.Fatal("welcome should stay dismissed")
	}

	recordShown()
	recordCompleted()
	recordShown()
	recordSkipped()

	got, skip := todayCounts()
	if got != 1 || skip != 1 {
		t.Fatalf("today = %d completed, %d skipped", got, skip)
	}
	if s := formatTodayStats(got, skip); s != "Today 1 lookaway · 1 skipped" {
		t.Fatalf("format = %q", s)
	}
	if s := formatTodayStats(3, 0); s != "Today 3 lookaways" {
		t.Fatalf("format = %q", s)
	}
}
