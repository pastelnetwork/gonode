package utils

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
