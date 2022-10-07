//go:build windows
// +build windows

package utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// DiskUsage returns free space in bytes of current working directory
func DiskUsage(dir string) (DiskStatus, error) {
	var disk DiskStatus

	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	var freeBytes int64
	var totalBytes int64
	var notUsed int64

	_, _, err := c.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(dir))),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&notUsed)))

	if err != nil {
		return disk, err
	}

	disk.All = BytesToMB(uint64(totalBytes))
	disk.Free = BytesToMB(uint64(freeBytes))
	disk.Used = BytesToMB(uint64(totalBytes - freeBytes))

	return disk, nil
}
