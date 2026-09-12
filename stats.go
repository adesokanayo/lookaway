package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type dayCounts struct {
	Shown     int `json:"shown"`
	Completed int `json:"completed"`
	Skipped   int `json:"skipped"`
}

type statsFile struct {
	WelcomeShown bool                 `json:"welcome_shown"`
	Days         map[string]dayCounts `json:"days"`
}

var (
	statsMu   sync.Mutex
	statsPath = ""
)

func statsFilePath() (string, error) {
	if statsPath != "" {
		return statsPath, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "Lookaway")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "stats.json"), nil
}

func loadStats() statsFile {
	s := statsFile{Days: map[string]dayCounts{}}
	path, err := statsFilePath()
	if err != nil {
		return s
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	if json.Unmarshal(raw, &s) != nil || s.Days == nil {
		s.Days = map[string]dayCounts{}
	}
	return s
}

func saveStats(s statsFile) {
	path, err := statsFilePath()
	if err != nil {
		return
	}
	cutoff := time.Now().AddDate(0, 0, -40).Format("2006-01-02")
	for day := range s.Days {
		if day < cutoff {
			delete(s.Days, day)
		}
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, raw, 0o644)
}

func todayKey() string {
	return time.Now().Format("2006-01-02")
}

func needsWelcome() bool {
	statsMu.Lock()
	defer statsMu.Unlock()
	return !loadStats().WelcomeShown
}

func markWelcome() {
	statsMu.Lock()
	defer statsMu.Unlock()
	s := loadStats()
	s.WelcomeShown = true
	saveStats(s)
}

func recordShown() {
	bumpToday(func(d dayCounts) dayCounts {
		d.Shown++
		return d
	})
}

func recordCompleted() {
	bumpToday(func(d dayCounts) dayCounts {
		d.Completed++
		return d
	})
}

func recordSkipped() {
	bumpToday(func(d dayCounts) dayCounts {
		d.Skipped++
		return d
	})
}

func bumpToday(fn func(dayCounts) dayCounts) {
	statsMu.Lock()
	defer statsMu.Unlock()
	s := loadStats()
	key := todayKey()
	s.Days[key] = fn(s.Days[key])
	saveStats(s)
}

func todayCounts() (completed, skipped int) {
	statsMu.Lock()
	defer statsMu.Unlock()
	d := loadStats().Days[todayKey()]
	return d.Completed, d.Skipped
}
