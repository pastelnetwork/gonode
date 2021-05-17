package state

// List of task statuses.
const (
	StatusTaskStarted StatusType = iota

	// Primary node statuses
	StatusSessionPrimaryNode
	StatusAcceptedNodes

	// Secondary node statuses
	StatusSessionSecondaryNode
	StatusConnectedToNode

	StatusImageUploaded

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[StatusType]string{
	StatusTaskStarted:          "Task started",
	StatusSessionPrimaryNode:   "Session Primary Node",
	StatusSessionSecondaryNode: "Session Secondary Node",
	StatusAcceptedNodes:        "Accepted Secondary Nodes",
	StatusConnectedToNode:      "Connected to Primary Node",
	StatusImageUploaded:        "Image Uploaded",
	StatusTaskCanceled:         "Task Canceled",
	StatusTaskCompleted:        "Task Completed",
}

// StatusType represents statusType type of the state.
type StatusType byte

func (statusType StatusType) String() string {
	if name, ok := statusNames[statusType]; ok {
		return name
	}
	return ""
}
