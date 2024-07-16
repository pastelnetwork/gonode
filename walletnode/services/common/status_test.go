package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectedStatues []Status
	}{
		{
			expectedStatues: []Status{
				StatusTaskStarted,
				StatusConnected,
				StatusValidateDuplicateTickets,
				StatusValidateBurnTxn,
				StatusBurnTxnValidated,
				// Sense and NFT reg
				StatusImageProbed,
				StatusImageAndThumbnailUploaded,
				StatusGenRaptorQSymbols,
				StatusPreburntRegistrationFee,
				// NFT Download
				StatusDownloaded,
				// Tickets
				StatusTicketAccepted,
				StatusTicketRegistered,
				StatusTicketActivated,
				//Errors
				StatusErrorMeshSetupFailed, // task for status logs
				StatusErrorSendingRegMetadata,
				StatusErrorUploadImageFailed,
				StatusErrorConvertingImageBytes,
				StatusErrorEncodingImage,
				StatusErrorCreatingTicket,
				StatusErrorSigningTicket,
				StatusErrorUploadingTicket,
				StatusErrorActivatingTicket,
				StatusErrorProbeImage,
				StatusErrorCheckDDServerAvailability,
				StatusErrorGenerateDDAndFPIds,
				StatusErrorCheckingStorageFee,
				StatusErrorInsufficientBalance,
				StatusErrorGetImageHash,
				StatusErrorSendSignTicketFailed,
				StatusErrorCheckBalance,
				StatusErrorPreburnRegFeeGetTicketID,
				StatusErrorValidateRegTxnIDFailed,
				StatusErrorValidateActivateTxnIDFailed,
				StatusErrorInsufficientFee,
				StatusErrorSignaturesNotMatch,
				StatusErrorFingerprintsNotMatch,
				StatusErrorThumbnailHashesNotMatch,
				StatusErrorGenRaptorQSymbolsFailed,
				StatusErrorFilesNotMatch,
				StatusErrorNotEnoughSuperNode,
				StatusErrorFindRespondingSNs,
				StatusErrorNotEnoughFiles,
				StatusErrorDownloadFailed,
				StatusErrorInvalidBurnTxID,
				StatusErrorOwnershipNotMatch,
				StatusErrorHashMismatch,
				// Final
				StatusTaskFailed,
				StatusTaskRejected,
				StatusTaskCompleted,
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			t.Parallel()

			var expectedNames []string
			for _, status := range testCase.expectedStatues {
				expectedNames = append(expectedNames, statusNames[status])
			}

			assert.Equal(t, expectedNames, StatusNames())
		})

	}
}

