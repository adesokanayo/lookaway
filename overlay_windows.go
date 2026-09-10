//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"
)

var overlayWndProc = syscall.NewCallback(overlayWindowProc)

type Overlay struct {
	hwnd     uintptr
	brush    uintptr
	onSkip   func()
	remain   time.Duration
	skipRect rect
}

func NewOverlay(onSkip func()) *Overlay {
	brush, _, _ := procCreateSolidBrush.Call(uintptr(rgb(10, 13, 18)))
	registerClass("LookawayOverlay", overlayWndProc, brush)
	return &Overlay{onSkip: onSkip, brush: brush}
}

func (o *Overlay) Show(remain time.Duration) {
	o.Hide()
	o.remain = remain
	x, y, w, h := virtualScreen()
	if w <= 0 || h <= 0 {
		w, h = 1920, 1080
	}
	hwnd := createWindow(wsExTopmost|wsExToolwindow, "LookawayOverlay", "Lookaway", wsPopup|wsVisible, x, y, w, h)
	o.hwnd = hwnd
	showWindow(hwnd)
	playBreakSound()
}

func (o *Overlay) Update(remain time.Duration) {
	o.remain = remain
	if o.hwnd != 0 {
		invalidate(o.hwnd)
	}
}

func (o *Overlay) Hide() {
	if o.hwnd != 0 {
		destroyWindow(o.hwnd)
		o.hwnd = 0
	}
}

func overlayWindowProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	o := winApp.overlay
	switch msg {
	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		var rc rect
		procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&rc)), o.brush)
		drawOverlay(hdc, rc, o)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case wmLButtonUp:
		x, y := loWord(lParam), hiWord(lParam)
		if hit(o.skipRect, x, y) && o.onSkip != nil {
			o.onSkip()
		}
		return 0
	case wmKeyDown:
		if wParam == vkEscape && o.onSkip != nil {
			o.onSkip()
		}
		return 0
	case wmClose:
		if o.onSkip != nil {
			o.onSkip()
		}
		return 0
	case wmDestroy:
		return 0
	}
	return defWindowProc(hwnd, msg, wParam, lParam)
}

func drawOverlay(hdc uintptr, rc rect, o *Overlay) {
	w := rc.right - rc.left
	h := rc.bottom - rc.top
	procSetBkMode.Call(hdc, transparent)

	titleFont := createFont(-48, true)
	bodyFont := createFont(-22, false)
	countFont := createFont(-96, true)
	hintFont := createFont(-18, false)
	defer func() {
		procDeleteObject.Call(titleFont)
		procDeleteObject.Call(bodyFont)
		procDeleteObject.Call(countFont)
		procDeleteObject.Call(hintFont)
	}()

	procSetTextColor.Call(hdc, uintptr(rgb(245, 247, 250)))
	drawCentered(hdc, titleFont, "Look away", 0, h/6, w, h/8)

	procSetTextColor.Call(hdc, uintptr(rgb(178, 188, 199)))
	drawCentered(hdc, bodyFont, "Look about 20 feet away.", 0, h/3, w, h/8)

	procSetTextColor.Call(hdc, uintptr(rgb(245, 247, 250)))
	drawCentered(hdc, countFont, formatSeconds(o.remain), 0, h/2-40, w, h/6)

	skip := rect{
		left:   w/2 - 110,
		top:    (h * 17) / 20,
		right:  w/2 + 110,
		bottom: (h*17)/20 + 40,
	}
	o.skipRect = skip
	drawCentered(hdc, hintFont, "Skip", skip.left, skip.top, skip.right-skip.left, skip.bottom-skip.top)
}

func drawCentered(hdc uintptr, font uintptr, text string, x, y, w, h int32) {
	old, _, _ := procSelectObject.Call(hdc, font)
	r := rect{left: x, top: y, right: x + w, bottom: y + h}
	procDrawTextW.Call(
		hdc,
		uintptr(unsafe.Pointer(utf16Ptr(text))),
		^uintptr(0), // -1
		uintptr(unsafe.Pointer(&r)),
		dtCenter|dtVCenter|dtWordBreak,
	)
	procSelectObject.Call(hdc, old)
}

func hit(r rect, x, y int32) bool {
	return x >= r.left && x <= r.right && y >= r.top && y <= r.bottom
}
