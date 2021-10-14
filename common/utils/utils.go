package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
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
