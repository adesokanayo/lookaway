//go:build darwin

package main

/*
#cgo CFLAGS: -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics
#include <stdlib.h>
#include "lookaway_cocoa_darwin.h"
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

var cocoaApp *App

func main() {
	runtime.LockOSThread()
	C.LookawayRunApp()
}

func cstr(s string) *C.char {
	return C.CString(s)
}

func free(p *C.char) { C.free(unsafe.Pointer(p)) }

type Overlay struct{}

func NewOverlay(func()) *Overlay { return &Overlay{} }

func (o *Overlay) Show(remain time.Duration) {
	s := cstr(formatSeconds(remain))
	defer free(s)
	C.LookawayShowOverlay(s)
}

func (o *Overlay) Update(remain time.Duration) {
	s := cstr(formatSeconds(remain))
	defer free(s)
	C.LookawayUpdateOverlay(s)
}

func (o *Overlay) Hide() { C.LookawayHideOverlay() }

type App struct {
	engine  *Engine
	overlay *Overlay
	locked  bool
}

//export lookawayOnReady
func lookawayOnReady() {
	work, rest := loadIntervals()
	cocoaApp = &App{engine: NewEngine(work, rest), overlay: NewOverlay(nil)}
	cocoaApp.refresh()
	if needsWelcome() {
		C.LookawayShowWelcome()
	}
	go func() {
		t := time.NewTicker(250 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			C.LookawayDispatchTick()
		}
	}()
}

//export lookawayOnTick
func lookawayOnTick() {
	a := cocoaApp
	if a == nil {
		return
	}
	if C.LookawayScreenIsLocked() {
		a.locked = true
		return
	}
	if a.locked {
		a.onUnlocked()
		return
	}
	ev := nextEvent(a.engine)
	switch ev {
	case EventBreakStarted:
		recordShown()
		a.overlay.Show(a.engine.Remaining())
	case EventBreakFinished:
		recordCompleted()
		a.overlay.Hide()
	}
	if a.engine.Phase() == PhaseBreak {
		a.overlay.Update(a.engine.Remaining())
	}
	a.refresh()
}

//export lookawayOnBreakNow
func lookawayOnBreakNow() {
	a := cocoaApp
	if a != nil && a.engine.StartBreak() == EventBreakStarted {
		recordShown()
		a.overlay.Show(a.engine.Remaining())
		a.refresh()
	}
}

//export lookawayOnPause
func lookawayOnPause() {
	a := cocoaApp
	if a == nil {
		return
	}
	if a.engine.Phase() == PhasePaused {
		a.engine.Resume()
	} else {
		a.engine.Pause()
	}
	a.refresh()
}

//export lookawayOnSkip
func lookawayOnSkip() {
	a := cocoaApp
	if a != nil && a.engine.SkipBreak() == EventBreakFinished {
		recordSkipped()
		a.overlay.Hide()
		a.refresh()
	}
}

//export lookawayOnQuit
func lookawayOnQuit() {
	C.LookawayQuit()
}

//export lookawayOnUnlocked
func lookawayOnUnlocked() {
	a := cocoaApp
	if a == nil {
		return
	}
	a.onUnlocked()
}

//export lookawayOnWelcomeDone
func lookawayOnWelcomeDone() {
	markWelcome()
	C.LookawayHideWelcome()
}

func (a *App) onUnlocked() {
	a.locked = false
	a.engine.ResetWork()
	a.overlay.Hide()
	a.refresh()
}

func inputIsQuiet() bool {
	return bool(C.LookawayInputIsQuiet())
}

func (a *App) refresh() {
	remain := a.engine.Remaining()
	status := formatClock(remain)
	next := "Next break " + formatClock(remain)
	pause := "Pause"
	warn, flash := false, false
	switch a.engine.Phase() {
	case PhaseWork:
		if w := warningSecond(remain); w > 0 {
			status = formatSeconds(time.Duration(w) * time.Second)
			next = "Break soon"
			warn = true
			flash = time.Now().UnixMilli()/400%2 == 0
		}
	case PhaseSeek:
		status = "…"
		next = "Waiting for a pause"
	case PhaseBreak:
		status = formatSeconds(remain) + "s"
		next = "Break " + formatSeconds(remain) + "s"
	case PhasePaused:
		status = "Paused"
		next = "Paused " + formatClock(remain)
		pause = "Resume"
	}
	C.LookawaySetEye(C.bool(warn), C.bool(flash))
	completed, skipped := todayCounts()
	stats := formatTodayStats(completed, skipped)
	cs, cn, ct, cp := cstr(status), cstr(next), cstr(stats), cstr(pause)
	defer free(cs)
	defer free(cn)
	defer free(ct)
	defer free(cp)
	C.LookawaySetMenu(cs, cn, ct, cp)
}
