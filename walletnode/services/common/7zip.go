package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

type FileSplitter struct {
	PartSizeMB int
}

func get7zPath() string {
	basePath := "/usr/bin/" // Adjust the base path as needed. For example: ./bin/ for local dir
	switch runtime.GOOS {
	case "windows":
		return basePath + "7za.exe"
	default:
		return basePath + "7za"
	}
}

func (fs *FileSplitter) SplitFile(filePath string) error {
	sevenZPath := get7zPath()
	cmd := exec.Command(sevenZPath, "a", "-v"+strconv.Itoa(fs.PartSizeMB)+"m", filePath+".7z", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.New("7zip split error: " + err.Error())
	}

	return nil
}

func (fs *FileSplitter) JoinFiles(dirPath string) error {
	// Find the first part of the split files
	var firstPartPath string
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".001") {
			firstPartPath = filepath.Join(dirPath, file.Name())
			break
		}
	}

	if firstPartPath == "" {
		return errors.New("no split parts (.001 file) found in the directory")
	}

	// Join the split parts
	sevenZPath := get7zPath()
	fmt.Println("Found .001 file at:", firstPartPath)
	cmd := exec.Command(sevenZPath, "x", "-aoa", "-o"+dirPath, firstPartPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Executing join command:", cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("7zip join error: %v", err)
	}

	// The assumed path of the joined archive
	reassembledPath := strings.TrimSuffix(firstPartPath, ".7z.001")
	fmt.Println("Assumed path of the joined archive:", reassembledPath)

	// Verify the existence of the joined archive
	if _, err := os.Stat(reassembledPath); os.IsNotExist(err) {
		fmt.Println("Checking for potential naming error...")
		if _, err := os.Stat(reassembledPath + ".7z"); err == nil {
			reassembledPath += ".7z"
			fmt.Println("Corrected path to joined archive:", reassembledPath)
		} else {
			return fmt.Errorf("joined archive does not exist at expected path: %s", reassembledPath)
		}
	}

	// Extract the original file from the joined .7z archive
	cmd = exec.Command(sevenZPath, "x", "-aoa", "-o"+dirPath, reassembledPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Executing extract command:", cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("7zip extract error: %v", err)
	}

	return nil
}
