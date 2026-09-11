package main

import (
	"sync"
	"time"
)

type Phase int

const (
	PhaseWork Phase = iota
	PhaseBreak
	PhasePaused
)

type Event int

const (
	EventNone Event = iota
	EventTick
	EventBreakStarted
	EventBreakFinished
)

// Engine is the 20-20-20 timer.
type Engine struct {
	Work  time.Duration
	Break time.Duration

	mu          sync.Mutex
	now         func() time.Time
	phase       Phase
	priorPhase  Phase
	deadline    time.Time
	savedRemain time.Duration
}

func NewEngine(work, rest time.Duration) *Engine {
	e := &Engine{
		Work:  work,
		Break: rest,
		now:   time.Now,
		phase: PhaseWork,
	}
	e.deadline = e.now().Add(work)
	return e
}

func (e *Engine) Phase() Phase {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.phase
}

func (e *Engine) Remaining() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase == PhasePaused {
		return e.savedRemain
	}
	remain := e.deadline.Sub(e.now())
	if remain < 0 {
		return 0
	}
	return remain
}

func (e *Engine) Tick() Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase == PhasePaused {
		return EventTick
	}
	if e.now().Before(e.deadline) {
		return EventTick
	}
	switch e.phase {
	case PhaseWork:
		e.enterBreakLocked()
		return EventBreakStarted
	case PhaseBreak:
		e.enterWorkLocked()
		return EventBreakFinished
	default:
		return EventNone
	}
}

func (e *Engine) StartBreak() Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase == PhaseBreak {
		return EventNone
	}
	e.enterBreakLocked()
	return EventBreakStarted
}

func (e *Engine) SkipBreak() Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase != PhaseBreak {
		return EventNone
	}
	e.enterWorkLocked()
	return EventBreakFinished
}

func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase == PhasePaused {
		return
	}
	e.savedRemain = e.deadline.Sub(e.now())
	if e.savedRemain < 0 {
		e.savedRemain = 0
	}
	e.priorPhase = e.phase
	e.phase = PhasePaused
}

func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase != PhasePaused {
		return
	}
	e.phase = e.priorPhase
	e.deadline = e.now().Add(e.savedRemain)
}

// ResetWork starts a fresh work interval. Used after unlock or wake
// so lock-screen time does not count as looking at the display.
func (e *Engine) ResetWork() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.priorPhase = PhaseWork
	e.savedRemain = 0
	e.enterWorkLocked()
}

func (e *Engine) enterBreakLocked() {
	e.phase = PhaseBreak
	e.deadline = e.now().Add(e.Break)
}

func (e *Engine) enterWorkLocked() {
	e.phase = PhaseWork
	e.deadline = e.now().Add(e.Work)
}
