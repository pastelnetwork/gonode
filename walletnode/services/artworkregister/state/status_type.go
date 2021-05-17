package state

// List of task statuses.
const (
	StatusTaskStarted StatusType = iota
	StatusConnected
	StatusUploadedImage
	// Ticket
	StatusTicketAccepted
	StatusTicketRegistered
	StatusTicketActivated
	// Error
	StatusErrorTooLowFee
	StatusErrorFGPTNotMatch
	// Final
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[StatusType]string{
	StatusTaskStarted:       "Task Started",
	StatusConnected:         "Connected",
	StatusUploadedImage:     "Uploaded Image",
	StatusTicketAccepted:    "Ticket Accepted",
	StatusTicketRegistered:  "Ticket Registered",
	StatusTicketActivated:   "Ticket Activated",
	StatusErrorTooLowFee:    "Error Too Low Fee",
	StatusErrorFGPTNotMatch: "Error FGPT Not Match",
	StatusTaskRejected:      "Task Rejected",
	StatusTaskCompleted:     "Task Completed",
}

// StatusType represents status of the task
type StatusType byte

func (status StatusType) String() string {
	if name, ok := statusNames[status]; ok {
		return name
	}
	return ""
}

// StatusTypeNames returns a sorted list of status names.
func StatusTypeNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
