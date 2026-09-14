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
	PhaseSeek
)

const defaultSeek = 60 * time.Second
const quietIdle = 2 * time.Second

type Event int

const (
	EventNone Event = iota
	EventTick
	EventBreakStarted
	EventBreakFinished
	EventSeekStarted
)

// Engine is the 20-20-20 timer.
type Engine struct {
	Work  time.Duration
	Break time.Duration
	Seek  time.Duration

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
		Seek:  defaultSeek,
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
		e.enterSeekLocked()
		return EventSeekStarted
	case PhaseSeek:
		e.enterBreakLocked()
		return EventBreakStarted
	case PhaseBreak:
		e.enterWorkLocked()
		return EventBreakFinished
	default:
		return EventNone
	}
}

// FireBreak starts the overlay during the quiet-seek window.
func (e *Engine) FireBreak() Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.phase != PhaseSeek {
		return EventNone
	}
	e.enterBreakLocked()
	return EventBreakStarted
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
	if e.phase == PhaseSeek {
		e.priorPhase = PhaseWork
		e.savedRemain = e.Work
		e.phase = PhasePaused
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

func (e *Engine) enterSeekLocked() {
	e.phase = PhaseSeek
	seek := e.Seek
	if seek <= 0 {
		seek = defaultSeek
	}
	e.deadline = e.now().Add(seek)
}

func (e *Engine) enterBreakLocked() {
	e.phase = PhaseBreak
	e.deadline = e.now().Add(e.Break)
}

func (e *Engine) enterWorkLocked() {
	e.phase = PhaseWork
	e.deadline = e.now().Add(e.Work)
}

func nextEvent(e *Engine) Event {
	ev := e.Tick()
	if e.Phase() == PhaseSeek && inputIsQuiet() {
		if fired := e.FireBreak(); fired == EventBreakStarted {
			return fired
		}
	}
	return ev
}
