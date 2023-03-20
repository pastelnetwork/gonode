package common

// List of task statuses.
const (
	StatusTaskStarted Status = iota

	// Mode
	StatusPrimaryMode
	StatusSecondaryMode

	// Process
	StatusConnected

	StatusImageProbed
	StatusAssetUploaded
	StatusImageAndThumbnailCoordinateUploaded
	StatusRegistrationFeeCalculated
	StatusFileDecoded

	// Error
	StatusErrorInvalidBurnTxID
	StatusRequestTooLate
	StatusNftRegGettingFailed
	StatusNftRegDecodingFailed
	StatusNftRegTicketInvalid
	StatusListTradeTicketsFailed
	StatusTradeTicketsNotFound
	StatusTradeTicketMismatched
	StatusTimestampVerificationFailed
	StatusTimestampInvalid
	StatusRQServiceConnectionFailed
	StatusSymbolFileNotFound
	StatusSymbolFileInvalid
	StatusSymbolNotFound
	StatusSymbolMismatched
	StatusSymbolsNotEnough
	StatusFileDecodingFailed
	StatusFileReadingFailed
	StatusFileMismatched
	StatusFileEmpty
	StatusKeyNotFound
	StatusFileRestoreFailed
	StatusFileExists

	// Final
	StatusTaskCanceled
	StatusTaskCompleted
)

var statusNames = map[Status]string{
	StatusTaskStarted:                         "Task started",
	StatusPrimaryMode:                         "Primary Mode",
	StatusSecondaryMode:                       "Secondary Mode",
	StatusConnected:                           "Connected",
	StatusImageProbed:                         "Image Probed",
	StatusAssetUploaded:                       "Asset Uploaded",
	StatusImageAndThumbnailCoordinateUploaded: "Imaged And Thumbnail Coordinate Uploaded",
	StatusRegistrationFeeCalculated:           "Registration Fee Caculated",
	StatusFileDecoded:                         "File Decoded",
	StatusErrorInvalidBurnTxID:                "Error Invalid Burn TxID",
	StatusRequestTooLate:                      "Request too late",
	StatusNftRegGettingFailed:                 "NFT registered getting failed",
	StatusNftRegDecodingFailed:                "NFT registered decoding failed",
	StatusNftRegTicketInvalid:                 "NFT registered ticket invalid",
	StatusListTradeTicketsFailed:              "Could not get available trade tickets",
	StatusTradeTicketsNotFound:                "Trade tickets not found",
	StatusTradeTicketMismatched:               "Trade ticket mismatched",
	StatusTimestampVerificationFailed:         "Could not verify timestamp",
	StatusTimestampInvalid:                    "Timestamp invalid",
	StatusRQServiceConnectionFailed:           "RQ Service connection failed",
	StatusSymbolFileNotFound:                  "Symbol file not found",
	StatusSymbolFileInvalid:                   "Symbol file invalid",
	StatusSymbolNotFound:                      "Symbol not found",
	StatusSymbolMismatched:                    "Symbol mismatched",
	StatusSymbolsNotEnough:                    "Symbols not enough",
	StatusFileDecodingFailed:                  "File decoding failed",
	StatusFileReadingFailed:                   "File reading failed",
	StatusFileEmpty:                           "File empty",
	StatusFileMismatched:                      "File mismatched",
	StatusKeyNotFound:                         "Key not found",
	StatusFileExists:                          "File hash exists",
	StatusFileRestoreFailed:                   "File restore failed",
	StatusTaskCanceled:                        "Task Canceled",
	StatusTaskCompleted:                       "Task Completed",
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
	return status == StatusTaskCanceled || status == StatusTaskCompleted
}

// IsFailure returns true if the task failed due to an error
func (status Status) IsFailure() bool {
	return status == StatusTaskCanceled || status == StatusRequestTooLate ||
		status == StatusNftRegGettingFailed || status == StatusNftRegDecodingFailed ||
		status == StatusNftRegTicketInvalid || status == StatusListTradeTicketsFailed ||
		status == StatusTradeTicketsNotFound || status == StatusTradeTicketMismatched ||
		status == StatusTimestampVerificationFailed || status == StatusTimestampInvalid ||
		status == StatusRQServiceConnectionFailed || status == StatusSymbolFileNotFound ||
		status == StatusSymbolFileInvalid || status == StatusSymbolNotFound ||
		status == StatusSymbolMismatched || status == StatusSymbolsNotEnough ||
		status == StatusFileDecodingFailed || status == StatusFileReadingFailed ||
		status == StatusFileEmpty || status == StatusFileMismatched ||
		status == StatusKeyNotFound || status == StatusFileRestoreFailed || status == StatusFileExists
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
