package artworkregister

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/crypto/sha3"
)

func sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func safeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// createRQIDSFile create a file with the following content
// RANDOM-GUID1..X
// BLOCK_HASH
// raptroq id1
// ...
// raptroq idN
// Return absolute path to the created file
func createRQIDSFile(rqids []string, blockHash []byte) (string, error) {
	file, err := os.CreateTemp(os.TempDir(), "rqids")
	if err != nil {
		return "", errors.Errorf("failed to create file at %s %w", os.TempDir(), err)
	}
	defer file.Close()

	var content []string

	uuid := uuid.New()
	content = append(content, uuid.String())

	content = append(content, fmt.Sprintf("%x", blockHash))

	hasher := sha3.New256()

	for _, rqid := range rqids {
		_, err := hasher.Write([]byte(rqid))
		if err != nil {
			return "", errors.Errorf("failed to hash rqid %s %w", rqid, err)
		}

		hash := hasher.Sum(nil)
		hasher.Reset()

		content = append(content, base64.StdEncoding.EncodeToString(hash))
	}

	if _, err := file.Write([]byte(strings.Join(content, "\n"))); err != nil {
		return "", errors.Errorf("failed to write to file %s %w", file.Name(), err)
	}

	absPath, err := filepath.Abs(path.Join(os.TempDir(), file.Name()))
	if err != nil {
		return "", errors.Errorf("failed to get absoulte path of file %s %w", file.Name(), err)
	}
	return absPath, nil
}
