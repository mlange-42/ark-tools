//go:build windows

package app

import "syscall"

// initTimer adjusts time resolution for Windows.
// Does nothing on Unix systems.
// See https://github.com/golang/go/issues/44343
func initTimer() {
	winmmDLL := syscall.NewLazyDLL("winmm.dll")
	procTimeBeginPeriod := winmmDLL.NewProc("timeBeginPeriod")
	if _, _, err := procTimeBeginPeriod.Call(uintptr(1)); err != syscall.Errno(0) {
		panic(err)
	}
}
