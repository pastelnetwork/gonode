//go:build linux
// +build linux

package utils

import (
	"golang.org/x/sys/unix"
)

// DiskUsage returns free space in bytes of current working directory
func DiskUsage(dir string) (DiskStatus, error) {
	var disk DiskStatus
	stat := unix.Statfs_t{}
	if err := unix.Statfs(dir, &stat); err != nil {
		return disk, err
	}

	disk.All = stat.Blocks * uint64(stat.Bsize)
	disk.Free = stat.Bfree * uint64(stat.Bsize)
	disk.Used = disk.All - disk.Free
	return disk, nil
}
