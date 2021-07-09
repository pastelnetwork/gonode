package artworkdownload

// List of task statuses.
const (
	StatusTaskStarted Status = iota
	StatusConnected
	StatusDownloaded
	// Error
	StatusErrorFilesNotMatch
	StatusErrorNotEnoughSuperNode
	StatusErrorNotEnoughFiles
	StatusErrorDownloadFailed
	// Final
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:             "Task Started",
	StatusConnected:               "Connected",
	StatusDownloaded:              "Downloaded",
	StatusErrorFilesNotMatch:      "Error Fingerprints Dont Match",
	StatusErrorNotEnoughSuperNode: "Error Not Enough SuperNode",
	StatusErrorNotEnoughFiles:     "Error Not Enough Downloaded Filed",
	StatusErrorDownloadFailed:     "Error Download Failed",
	StatusTaskRejected:            "Task Rejected",
	StatusTaskCompleted:           "Task Completed",
}

// Status represents status of the task
type Status byte

func (status Status) String() string {
	if name, ok := statusNames[status]; ok {
		return name
	}
	return ""
}

// IsFinal returns true if the status is the final.
func (status Status) IsFinal() bool {
	return status == StatusTaskCompleted || status == StatusTaskRejected
}

// IsFailure returns true if the task failed due to an error
func (status Status) IsFailure() bool {
	return (status == StatusErrorFilesNotMatch ||
		status == StatusErrorNotEnoughSuperNode ||
		status == StatusErrorNotEnoughFiles ||
		status == StatusErrorDownloadFailed ||
		status == StatusTaskRejected)
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
