//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")
	winmm    = windows.NewLazySystemDLL("winmm.dll")

	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procPostMessageW        = user32.NewProc("PostMessageW")
	procDestroyWindow       = user32.NewProc("DestroyWindow")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procInvalidateRect      = user32.NewProc("InvalidateRect")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procSetFocus            = user32.NewProc("SetFocus")
	procGetClientRect       = user32.NewProc("GetClientRect")
	procBeginPaint          = user32.NewProc("BeginPaint")
	procEndPaint            = user32.NewProc("EndPaint")
	procFillRect            = user32.NewProc("FillRect")
	procDrawTextW           = user32.NewProc("DrawTextW")
	procLoadCursorW         = user32.NewProc("LoadCursorW")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procSetCursor           = user32.NewProc("SetCursor")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procSetProcessDPIAware  = user32.NewProc("SetProcessDPIAware")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")

	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procCreateFontW      = gdi32.NewProc("CreateFontW")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procSetBkMode        = gdi32.NewProc("SetBkMode")
	procSetTextColor     = gdi32.NewProc("SetTextColor")

	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
	procPlaySoundW       = winmm.NewProc("PlaySoundW")
)

const (
	wsPopup        = 0x80000000
	wsVisible      = 0x10000000
	wsExTopmost    = 0x00000008
	wsExToolwindow = 0x00000080
	hwndTopmost    = ^uintptr(0) // -1
	swShow         = 5
	swpShowWindow  = 0x0040
	swpNoMove      = 0x0002
	swpNoSize      = 0x0001

	smXVirtualScreen  = 76
	smYVirtualScreen  = 77
	smCXVirtualScreen = 78
	smCYVirtualScreen = 79

	wmDestroy    = 0x0002
	wmClose      = 0x0010
	wmPaint      = 0x000F
	wmCommand    = 0x0111
	wmLButtonUp  = 0x0202
	wmRButtonUp  = 0x0205
	wmKeyDown    = 0x0100
	wmSetCursor  = 0x0020
	wmApp        = 0x8000
	wmTray       = wmApp + 1
	wmAppRefresh = wmApp + 2
	wmNull       = 0x0000

	nimAdd     = 0
	nimModify  = 1
	nimDelete  = 2
	nifMessage = 0x00000001
	nifIcon    = 0x00000002
	nifTip     = 0x00000004

	mfString       = 0x00000000
	mfGrayed       = 0x00000001
	mfSeparator    = 0x00000800
	tpmRightButton = 0x0002

	idiApplication = 32512
	idcArrow       = 32512
	idcHand        = 32649

	dtCenter     = 0x00000001
	dtVCenter    = 0x00000004
	dtSingleLine = 0x00000020
	dtWordBreak  = 0x00000010
	transparent  = 1

	vkEscape = 0x1B

	sndAsync = 0x0001
	sndAlias = 0x00010000

	idStartBreak = 1001
	idPause      = 1002
	idSkip       = 1003
	idQuit       = 1004
)

type point struct {
	x, y int32
}

type rect struct {
	left, top, right, bottom int32
}

type msg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
}

type paintStruct struct {
	hdc         uintptr
	erase       int32
	rcPaint     rect
	restore     int32
	incUpdate   int32
	rgbReserved [32]byte
}

type wndClassEx struct {
	size       uint32
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   uintptr
	icon       uintptr
	cursor     uintptr
	background uintptr
	menuName   *uint16
	className  *uint16
	iconSm     uintptr
}

type notifyIconData struct {
	size            uint32
	wnd             uintptr
	id              uint32
	flags           uint32
	callbackMessage uint32
	icon            uintptr
	tip             [128]uint16
}

func utf16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		p, _ = windows.UTF16PtrFromString("")
	}
	return p
}

func loWord(v uintptr) int32 { return int32(uint16(v)) }
func hiWord(v uintptr) int32 { return int32(uint16(v >> 16)) }

func rgb(r, g, b uint8) uint32 {
	return uint32(r) | uint32(g)<<8 | uint32(b)<<16
}

func getModuleHandle() uintptr {
	h, _, _ := procGetModuleHandleW.Call(0)
	return h
}

func setProcessDPIAware() {
	procSetProcessDPIAware.Call()
}

func loadAppIcon() uintptr {
	h, _, _ := procLoadIconW.Call(0, uintptr(idiApplication))
	return h
}

func loadCursor(id uintptr) uintptr {
	h, _, _ := procLoadCursorW.Call(0, id)
	return h
}

func registerClass(name string, proc uintptr, bg uintptr) {
	cls := wndClassEx{
		wndProc:    proc,
		instance:   getModuleHandle(),
		icon:       loadAppIcon(),
		cursor:     loadCursor(idcArrow),
		background: bg,
		className:  utf16Ptr(name),
	}
	cls.size = uint32(unsafe.Sizeof(cls))
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&cls)))
}

func createWindow(exStyle uint32, class, title string, style uint32, x, y, w, h int32) uintptr {
	hwnd, _, _ := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(utf16Ptr(class))),
		uintptr(unsafe.Pointer(utf16Ptr(title))),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(w),
		uintptr(h),
		0, 0, getModuleHandle(), 0,
	)
	return hwnd
}

func destroyWindow(hwnd uintptr) {
	if hwnd != 0 {
		procDestroyWindow.Call(hwnd)
	}
}

func showWindow(hwnd uintptr) {
	procShowWindow.Call(hwnd, swShow)
	procUpdateWindow.Call(hwnd)
	procSetWindowPos.Call(hwnd, hwndTopmost, 0, 0, 0, 0, swpNoMove|swpNoSize|swpShowWindow)
	procSetForegroundWindow.Call(hwnd)
	procSetFocus.Call(hwnd)
}

func invalidate(hwnd uintptr) {
	procInvalidateRect.Call(hwnd, 0, 1)
}

func postMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) {
	procPostMessageW.Call(hwnd, uintptr(msg), wParam, lParam)
}

func postQuit() { procPostQuitMessage.Call(0) }

func messageLoop() {
	var m msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func virtualScreen() (x, y, w, h int32) {
	x = int32(sysMetric(smXVirtualScreen))
	y = int32(sysMetric(smYVirtualScreen))
	w = int32(sysMetric(smCXVirtualScreen))
	h = int32(sysMetric(smCYVirtualScreen))
	return
}

func sysMetric(idx int) uintptr {
	v, _, _ := procGetSystemMetrics.Call(uintptr(idx))
	return v
}

func createFont(height int32, bold bool) uintptr {
	weight := uintptr(400)
	if bold {
		weight = 700
	}
	h, _, _ := procCreateFontW.Call(
		uintptr(height), 0, 0, 0, weight, 0, 0, 0,
		1, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(utf16Ptr("Segoe UI"))),
	)
	return h
}

func playBreakSound() {
	procPlaySoundW.Call(
		uintptr(unsafe.Pointer(utf16Ptr("SystemAsterisk"))),
		0,
		sndAsync|sndAlias,
	)
}

func notifyIcon(action uint32, nid *notifyIconData) {
	nid.size = uint32(unsafe.Sizeof(*nid))
	procShellNotifyIconW.Call(uintptr(action), uintptr(unsafe.Pointer(nid)))
}

func setTip(nid *notifyIconData, tip string) {
	u, _ := windows.UTF16FromString(tip)
	copy(nid.tip[:], u)
}

func appendMenu(menu uintptr, flags uint32, id uintptr, title string) {
	procAppendMenuW.Call(menu, uintptr(flags), id, uintptr(unsafe.Pointer(utf16Ptr(title))))
}

func defWindowProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return r
}
