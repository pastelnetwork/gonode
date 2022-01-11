package senseregister

// List of task statuses.
const (
	StatusTaskStarted Status = iota
	StatusConnected
	StatusImageProbed

	// Ticket
	StatusTicketAccepted
	StatusTicketRegistered
	StatusTicketActivated
	// Error
	ErrorInsufficientFee
	StatusErrorFindTopNodes
	StatusErrorFindResponsdingSNs
	StatusErrorSignaturesNotMatch
	StatsuErrorInvalidBurnTxID
	StatusErrorFingerprintsNotMatch
	// Final
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:               "Task Started",
	StatusConnected:                 "Connected",
	StatusImageProbed:               "Image Probed",
	StatusTicketAccepted:            "Ticket Accepted",
	StatusTicketRegistered:          "Ticket Registered",
	StatusTicketActivated:           "Ticket Activated",
	ErrorInsufficientFee:            "Error Insufficient Fee",
	StatusErrorFindTopNodes:         "Error Find Top Nodes",
	StatusErrorFindResponsdingSNs:   "Error Find Responsding SNs",
	StatusErrorSignaturesNotMatch:   "Error Signatures Dont Match",
	StatsuErrorInvalidBurnTxID:      "Error Invalid Burn TxID",
	StatusErrorFingerprintsNotMatch: "Error Fingerprints Dont Match",
	StatusTaskRejected:              "Task Rejected",
	StatusTaskCompleted:             "Task Completed",
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

// IsFailure returns true if the status is the failure.
func (status Status) IsFailure() bool {
	return status == ErrorInsufficientFee || status == StatusTaskRejected || status == StatusErrorFingerprintsNotMatch || status == StatusErrorSignaturesNotMatch
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
