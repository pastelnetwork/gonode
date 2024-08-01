package cloud

import (
	"bytes"
	"fmt"

	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pastelnetwork/gonode/common/log"
)

type Storage interface {
	Store(key string, data []byte) (string, error)
	Fetch(key string) ([]byte, error)
	StoreBatch(data [][]byte) error
	FetchBatch(keys []string) (map[string][]byte, error)
}

type RcloneStorage struct {
	bucketName string
	specName   string
}

func NewRcloneStorage(bucketName, specName string) *RcloneStorage {
	return &RcloneStorage{
		bucketName: bucketName,
		specName:   specName,
	}
}

func (r *RcloneStorage) Store(key string, data []byte) (string, error) {
	filePath := filepath.Join(os.TempDir(), key)

	// Write data to a temporary file using os.WriteFile
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write data to file: %w", err)
	}

	// Construct the remote path where the file will be stored
	// This example places the file at the root of the remote, but you can modify the path as needed
	remotePath := fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key)

	// Use rclone to copy the file to the remote
	cmd := exec.Command("rclone", "copyto", filePath, remotePath)
	if err := cmd.Run(); err != nil {
		// Clean up the local file if the upload fails
		os.Remove(filePath)
		return "", fmt.Errorf("rclone command failed: %w", err)
	}

	// Delete the local file after successful upload
	go func() {
		if err := os.Remove(filePath); err != nil {
			log.Error("failed to delete local file", "path", filePath, "error", err)
		}
	}()

	// Return the remote path where the file was stored
	return remotePath, nil
}

func (r *RcloneStorage) Fetch(key string) ([]byte, error) {
	// Construct the rclone command to fetch the file
	cmd := exec.Command("rclone", "cat", fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("rclone command failed: %w - out %s", err, out.String())
	}

	return out.Bytes(), nil
}

func (r *RcloneStorage) StoreBatch(data [][]byte) error {
	// Placeholder for StoreBatch implementation
	return nil
}

func (r *RcloneStorage) FetchBatch(keys []string) (map[string][]byte, error) {
	results := make(map[string][]byte)
	errs := make(map[string]error)
	var mu sync.Mutex

	semaphore := make(chan struct{}, 50)

	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a token

		go func(key string) {
			defer wg.Done()
			data, err := r.Fetch(key)

			func() {
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					errs[key] = err
				} else {
					results[key] = data
				}

			}()
			<-semaphore // Release the token
		}(key)
	}

	wg.Wait()

	if len(results) > 0 {
		return results, nil
	}

	if len(errs) > 0 {
		combinedError := fmt.Errorf("errors occurred in fetching keys")
		for k, e := range errs {
			combinedError = fmt.Errorf("%v; key %s error: %v", combinedError, k, e)
		}
		return nil, combinedError
	}

	return results, nil
}
