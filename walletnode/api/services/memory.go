//go:build !windows
// +build !windows

package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getAvailableRAMInMB() (int, error) {
	var cmdString string
	if runtime.GOOS == "linux" {
		cmdString = "free -m | grep Mem | awk '{print $7}'"
	} else if runtime.GOOS == "darwin" {
		cmdString = "vm_stat | grep 'Pages free' | awk '{print $3}' | sed 's/\\.$//'"
	} else {
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd := exec.Command("sh", "-c", cmdString)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to execute command: %w", err)
	}

	output := strings.TrimSpace(out.String())
	availableRAM, err := strconv.Atoi(output)
	if err != nil {
		return 0, fmt.Errorf("failed to convert output to integer: %w", err)
	}

	if runtime.GOOS == "darwin" {
		// Convert pages to MB (assuming 4096 bytes per page)
		availableRAM = availableRAM * 4096 / 1024 / 1024
	}

	return availableRAM, nil
}
