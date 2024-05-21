package storagechallenge

import (
	"os"
)

func createTempDirInPath(basePath string) (string, error) {
	dir, err := os.MkdirTemp(basePath, ".pastel")
	if err != nil {
		return "", err
	}
	return dir, nil
}
