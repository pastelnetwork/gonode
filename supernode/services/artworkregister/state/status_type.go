package state

// List of task statuses.
const (
	StatusTaskStarted StatusType = iota

	// Primary node statuses
	StatusHandshakePrimaryNode
	StatusAcceptedNodes

	// Secondary node statuses
	StatusHandshakeSecondaryNode
	StatusConnectedToNode

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[StatusType]string{
	StatusTaskStarted:            "Task started",
	StatusHandshakePrimaryNode:   "Handshake Primary Node",
	StatusHandshakeSecondaryNode: "Handshake Secondary Node",
	StatusAcceptedNodes:          "Accepted Secondary Nodes",
	StatusConnectedToNode:        "Connected to Primary Node",
	StatusTaskCanceled:           "Task Canceled",
	StatusTaskCompleted:          "Task Completed",
}

// StatusType represents statusType type of the state.
type StatusType byte

func (statusType StatusType) String() string {
	if name, ok := statusNames[statusType]; ok {
		return name
	}
	return ""
}
