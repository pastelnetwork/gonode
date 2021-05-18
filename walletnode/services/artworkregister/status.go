package artworkregister

// List of task statuses.
const (
	StatusTaskStarted Status = iota
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

var statusNames = map[Status]string{
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

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
