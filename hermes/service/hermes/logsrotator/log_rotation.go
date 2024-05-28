package logsrotator

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	runTaskInterval = 30 * time.Minute
	maxLogSize      = 10 * 1024 * 1024 // 10 MB
	destDir         = "/tmp"
)

func (lrs *logRotationService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			log.WithContext(ctx).Debug("log-rotation service run() has been invoked")

			if err := lrs.rotate(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("logs rotation failed")
			}
		}
	}
}

func (lrs *logRotationService) rotate(ctx context.Context) error {
	for appName, logPath := range lrs.appLogMap {
		expandedLogPath := expandHomeDir(logPath)

		fileInfo, err := os.Stat(expandedLogPath)
		if err != nil {
			log.WithContext(ctx).WithField("app_name", appName).
				WithError(err).Error("error getting log file-info")
			continue
		}

		if fileInfo.Size() > maxLogSize {
			// Compress and move the log file
			err = compressAndMoveLog(appName, expandedLogPath)
			if err != nil {
				log.WithContext(ctx).WithField("path", expandedLogPath).
					WithError(err).Error("error rotating log")
			}
		}
	}

	return nil
}

func expandHomeDir(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}
	return path
}

func compressAndMoveLog(appName, logPath string) error {
	timestamp := time.Now().Format("20060102-150405")
	destPath := filepath.Join(destDir, appName, fmt.Sprintf("%s-%s.log.gz", appName, timestamp))

	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return fmt.Errorf("error creating destination directory: %v", err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	srcFile, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, srcFile)
	if err != nil {
		return fmt.Errorf("error compressing log file: %v", err)
	}

	err = os.Truncate(logPath, 0)
	if err != nil {
		return fmt.Errorf("error truncating log file: %v", err)
	}

	return nil
}

func (lrs *logRotationService) Stats(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}
