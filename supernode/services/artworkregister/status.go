package artworkregister

// List of task statuses.
const (
	StatusTaskStarted Status = iota

	// Mode
	StatusPrimaryMode
	StatusSecondaryMode

	// Process
	StatusConnected
	StatusImageUploaded

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:   "Task started",
	StatusPrimaryMode:   "Primary Mode",
	StatusSecondaryMode: "Secondary Mode",
	StatusConnected:     "Connected",
	StatusImageUploaded: "Image Uploaded",
	StatusTaskCanceled:  "Task Canceled",
	StatusTaskCompleted: "Task Completed",
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
	return status == StatusTaskCanceled || status == StatusTaskCompleted
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
