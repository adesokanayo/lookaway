//go:build windows

package main

import (
	"time"
	"unsafe"
)

var (
	procGetLastInputInfo = user32.NewProc("GetLastInputInfo")
	procGetTickCount     = kernel32.NewProc("GetTickCount")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

type lastInputInfo struct {
	size uint32
	time uint32
}

const (
	vkLButton = 0x01
	vkRButton = 0x02
	vkMButton = 0x04
	vkShift   = 0x10
	vkControl = 0x11
	vkLWin    = 0x5B
	vkRWin    = 0x5C
	keyDown   = 0x8000
)

func keyHeld(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return r&keyDown != 0
}

func inputIsQuiet() bool {
	info := lastInputInfo{}
	info.size = uint32(unsafe.Sizeof(info))
	ok, _, _ := procGetLastInputInfo.Call(uintptr(unsafe.Pointer(&info)))
	if ok == 0 {
		return false
	}
	tick, _, _ := procGetTickCount.Call()
	idle := time.Duration(uint32(tick)-info.time) * time.Millisecond
	if idle < quietIdle {
		return false
	}
	if keyHeld(vkLButton) || keyHeld(vkRButton) || keyHeld(vkMButton) {
		return false
	}
	if keyHeld(vkControl) || keyHeld(vkShift) || keyHeld(vkLWin) || keyHeld(vkRWin) {
		return false
	}
	return true
}
