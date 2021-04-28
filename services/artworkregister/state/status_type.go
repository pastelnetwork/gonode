package state

// List of task statuses.
const (
	StatusTaskStarted StatusType = iota

	// Primary node statuses
	StatusRegisterPrimaryConnection
	StatusWaitForSecondaryNodes
	StatusAcceptedSecondaryNodes

	// Secondary node statuses
	StatusRegisterSecondaryConnection
	StatusConnectedToPrimaryNode

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[StatusType]string{
	StatusTaskStarted:                 "Task started",
	StatusRegisterPrimaryConnection:   "Register Primary Connection",
	StatusWaitForSecondaryNodes:       "Wait for Secondary Nodes",
	StatusAcceptedSecondaryNodes:      "Accepted Secondary Nodes",
	StatusRegisterSecondaryConnection: "Register Secondary Connection",
	StatusConnectedToPrimaryNode:      "Connected to Primary Node",
	StatusTaskCanceled:                "Task Canceled",
	StatusTaskCompleted:               "Task Completed",
}

// StatusType represents statusType type of the state.
type StatusType byte

func (statusType StatusType) String() string {
	if name, ok := statusNames[statusType]; ok {
		return name
	}
	return ""
}
