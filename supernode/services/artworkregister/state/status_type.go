package state

// List of task statuses.
const (
	StatusTaskStarted StatusType = iota

	// Primary node statuses
	StatusWaitForSecondaryNodes
	StatusAcceptedSecondaryNodes

	// Secondary node statuses
	StatusConnectedToPrimaryNode

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[StatusType]string{
	StatusTaskStarted:            "Task started",
	StatusWaitForSecondaryNodes:  "Wait for Secondary Nodes",
	StatusAcceptedSecondaryNodes: "Accepted Secondary Nodes",
	StatusConnectedToPrimaryNode: "Connected to Primary Node",
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
