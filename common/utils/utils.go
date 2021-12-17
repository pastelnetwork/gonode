package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/crypto/sha3"
)

// DiskStatus cotains info of disk storage
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// SafeErrStr returns err string
func SafeErrStr(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

// Sha3256hash returns SHA-256 hash of input message
func Sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// SafeString returns value of str ptr or empty string if ptr is nil
func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// IsContextErr checks if err is related to context
func IsContextErr(err error) bool {
	if err != nil {
		errStr := strings.ToLower(err.Error())
		return strings.Contains(errStr, "context") || strings.Contains(errStr, "ctx")
	}

	return false
}

// GetExternalIPAddress returns external IP address
func GetExternalIPAddress() (externalIP string, err error) {

	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func B64Encode(in []byte) (out []byte) {
	out = make([]byte, base64.StdEncoding.EncodedLen(len(in)))
	base64.StdEncoding.Encode(out, in)

	return out
}

func B64Decode(in []byte) (out []byte, err error) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(in)))
	n, err := base64.StdEncoding.Decode(out, in)
	if err != nil {
		return nil, errors.Errorf("b64 decode: %w", err)
	}

	return out[:n], nil
}

// EqualStrList tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func EqualStrList(a, b []string) error {
	if len(a) != len(b) {
		return errors.Errorf("unequal length: %d  & %d", len(a), len(b))
	}

	for i, v := range a {
		if v != b[i] {
			return errors.Errorf("index %d mismatch, vals: %s   &  %s", i, v, b[i])
		}
	}

	return nil
}
