//go:build windows
// +build windows

package services

import (
	"syscall"
	"unsafe"
)

func getAvailableRAMInMB() (int, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	var memStatus syscall.MEMORYSTATUSEX
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))
	_, _, _ = globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))

	// Convert from bytes to MB
	return int(memStatus.AvailPhys / 1024 / 1024), nil
}
