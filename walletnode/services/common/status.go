package common

// List of task statuses.
const (
	StatusTaskStarted Status = iota
	StatusConnected

	StatusValidateDuplicateTickets
	StatusValidateBurnTxn
	StatusBurnTxnValidated

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
	StatusErrorMeshSetupFailed // task for status logs
	StatusErrorSendingRegMetadata
	StatusErrorUploadImageFailed
	StatusErrorConvertingImageBytes
	StatusErrorEncodingImage
	StatusErrorCreatingTicket
	StatusErrorSigningTicket
	StatusErrorUploadingTicket
	StatusErrorActivatingTicket
	StatusErrorProbeImage
	StatusErrorCheckDDServerAvailability
	StatusErrorGenerateDDAndFPIds
	StatusErrorCheckingStorageFee
	StatusErrorInsufficientBalance
	StatusErrorGetImageHash
	StatusErrorSendSignTicketFailed
	StatusErrorCheckBalance
	StatusErrorPreburnRegFeeGetTicketID
	StatusErrorValidateRegTxnIDFailed
	StatusErrorValidateActivateTxnIDFailed

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
	StatusErrorOwnershipNotMatch
	StatusErrorHashMismatch

	// Final
	StatusTaskFailed
	StatusTaskRejected
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:              "Task Started",
	StatusConnected:                "Connected",
	StatusValidateBurnTxn:          "Validating Burn Txn",
	StatusBurnTxnValidated:         "Burn Txn Validated",
	StatusValidateDuplicateTickets: "Validated Duplicate Reg Tickets",

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
	StatusErrorMeshSetupFailed:             "Error Setting up mesh of supernodes",
	StatusErrorSendingRegMetadata:          "Error Sending Reg Metadata",
	StatusErrorUploadImageFailed:           "Error Uploading Image",
	StatusErrorConvertingImageBytes:        "Error Converting Image to Bytes",
	StatusErrorEncodingImage:               "Error Encoding Image",
	StatusErrorCreatingTicket:              "Error Creating Ticket",
	StatusErrorSigningTicket:               "Error Signing Ticket",
	StatusErrorUploadingTicket:             "Error Uploading Ticket",
	StatusErrorActivatingTicket:            "Error Activating Ticket",
	StatusErrorProbeImage:                  "Error Probing Image",
	StatusErrorCheckDDServerAvailability:   "Error checking dd-server availability before probe image",
	StatusErrorGenerateDDAndFPIds:          "Error Generating DD and Fingerprint IDs",
	StatusErrorCheckingStorageFee:          "Error comparing suitable storage fee with task request maximum fee",
	StatusErrorInsufficientBalance:         "Error balance not sufficient",
	StatusErrorGetImageHash:                "Error getting hash of the image",
	StatusErrorSendSignTicketFailed:        "Error sending signed ticket to SNs",
	StatusErrorCheckBalance:                "Error checking balance",
	StatusErrorPreburnRegFeeGetTicketID:    "Error burning reg fee to get reg ticket id",
	StatusErrorValidateActivateTxnIDFailed: "Error validating activate ticket txn id",
	StatusErrorValidateRegTxnIDFailed:      "Error validating reg ticket txn id",

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

	StatusErrorOwnershipNotMatch: "Error Ownership Not Match",
	StatusErrorHashMismatch:      "Error Hash Mismatch",

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
		status == StatusTaskFailed ||
		status == StatusErrorOwnershipNotMatch ||
		status == StatusErrorHashMismatch
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}

// EphemeralStatus allows a task to add something like "Created txid: <txid>" to the task history
type EphemeralStatus struct {
	StatusTitle   string
	StatusString  string
	IsFinalBool   bool
	IsFailureBool bool
}

// String returns the string representation of the status by concatenation
func (status EphemeralStatus) String() string {
	return status.StatusTitle + status.StatusString
}

// IsFinal tracks whether this should be a final status (almost always no)
func (status EphemeralStatus) IsFinal() bool {
	return status.IsFinalBool
}

// IsFailure tracks whether this should be a failure status (almost always no)
func (status EphemeralStatus) IsFailure() bool {
	return status.IsFailureBool
}
