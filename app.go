//go:build darwin

package main

import (
	"time"

	"github.com/progrium/darwinkit/dispatch"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/objc"
)

type App struct {
	engine    *Engine
	overlay   *Overlay
	status    appkit.StatusItem
	nextItem  appkit.MenuItem
	pauseItem appkit.MenuItem
}

func startApp(nsapp appkit.Application) {
	work, rest := loadIntervals()
	a := &App{engine: NewEngine(work, rest)}
	a.overlay = NewOverlay(func() {
		dispatch.MainQueue().DispatchAsync(func() {
			if a.engine.SkipBreak() == EventBreakFinished {
				a.overlay.Hide()
				a.refresh()
			}
		})
	})
	a.installMenu(nsapp)
	a.refresh()

	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			ev := a.engine.Tick()
			dispatch.MainQueue().DispatchAsync(func() {
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
			})
		}
	}()
}

func (a *App) installMenu(nsapp appkit.Application) {
	item := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(appkit.VariableStatusItemLength)
	objc.Retain(&item)
	a.status = item

	img := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("eye", "20-20-20 eye rest")
	if img.Ptr() != nil {
		img.SetTemplate(true)
		item.Button().SetImage(img)
	}
	item.Button().SetImagePosition(appkit.ImageLeft)

	menu := appkit.NewMenuWithTitle("Lookaway")
	a.nextItem = appkit.NewMenuItemWithTitleActionKeyEquivalent("Next break", objc.Selector{}, "")
	a.nextItem.SetEnabled(false)
	menu.AddItem(a.nextItem)

	menu.AddItem(appkit.MenuItem_SeparatorItem())
	menu.AddItem(appkit.NewMenuItemWithAction("Break now", "", func(sender objc.Object) {
		if a.engine.StartBreak() == EventBreakStarted {
			a.overlay.Show(a.engine.Remaining())
			a.refresh()
		}
	}))
	a.pauseItem = appkit.NewMenuItemWithAction("Pause", "", func(sender objc.Object) {
		if a.engine.Phase() == PhasePaused {
			a.engine.Resume()
		} else {
			a.engine.Pause()
		}
		a.refresh()
	})
	menu.AddItem(a.pauseItem)
	menu.AddItem(appkit.NewMenuItemWithAction("Skip", "", func(sender objc.Object) {
		if a.engine.SkipBreak() == EventBreakFinished {
			a.overlay.Hide()
			a.refresh()
		}
	}))
	menu.AddItem(appkit.MenuItem_SeparatorItem())
	menu.AddItem(appkit.NewMenuItemWithAction("Quit", "q", func(sender objc.Object) {
		a.overlay.Hide()
		nsapp.Terminate(nil)
	}))
	item.SetMenu(menu)

	nsapp.SetActivationPolicy(appkit.ApplicationActivationPolicyAccessory)
}

func (a *App) refresh() {
	remain := a.engine.Remaining()
	phase := a.engine.Phase()
	title := formatClock(remain)
	next := "Next break " + formatClock(remain)
	pause := "Pause"

	switch phase {
	case PhaseBreak:
		title = formatSeconds(remain) + "s"
		next = "Break " + formatSeconds(remain) + "s"
	case PhasePaused:
		title = "Paused"
		next = "Paused " + formatClock(remain)
		pause = "Resume"
	}

	a.status.Button().SetTitle(title)
	a.nextItem.SetTitle(next)
	a.pauseItem.SetTitle(pause)
}
