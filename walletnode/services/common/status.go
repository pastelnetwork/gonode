package common

// List of task statuses.
const (
	StatusTaskStarted Status = iota
	StatusConnected

	// Sense and NFT reg
	StatusImageProbed
	StatusImageAndThumbnailUploaded
	StatusGenRaptorQSymbols
	StatusPreburntRegistrationFee

	// NFT Search

	// NFT Download
	StatusDownloaded

	// Request
	StatusTicketAccepted
	StatusTicketRegistered
	StatusTicketActivated

	// Errors
	StatusErrorInsufficientFee
	StatusErrorSignaturesNotMatch
	StatusErrorFingerprintsNotMatch
	StatusErrorThumbnailHashesNotMatch
	StatusErrorGenRaptorQSymbolsFailed
	StatusErrorFilesNotMatch
	StatusErrorNotEnoughSuperNode
	StatusErrorFindRespondingSNs
	StatusErrorNotEnoughFiles
	StatusErrorDownloadFailed
	StatusErrorInvalidBurnTxID

	// Final
	StatusTaskFailed
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted: "Task Started",
	StatusConnected:   "Connected",

	// Sense and NFT reg
	StatusImageProbed:               "Image Probed",
	StatusImageAndThumbnailUploaded: "Image And Thumbnail Uploaded",
	StatusGenRaptorQSymbols:         "Status Gen ReptorQ Symbols",
	StatusPreburntRegistrationFee:   "Preburn Registration Fee",

	// NFT Search

	// NFT Download
	StatusDownloaded: "Downloaded",

	StatusTicketAccepted:   "Request Accepted",
	StatusTicketRegistered: "Request Registered",
	StatusTicketActivated:  "Request Activated",

	// Errors
	StatusErrorInsufficientFee:         "Error Insufficient Fee",
	StatusErrorSignaturesNotMatch:      "Error Signatures Dont Match",
	StatusErrorFingerprintsNotMatch:    "Error Fingerprints Dont Match",
	StatusErrorThumbnailHashesNotMatch: "Error ThumbnailHashes Dont Match",
	StatusErrorGenRaptorQSymbolsFailed: "Error GenRaptorQ Symbols Failed",

	StatusErrorFilesNotMatch:      "Error File Don't Match",
	StatusErrorNotEnoughSuperNode: "Error Not Enough SuperNode",
	StatusErrorFindRespondingSNs:  "Error Find Responding SNs",
	StatusErrorNotEnoughFiles:     "Error Not Enough Downloaded Filed",
	StatusErrorDownloadFailed:     "Error Download Failed",
	StatusErrorInvalidBurnTxID:    "Error Invalid Burn TxID",

	// Final
	StatusTaskFailed:    "Task Failed",
	StatusTaskRejected:  "Task Rejected",
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
	return status == StatusTaskCompleted ||
		status == StatusTaskRejected ||
		status == StatusTaskFailed
}

// IsFailure returns true if the task failed due to an error
func (status Status) IsFailure() bool {
	return status == StatusErrorInsufficientFee ||
		status == StatusErrorSignaturesNotMatch ||
		status == StatusErrorFingerprintsNotMatch ||
		status == StatusErrorThumbnailHashesNotMatch ||
		status == StatusErrorGenRaptorQSymbolsFailed ||
		status == StatusErrorFilesNotMatch ||
		status == StatusErrorNotEnoughSuperNode ||
		status == StatusErrorNotEnoughFiles ||
		status == StatusErrorDownloadFailed ||
		status == StatusErrorFindRespondingSNs ||
		status == StatusErrorInvalidBurnTxID ||
		status == StatusTaskRejected ||
		status == StatusTaskFailed
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
