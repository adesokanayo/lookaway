package main

import (
	"testing"
	"time"
)

func TestEngineStartsWorkThenBreaks(t *testing.T) {
	now := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	e := NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)

	if e.Phase() != PhaseWork {
		t.Fatalf("start phase = %v", e.Phase())
	}
	if got := e.Remaining(); got != 20*time.Minute {
		t.Fatalf("remaining = %v", got)
	}

	now = now.Add(19 * time.Minute)
	if ev := e.Tick(); ev != EventTick {
		t.Fatalf("early tick = %v", ev)
	}

	now = now.Add(time.Minute)
	if ev := e.Tick(); ev != EventSeekStarted {
		t.Fatalf("expected seek start, got %v", ev)
	}
	if e.Phase() != PhaseSeek {
		t.Fatalf("phase after work = %v", e.Phase())
	}
	now = now.Add(e.Seek)
	if ev := e.Tick(); ev != EventBreakStarted {
		t.Fatalf("expected break start, got %v", ev)
	}
	if e.Phase() != PhaseBreak {
		t.Fatalf("phase after work = %v", e.Phase())
	}
	if got := e.Remaining(); got != 20*time.Second {
		t.Fatalf("break remaining = %v", got)
	}

	now = now.Add(20 * time.Second)
	if ev := e.Tick(); ev != EventBreakFinished {
		t.Fatalf("expected break end, got %v", ev)
	}
	if e.Phase() != PhaseWork {
		t.Fatalf("phase after break = %v", e.Phase())
	}
}

func TestSkipAndManualBreak(t *testing.T) {
	now := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	e := NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)

	if ev := e.StartBreak(); ev != EventBreakStarted {
		t.Fatalf("manual break = %v", ev)
	}
	if ev := e.SkipBreak(); ev != EventBreakFinished {
		t.Fatalf("skip = %v", ev)
	}
	if e.Phase() != PhaseWork {
		t.Fatalf("phase after skip = %v", e.Phase())
	}
}

func TestPauseFreezesRemaining(t *testing.T) {
	now := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	e := NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)

	now = now.Add(5 * time.Minute)
	e.Pause()
	now = now.Add(10 * time.Minute)
	if got := e.Remaining(); got != 15*time.Minute {
		t.Fatalf("paused remaining = %v", got)
	}
	if ev := e.Tick(); ev != EventTick {
		t.Fatalf("paused tick = %v", ev)
	}
	e.Resume()
	if got := e.Remaining(); got != 15*time.Minute {
		t.Fatalf("resumed remaining = %v", got)
	}
}

func TestResetWorkAfterLock(t *testing.T) {
	now := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	e := NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)

	now = now.Add(19 * time.Minute)
	if ev := e.Tick(); ev != EventTick {
		t.Fatalf("pre-lock tick = %v", ev)
	}

	now = now.Add(2 * time.Hour)
	e.ResetWork()
	if e.Phase() != PhaseWork {
		t.Fatalf("phase after unlock = %v", e.Phase())
	}
	if got := e.Remaining(); got != 20*time.Minute {
		t.Fatalf("remaining after unlock = %v", got)
	}
	if ev := e.Tick(); ev != EventTick {
		t.Fatalf("tick after unlock = %v", ev)
	}

	e.Pause()
	now = now.Add(time.Hour)
	e.ResetWork()
	if e.Phase() != PhaseWork {
		t.Fatalf("paused unlock phase = %v", e.Phase())
	}
	if got := e.Remaining(); got != 20*time.Minute {
		t.Fatalf("paused unlock remaining = %v", got)
	}
}

func TestSeekFiresOnQuietOrTimeout(t *testing.T) {
	now := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	e := NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)
	e.Seek = time.Minute

	now = now.Add(e.Work)
	if ev := e.Tick(); ev != EventSeekStarted {
		t.Fatalf("seek = %v", ev)
	}
	if ev := e.FireBreak(); ev != EventBreakStarted {
		t.Fatalf("quiet fire = %v", ev)
	}
	if e.Phase() != PhaseBreak {
		t.Fatalf("phase = %v", e.Phase())
	}

	e = NewEngine(20*time.Minute, 20*time.Second)
	e.now = func() time.Time { return now }
	e.deadline = now.Add(e.Work)
	e.Seek = time.Minute
	now = now.Add(e.Work)
	e.Tick()
	e.Pause()
	if e.Phase() != PhasePaused {
		t.Fatalf("pause during seek = %v", e.Phase())
	}
	e.Resume()
	if e.Phase() != PhaseWork {
		t.Fatalf("resume after seek pause = %v", e.Phase())
	}
	if got := e.Remaining(); got != e.Work {
		t.Fatalf("resume remaining = %v", got)
	}
}
