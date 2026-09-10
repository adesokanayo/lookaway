//go:build darwin

package main

import (
	"time"

	"github.com/progrium/darwinkit/helper/action"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

type overlayScreen struct {
	window    appkit.Panel
	countdown appkit.TextField
	isPrimary bool
}

type Overlay struct {
	screens []overlayScreen
	onSkip  func()
}

func NewOverlay(onSkip func()) *Overlay {
	return &Overlay{onSkip: onSkip}
}

func (o *Overlay) Show(remain time.Duration) {
	o.Hide()

	bg := appkit.Color_ColorWithSRGBRedGreenBlueAlpha(0.04, 0.05, 0.07, 1)
	fg := appkit.Color_ColorWithSRGBRedGreenBlueAlpha(0.96, 0.97, 0.98, 1)
	muted := appkit.Color_ColorWithSRGBRedGreenBlueAlpha(0.70, 0.74, 0.78, 1)

	main := appkit.Screen_MainScreen()
	for _, screen := range appkit.Screen_Screens() {
		frame := screen.Frame()
		panel := appkit.NewPanelWithContentRectStyleMaskBackingDefer(
			frame,
			appkit.WindowStyleMaskBorderless|appkit.WindowStyleMaskNonactivatingPanel,
			appkit.BackingStoreBuffered,
			false,
		)
		objc.Retain(&panel)
		panel.SetLevel(appkit.ScreenSaverWindowLevel)
		panel.SetCollectionBehavior(
			appkit.WindowCollectionBehaviorCanJoinAllSpaces |
				appkit.WindowCollectionBehaviorFullScreenAuxiliary |
				appkit.WindowCollectionBehaviorStationary,
		)
		panel.SetOpaque(true)
		panel.SetHasShadow(false)
		panel.SetBackgroundColor(bg)
		panel.SetIgnoresMouseEvents(false)
		panel.SetHidesOnDeactivate(false)
		panel.SetAnimationBehavior(appkit.WindowAnimationBehaviorNone)
		panel.SetReleasedWhenClosed(false)

		isPrimary := sameScreen(screen, main)
		item := overlayScreen{window: panel, isPrimary: isPrimary}

		content := panel.ContentView()
		w := frame.Size.Width
		h := frame.Size.Height

		title := styledLabel("Look away", fg, 42, true)
		title.SetFrame(rectOf((w-720)/2, h*0.58, 720, 56))
		content.AddSubview(title)

		body := styledLabel("Focus on something about 20 feet away.\nA far wall, a window, or outside.", muted, 20, false)
		body.SetFrame(rectOf((w-720)/2, h*0.48, 720, 64))
		content.AddSubview(body)

		count := styledLabel(formatSeconds(remain), fg, 96, true)
		count.SetFont(appkit.Font_MonospacedDigitSystemFontOfSizeWeight(96, appkit.FontWeightSemibold))
		count.SetFrame(rectOf((w-240)/2, h*0.30, 240, 110))
		content.AddSubview(count)
		item.countdown = count

		if isPrimary {
			hint := styledLabel("Blink fully. Let the focusing muscles in your eyes rest.", muted, 16, false)
			hint.SetFrame(rectOf((w-720)/2, h*0.22, 720, 28))
			content.AddSubview(hint)

			skip := appkit.NewButtonWithTitle("Skip this break")
			skip.SetBezelStyle(appkit.BezelStyleInline)
			skip.SetFrame(rectOf((w-180)/2, h*0.14, 180, 32))
			action.Set(skip, func(sender objc.Object) {
				if o.onSkip != nil {
					o.onSkip()
				}
			})
			content.AddSubview(skip)
		}

		panel.OrderFrontRegardless()
		o.screens = append(o.screens, item)
	}

	if snd := appkit.Sound_SoundNamed("Tink"); snd.Ptr() != nil {
		snd.Play()
	}
}

func (o *Overlay) Update(remain time.Duration) {
	text := formatSeconds(remain)
	for _, screen := range o.screens {
		if screen.countdown.Ptr() != nil {
			screen.countdown.SetStringValue(text)
		}
	}
}

func (o *Overlay) Hide() {
	for i := range o.screens {
		o.screens[i].window.OrderOut(nil)
		o.screens[i].window.Close()
	}
	o.screens = nil
}

func styledLabel(text string, color appkit.Color, size float64, bold bool) appkit.TextField {
	label := appkit.NewLabel(text)
	if bold {
		label.SetFont(appkit.Font_BoldSystemFontOfSize(size))
	} else {
		label.SetFont(appkit.Font_SystemFontOfSize(size))
	}
	label.SetTextColor(color)
	label.SetAlignment(appkit.TextAlignmentCenter)
	label.SetDrawsBackground(false)
	label.SetEditable(false)
	label.SetSelectable(false)
	return label
}

func sameScreen(a, b appkit.Screen) bool {
	af := a.Frame()
	bf := b.Frame()
	return af.Origin.X == bf.Origin.X &&
		af.Origin.Y == bf.Origin.Y &&
		af.Size.Width == bf.Size.Width &&
		af.Size.Height == bf.Size.Height
}

func rectOf(x, y, width, height float64) foundation.Rect {
	return foundation.Rect{
		Origin: foundation.Point{X: x, Y: y},
		Size:   foundation.Size{Width: width, Height: height},
	}
}
