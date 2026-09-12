//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	winApp      *App
	trayWndProc = syscall.NewCallback(trayWindowProc)
)

type App struct {
	engine  *Engine
	overlay *Overlay
	hwnd    uintptr
	nid     notifyIconData
}

func startApp() {
	if !claimInstance() {
		return
	}
	setProcessDPIAware()
	work, rest := loadIntervals()
	a := &App{engine: NewEngine(work, rest)}
	winApp = a
	a.overlay = NewOverlay(func() {
		postMessage(a.hwnd, wmAppRefresh, idSkip, 0)
	})

	brush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(10, 13, 18)))
	registerClass("LookawayTray", trayWndProc, brush)
	a.hwnd = createWindow(0, "LookawayTray", "Lookaway", 0, 0, 0, 0, 0)

	a.nid = notifyIconData{
		wnd:             a.hwnd,
		id:              1,
		flags:           nifMessage | nifIcon | nifTip,
		callbackMessage: wmTray,
		icon:            loadAppIcon(),
	}
	setTip(&a.nid, "Lookaway")
	notifyIcon(nimAdd, &a.nid)
	registerSessionNotify(a.hwnd)
	a.refresh()

	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			ev := a.engine.Tick()
			postMessage(a.hwnd, wmAppRefresh, uintptr(ev), 0)
		}
	}()

	messageLoop()
	unregisterSessionNotify(a.hwnd)
	notifyIcon(nimDelete, &a.nid)
	a.overlay.Hide()
}

func (a *App) onUnlocked() {
	a.engine.ResetWork()
	a.overlay.Hide()
	a.refresh()
}

func (a *App) handleEvent(ev Event) {
	switch ev {
	case EventBreakStarted:
		a.overlay.Show(a.engine.Remaining())
	case EventBreakFinished:
		a.overlay.Hide()
	}
	if a.engine.Phase() == PhaseBreak {
		a.overlay.Update(a.engine.Remaining())
	}
	a.refresh()
}

func (a *App) refresh() {
	remain := a.engine.Remaining()
	tip := "Lookaway " + formatClock(remain)
	switch a.engine.Phase() {
	case PhaseBreak:
		tip = "Lookaway " + formatSeconds(remain) + "s"
	case PhasePaused:
		tip = "Lookaway paused " + formatClock(remain)
	}
	setTip(&a.nid, tip)
	notifyIcon(nimModify, &a.nid)
}

func (a *App) showMenu() {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	menu, _, _ := procCreatePopupMenu.Call()

	status := "Next " + formatClock(a.engine.Remaining())
	pause := "Pause"
	switch a.engine.Phase() {
	case PhaseBreak:
		status = "Break " + formatSeconds(a.engine.Remaining()) + "s"
	case PhasePaused:
		status = "Paused " + formatClock(a.engine.Remaining())
		pause = "Resume"
	}

	appendMenu(menu, mfString|mfGrayed, 0, status)
	appendMenu(menu, mfSeparator, 0, "")
	appendMenu(menu, mfString, idStartBreak, "Break now")
	appendMenu(menu, mfString, idPause, pause)
	appendMenu(menu, mfString, idSkip, "Skip")
	appendMenu(menu, mfSeparator, 0, "")
	appendMenu(menu, mfString, idQuit, "Quit")

	procSetForegroundWindow.Call(a.hwnd)
	procTrackPopupMenu.Call(menu, tpmRightButton, uintptr(pt.x), uintptr(pt.y), 0, a.hwnd, 0)
	postMessage(a.hwnd, wmNull, 0, 0)
	procDestroyMenu.Call(menu)
}

func (a *App) onCommand(id uintptr) {
	switch id {
	case idStartBreak:
		if a.engine.StartBreak() == EventBreakStarted {
			a.overlay.Show(a.engine.Remaining())
		}
	case idPause:
		if a.engine.Phase() == PhasePaused {
			a.engine.Resume()
		} else {
			a.engine.Pause()
		}
	case idSkip:
		if a.engine.SkipBreak() == EventBreakFinished {
			a.overlay.Hide()
		}
	case idQuit:
		a.overlay.Hide()
		notifyIcon(nimDelete, &a.nid)
		postQuit()
		return
	}
	a.refresh()
}

func trayWindowProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	a := winApp
	switch msg {
	case wmTray:
		if lParam == wmRButtonUp || lParam == wmLButtonUp {
			a.showMenu()
		}
		return 0
	case wmCommand:
		a.onCommand(wParam)
		return 0
	case wmAppRefresh:
		if wParam == idSkip {
			a.onCommand(idSkip)
			return 0
		}
		a.handleEvent(Event(wParam))
		return 0
	case wmWTSSessionChange:
		if wParam == wtsSessionUnlock {
			a.onUnlocked()
		}
		return 0
	case wmPowerBroadcast:
		if wParam == pbtAPMResumeAutomatic || wParam == pbtAPMResumeSuspend {
			a.onUnlocked()
		}
		return 0
	case wmDestroy:
		unregisterSessionNotify(hwnd)
		postQuit()
		return 0
	}
	return defWindowProc(hwnd, msg, wParam, lParam)
}
