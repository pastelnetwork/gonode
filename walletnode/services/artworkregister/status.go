package artworkregister

// List of task statuses.
const (
	StatusTaskStarted Status = iota
	StatusConnected
	StatusImageProbed
	StatusImageAndThumbnailUploaded
	StatusGenRaptorQSymbols

	// Ticket
	StatusPreburntRegistrationFee
	StatusTicketAccepted
	StatusTicketRegistered
	StatusTicketActivated
	// Error
	ErrorInsufficientFee
	StatusErrorFingerprintsNotMatch
	StatusErrorThumbnailHashsesNotMatch
	StatusErrorGenRaptorQSymbolsFailed
	// Final
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:                   "Task Started",
	StatusConnected:                     "Connected",
	StatusImageProbed:                   "Image Probed",
	StatusImageAndThumbnailUploaded:     "Image And Thumbnail Uploaded",
	StatusGenRaptorQSymbols:             "Status Gen ReptorQ Symbols",
	StatusPreburntRegistrationFee:       "Preburn Registration Fee",
	StatusTicketAccepted:                "Ticket Accepted",
	StatusTicketRegistered:              "Ticket Registered",
	StatusTicketActivated:               "Ticket Activated",
	ErrorInsufficientFee:                "Error Insufficient Fee",
	StatusErrorFingerprintsNotMatch:     "Error Fingerprints Dont Match",
	StatusErrorThumbnailHashsesNotMatch: "Error ThumbnailHashes Dont Match",
	StatusErrorGenRaptorQSymbolsFailed:  "Error GenRaptorQ Symbols Failed",
	StatusTaskRejected:                  "Task Rejected",
	StatusTaskCompleted:                 "Task Completed",
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
	return status == ErrorInsufficientFee || status == StatusTaskRejected || status == StatusErrorFingerprintsNotMatch
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
