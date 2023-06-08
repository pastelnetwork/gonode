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
	defaultFileTTL       = 1 * time.Hour
	defaultCheckInterval = 1 * time.Minute
)

type CleanupService struct {
	fileTTL time.Duration
	fileDir string
}

func (service *CleanupService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.WithContext(ctx).Info("context canceled, returning cleanup service")
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

	cutoffTime := time.Now().Add(-service.fileTTL)

	for _, f := range files {
		// Skip if it's a directory
		if f.IsDir() {
			continue
		}

		// If the file's modification time is before one hour ago
		if f.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(service.fileDir, f.Name())
			err := os.Remove(fullPath)
			if err != nil {
				log.WithContext(ctx).WithField("cutoff time", cutoffTime).WithError(err).Error("cleanup file error")
				continue
			}
			count++
		}
	}

	return count, nil
}

func NewCleanupService(fileDir string) *CleanupService {
	return &CleanupService{
		fileTTL: defaultFileTTL,
		fileDir: fileDir,
	}
}
