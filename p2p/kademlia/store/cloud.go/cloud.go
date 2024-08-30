package cloud

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	maxConcurrentUploads = 50
)

type Storage interface {
	Store(key string, data []byte) (string, error)
	Fetch(key string) ([]byte, error)
	StoreBatch(data [][]byte) error
	FetchBatch(keys []string) (map[string][]byte, error)
	Upload(key string, data []byte) (string, error)
	UploadBatch(keys []string, data [][]byte) ([]string, error)
	CheckCloudConnection() error
	Delete(key string) error
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
		log.WithField("filePath", filePath).WithError(err).Error("Failed to write data to file")
		return "", fmt.Errorf("failed to write data to file: %w", err)
	}

	// Construct the remote path where the file will be stored
	// This example places the file at the root of the remote, but you can modify the path as needed
	remotePath := fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key)

	// Use rclone to copy the file to the remote
	cmd := exec.Command("rclone", "copyto", filePath, remotePath)
	cmdOutput := &bytes.Buffer{}
	cmd.Stderr = cmdOutput // Capture standard error
	if err := cmd.Run(); err != nil {
		log.WithField("key", key).WithError(err).Error("rclone command failed to upload file")
		return "", fmt.Errorf("rclone command failed: %s, error: %w", cmdOutput.String(), err)
	}

	// Delete the local file after successful upload
	go func() {
		if err := os.Remove(filePath); err != nil {
			log.WithField("filePath", filePath).WithError(err).Error("Failed to delete local file after upload")
		}
	}()

	// Return the remote path where the file was stored
	return remotePath, nil
}

// Upload uploads the data to the remote storage without writing to a local file
func (r *RcloneStorage) Upload(key string, data []byte) (string, error) {
	// Construct the remote path where the file will be stored
	remotePath := fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key)

	// Use rclone to copy the data to the remote
	cmd := exec.Command("rclone", "rcat", remotePath)
	cmdOutput := &bytes.Buffer{}
	cmd.Stderr = cmdOutput            // Capture standard error
	cmd.Stdin = bytes.NewReader(data) // Provide data as stdin

	if err := cmd.Run(); err != nil {
		log.WithField("remotePath", remotePath).Errorf("rclone command failed to upload file:%s", cmdOutput.String())
		return "", fmt.Errorf("rclone command failed: %s, error: %w", cmdOutput.String(), err)
	}

	// Return the remote path where the file was stored
	return remotePath, nil
}

func (r *RcloneStorage) Fetch(key string) ([]byte, error) {
	// Construct the rclone command to fetch the file
	cmd := exec.Command("rclone", "cat", fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmdOutput := &bytes.Buffer{}
	cmd.Stderr = cmdOutput // Capture standard error
	err := cmd.Run()
	if err != nil {
		log.WithField("key", key).WithError(err).Errorf("rclone command failed to fetch file:%s", cmdOutput.String())
		return nil, fmt.Errorf("rclone command failed: %s, error: %w", cmdOutput.String(), err)
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

func (r *RcloneStorage) UploadBatch(keys []string, data [][]byte) ([]string, error) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentUploads) // Semaphore to limit concurrent goroutines

	var mu sync.Mutex
	var lastError error
	successfulKeys := []string{} // Slice to store successfully uploaded keys

	for i := range data {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(key string, data []byte) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			if _, err := r.Upload(key, data); err != nil {
				func() {
					mu.Lock()
					defer mu.Unlock()

					lastError = err // Store the last error encountered
				}()
			} else {
				func() {
					mu.Lock()
					defer mu.Unlock()

					successfulKeys = append(successfulKeys, key) // Append the key if upload was successful
				}()
			}
		}(keys[i], data[i])
	}

	wg.Wait() // Wait for all goroutines to complete

	if lastError != nil {
		return successfulKeys, fmt.Errorf("failed to upload some files: %w", lastError)
	}

	return successfulKeys, nil
}

// Delete - deletes the file from the bucket of the configured spec
func (r *RcloneStorage) Delete(key string) error {
	remotePath := fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key)

	cmd := exec.Command("rclone", "deletefile", remotePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmdOutput := &bytes.Buffer{}
	cmd.Stderr = cmdOutput

	// Execute the command
	err := cmd.Run()
	if err != nil {
		log.WithField("key", key).WithError(err).Errorf("rclone command failed to delete file:%s", cmdOutput.String())
		return fmt.Errorf("rclone command failed: %s, error: %w", cmdOutput.String(), err)
	}

	return nil
}

// CheckCloudConnection verifies the R-clone connection by storing and fetching a test file.
func (r *RcloneStorage) CheckCloudConnection() error {
	testKey := uuid.NewString()
	testData := []byte("this is test data to verify cloud storage connectivity through r-clone")

	_, err := r.Store(testKey, testData)
	if err != nil {
		return err
	}

	_, err = r.Fetch(testKey)
	if err != nil {
		return err
	}

	err = r.Delete(testKey)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBatch deletes a batch of files from the cloud identified by their keys.
func (r *RcloneStorage) DeleteBatch(keys []string) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentUploads) // Use same max concurrency as uploads

	var mu sync.Mutex
	var lastError error

	for _, key := range keys {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore to limit concurrency

		go func(key string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Ensure to release semaphore

			// Construct the remote path where the file is stored
			remotePath := fmt.Sprintf("%s:%s/%s", r.specName, r.bucketName, key)

			// Use rclone to delete the file from the remote
			cmd := exec.Command("rclone", "deletefile", remotePath)
			cmdOutput := &bytes.Buffer{}
			cmd.Stderr = cmdOutput // Capture standard error

			if err := cmd.Run(); err != nil {
				log.WithField("key", key).WithError(err).Errorf("rclone command failed to delete file:%s", cmdOutput.String())
				mu.Lock()
				if lastError == nil {
					lastError = fmt.Errorf("failed to delete some files")
				}
				lastError = fmt.Errorf("%v; key %s error: %v", lastError, key, err)
				mu.Unlock()
			}
		}(key)
	}

	wg.Wait() // Wait for all goroutines to finish

	return lastError // Return the last error encountered, if any
}
