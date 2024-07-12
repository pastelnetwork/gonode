package download

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

const (
	defaultFileTTL       = 12 * time.Hour
	defaultCheckInterval = 5 * time.Minute
)

// CleanupService cleans up the static dir periodically.
type CleanupService struct {
	fileTTL time.Duration
	fileDir string
}

// Run runs the cleanup service.
func (service *CleanupService) Run(ctx context.Context) error {
	if r := recover(); r != nil {
		log.Errorf("Recovered from panic in cleanup run: %v", r)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled, returning cleanup service")
			return nil
		case <-time.After(defaultCheckInterval):
			count, err := service.cleanup(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("cleanup static dir error")
			} else if count > 0 {
				log.WithContext(ctx).WithField("count", count).Info("cleanup static dir success")
			}
		}
	}
}

func (service *CleanupService) cleanup(ctx context.Context) (int, error) {
	count := 0
	files, err := ioutil.ReadDir(service.fileDir)
	if err != nil {
		return count, err
	}

	cutoffTime := time.Now().UTC().Add(-service.fileTTL)

	for _, f := range files {
		fullPath := filepath.Join(service.fileDir, f.Name())

		// Check if it's a directory
		if f.IsDir() {
			// Remove the directory and its contents if the modification time is before one hour ago
			if f.ModTime().Before(cutoffTime) {
				err := os.RemoveAll(fullPath)
				if err != nil {
					log.WithContext(ctx).WithField("cutoff time", cutoffTime).WithError(err).Error("cleanup directory error")
					continue
				}
				count++
			}
		} else {
			// Remove the file if the modification time is before one hour ago
			if f.ModTime().Before(cutoffTime) {
				err := os.Remove(fullPath)
				if err != nil {
					log.WithContext(ctx).WithField("cutoff time", cutoffTime).WithError(err).Error("cleanup file error")
					continue
				}
				count++
			}
		}
	}

	return count, nil
}

// NewCleanupService returns a new CleanupService instance.
func NewCleanupService(fileDir string) *CleanupService {
	return &CleanupService{
		fileTTL: defaultFileTTL,
		fileDir: fileDir,
	}
}