func TestStatusString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		status        Status
		expectedValue string
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: statusNames[StatusTaskStarted],
		}, {
			status:        StatusConnected,
			expectedValue: statusNames[StatusConnected],
		},
		// Sense and NFT reg
		{
			status:        StatusImageProbed,
			expectedValue: statusNames[StatusImageProbed],
		}, {
			status:        StatusImageAndThumbnailUploaded,
			expectedValue: statusNames[StatusImageAndThumbnailUploaded],
		}, {
			status:        StatusGenRaptorQSymbols,
			expectedValue: statusNames[StatusGenRaptorQSymbols],
		}, {
			status:        StatusPreburntRegistrationFee,
			expectedValue: statusNames[StatusPreburntRegistrationFee],
		},
		// NFT Download
		{
			status:        StatusDownloaded,
			expectedValue: statusNames[StatusDownloaded],
		},
		// Tickets
		{
			status:        StatusTicketAccepted,
			expectedValue: statusNames[StatusTicketAccepted],
		}, {
			status:        StatusTicketRegistered,
			expectedValue: statusNames[StatusTicketRegistered],
		}, {
			status:        StatusTicketActivated,
			expectedValue: statusNames[StatusTicketActivated],
		},
		// Errors
		{
			status:        StatusErrorInsufficientFee,
			expectedValue: statusNames[StatusErrorInsufficientFee],
		}, {
			status:        StatusErrorFingerprintsNotMatch,
			expectedValue: statusNames[StatusErrorFingerprintsNotMatch],
		}, {
			status:        StatusErrorThumbnailHashesNotMatch,
			expectedValue: statusNames[StatusErrorThumbnailHashesNotMatch],
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: statusNames[StatusErrorFilesNotMatch],
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: statusNames[StatusErrorNotEnoughSuperNode],
		}, {
			status:        StatusErrorFindRespondingSNs,
			expectedValue: statusNames[StatusErrorFindRespondingSNs],
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: statusNames[StatusErrorNotEnoughFiles],
		}, {
			status:        StatusErrorDownloadFailed,
			expectedValue: statusNames[StatusErrorDownloadFailed],
		}, {
			status:        StatusErrorInvalidBurnTxID,
			expectedValue: statusNames[StatusErrorInvalidBurnTxID],
		},
		// Final
		{
			status:        StatusTaskFailed,
			expectedValue: statusNames[StatusTaskFailed],
		}, {
			status:        StatusTaskRejected,
			expectedValue: statusNames[StatusTaskRejected],
		}, {
			status:        StatusTaskCompleted,
			expectedValue: statusNames[StatusTaskCompleted],
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("status:%v/value:%s", testCase.status, testCase.expectedValue), func(t *testing.T) {
			t.Parallel()

			value := testCase.status.String()
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}

func TestStatusIsFinal(t *testing.T) {
	// t.Parallel()

	testCases := []struct {
		status        Status
		expectedValue bool
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: false,
		}, {
			status:        StatusConnected,
			expectedValue: false,
		},
		// Sense and NFT reg
		{
			status:        StatusImageProbed,
			expectedValue: false,
		}, {
			status:        StatusImageAndThumbnailUploaded,
			expectedValue: false,
		}, {
			status:        StatusGenRaptorQSymbols,
			expectedValue: false,
		}, {
			status:        StatusPreburntRegistrationFee,
			expectedValue: false,
		},
		// NFT Download
		{
			status:        StatusDownloaded,
			expectedValue: false,
		},
		// Errors
		{
			status:        StatusErrorInsufficientFee,
			expectedValue: false,
		}, {
			status:        StatusErrorSignaturesNotMatch,
			expectedValue: false,
		}, {
			status:        StatusErrorFingerprintsNotMatch,
			expectedValue: false,
		}, {
			status:        StatusErrorThumbnailHashesNotMatch,
			expectedValue: false,
		}, {
			status:        StatusErrorGenRaptorQSymbolsFailed,
			expectedValue: false,
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: false,
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: false,
		}, {
			status:        StatusErrorFindRespondingSNs,
			expectedValue: false,
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: false,
		}, {
			status:        StatusErrorDownloadFailed,
			expectedValue: false,
		}, {
			status:        StatusErrorInvalidBurnTxID,
			expectedValue: false,
		},
		// Final
		{
			status:        StatusTaskFailed,
			expectedValue: true,
		}, {
			status:        StatusTaskRejected,
			expectedValue: true,
		}, {
			status:        StatusTaskCompleted,
			expectedValue: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("status:%v/value:%v", testCase.status, testCase.expectedValue), func(t *testing.T) {
			// t.Parallel()

			value := testCase.status.IsFinal()
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}

func TestStatusIsFailure(t *testing.T) {
	// t.Parallel()

	testCases := []struct {
		status        Status
		expectedValue bool
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: false,
		}, {
			status:        StatusConnected,
			expectedValue: false,
		},
		// Sense and NFT reg
		{
			status:        StatusImageProbed,
			expectedValue: false,
		}, {
			status:        StatusImageAndThumbnailUploaded,
			expectedValue: false,
		}, {
			status:        StatusGenRaptorQSymbols,
			expectedValue: false,
		}, {
			status:        StatusPreburntRegistrationFee,
			expectedValue: false,
		},
		// NFT Download
		{
			status:        StatusDownloaded,
			expectedValue: false,
		},
		// Errors
		{
			status:        StatusErrorInsufficientFee,
			expectedValue: true,
		}, {
			status:        StatusErrorSignaturesNotMatch,
			expectedValue: true,
		}, {
			status:        StatusErrorFingerprintsNotMatch,
			expectedValue: true,
		}, {
			status:        StatusErrorThumbnailHashesNotMatch,
			expectedValue: true,
		}, {
			status:        StatusErrorGenRaptorQSymbolsFailed,
			expectedValue: true,
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: true,
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: true,
		}, {
			status:        StatusErrorFindRespondingSNs,
			expectedValue: true,
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: true,
		}, {
			status:        StatusErrorDownloadFailed,
			expectedValue: true,
		}, {
			status:        StatusErrorInvalidBurnTxID,
			expectedValue: true,
		},
		// Final
		{
			status:        StatusTaskFailed,
			expectedValue: true,
		}, {
			status:        StatusTaskRejected,
			expectedValue: true,
		}, {
			status:        StatusTaskCompleted,
			expectedValue: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("status:%v/value:%v", testCase.status, testCase.expectedValue), func(t *testing.T) {
			// t.Parallel()

			value := testCase.status.IsFailure()
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}
