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
				StatusPrimaryMode,
				StatusSecondaryMode,

				// Process
				StatusConnected,
				StatusImageProbed,
				StatusAssetUploaded,
				StatusImageAndThumbnailCoordinateUploaded,
				StatusRegistrationFeeCalculated,
				StatusFileDecoded,

				// Error
				StatusErrorInvalidBurnTxID,
				StatusRequestTooLate,
				StatusNftRegGettingFailed,
				StatusNftRegDecodingFailed,
				StatusNftRegTicketInvalid,
				StatusListTradeTicketsFailed,
				StatusTradeTicketsNotFound,
				StatusTradeTicketMismatched,
				StatusTimestampVerificationFailed,
				StatusTimestampInvalid,
				StatusRQServiceConnectionFailed,
				StatusSymbolFileNotFound,
				StatusSymbolFileInvalid,
				StatusSymbolNotFound,
				StatusSymbolMismatched,
				StatusSymbolsNotEnough,
				StatusFileDecodingFailed,
				StatusFileReadingFailed,
				StatusFileMismatched,
				StatusFileEmpty,
				StatusKeyNotFound,
				StatusFileRestoreFailed,
				StatusFileExists,

				// Final
				StatusTaskCanceled,
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
				expectedNames = append(expectedNames, StatusNames()[status])
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
			expectedValue: StatusNames()[StatusTaskStarted],
		}, {
			status:        StatusFileDecoded,
			expectedValue: StatusNames()[StatusFileDecoded],
		}, {
			status:        StatusRequestTooLate,
			expectedValue: StatusNames()[StatusRequestTooLate],
		}, {
			status:        StatusNftRegGettingFailed,
			expectedValue: StatusNames()[StatusNftRegGettingFailed],
		}, {
			status:        StatusNftRegDecodingFailed,
			expectedValue: StatusNames()[StatusNftRegDecodingFailed],
		}, {
			status:        StatusNftRegTicketInvalid,
			expectedValue: StatusNames()[StatusNftRegTicketInvalid],
		}, {
			status:        StatusListTradeTicketsFailed,
			expectedValue: StatusNames()[StatusListTradeTicketsFailed],
		}, {
			status:        StatusTradeTicketsNotFound,
			expectedValue: StatusNames()[StatusTradeTicketsNotFound],
		}, {
			status:        StatusTradeTicketMismatched,
			expectedValue: StatusNames()[StatusTradeTicketMismatched],
		}, {
			status:        StatusTimestampVerificationFailed,
			expectedValue: StatusNames()[StatusTimestampVerificationFailed],
		}, {
			status:        StatusTimestampInvalid,
			expectedValue: StatusNames()[StatusTimestampInvalid],
		}, {
			status:        StatusRQServiceConnectionFailed,
			expectedValue: StatusNames()[StatusRQServiceConnectionFailed],
		}, {
			status:        StatusSymbolFileNotFound,
			expectedValue: StatusNames()[StatusSymbolFileNotFound],
		}, {
			status:        StatusSymbolFileInvalid,
			expectedValue: StatusNames()[StatusSymbolFileInvalid],
		}, {
			status:        StatusSymbolNotFound,
			expectedValue: StatusNames()[StatusSymbolNotFound],
		}, {
			status:        StatusSymbolMismatched,
			expectedValue: StatusNames()[StatusSymbolMismatched],
		}, {
			status:        StatusSymbolsNotEnough,
			expectedValue: StatusNames()[StatusSymbolsNotEnough],
		}, {
			status:        StatusFileDecodingFailed,
			expectedValue: StatusNames()[StatusFileDecodingFailed],
		}, {
			status:        StatusFileReadingFailed,
			expectedValue: StatusNames()[StatusFileReadingFailed],
		}, {
			status:        StatusFileMismatched,
			expectedValue: StatusNames()[StatusFileMismatched],
		}, {
			status:        StatusFileEmpty,
			expectedValue: StatusNames()[StatusFileEmpty],
		}, {
			status:        StatusTaskCanceled,
			expectedValue: StatusNames()[StatusTaskCanceled],
		}, {
			status:        StatusTaskCompleted,
			expectedValue: StatusNames()[StatusTaskCompleted],
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
	t.Parallel()

	testCases := []struct {
		status        Status
		expectedValue bool
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: false,
		}, {
			status:        StatusFileDecoded,
			expectedValue: false,
		}, {
			status:        StatusRequestTooLate,
			expectedValue: false,
		}, {
			status:        StatusNftRegGettingFailed,
			expectedValue: false,
		}, {
			status:        StatusNftRegDecodingFailed,
			expectedValue: false,
		}, {
			status:        StatusNftRegTicketInvalid,
			expectedValue: false,
		}, {
			status:        StatusListTradeTicketsFailed,
			expectedValue: false,
		}, {
			status:        StatusTradeTicketsNotFound,
			expectedValue: false,
		}, {
			status:        StatusTradeTicketMismatched,
			expectedValue: false,
		}, {
			status:        StatusTimestampVerificationFailed,
			expectedValue: false,
		}, {
			status:        StatusTimestampInvalid,
			expectedValue: false,
		}, {
			status:        StatusRQServiceConnectionFailed,
			expectedValue: false,
		}, {
			status:        StatusSymbolFileNotFound,
			expectedValue: false,
		}, {
			status:        StatusSymbolFileInvalid,
			expectedValue: false,
		}, {
			status:        StatusSymbolNotFound,
			expectedValue: false,
		}, {
			status:        StatusSymbolMismatched,
			expectedValue: false,
		}, {
			status:        StatusSymbolsNotEnough,
			expectedValue: false,
		}, {
			status:        StatusFileDecodingFailed,
			expectedValue: false,
		}, {
			status:        StatusFileReadingFailed,
			expectedValue: false,
		}, {
			status:        StatusFileMismatched,
			expectedValue: false,
		}, {
			status:        StatusFileEmpty,
			expectedValue: false,
		}, {
			status:        StatusTaskCanceled,
			expectedValue: true,
		}, {
			status:        StatusTaskCompleted,
			expectedValue: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("status:%v/value:%v", testCase.status, testCase.expectedValue), func(t *testing.T) {
			t.Parallel()

			value := testCase.status.IsFinal()
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}

func TestStatusIsFailure(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		status        Status
		expectedValue bool
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: false,
		}, {
			status:        StatusFileDecoded,
			expectedValue: false,
		}, {
			status:        StatusRequestTooLate,
			expectedValue: true,
		}, {
			status:        StatusNftRegGettingFailed,
			expectedValue: true,
		}, {
			status:        StatusNftRegDecodingFailed,
			expectedValue: true,
		}, {
			status:        StatusNftRegTicketInvalid,
			expectedValue: true,
		}, {
			status:        StatusListTradeTicketsFailed,
			expectedValue: true,
		}, {
			status:        StatusTradeTicketsNotFound,
			expectedValue: true,
		}, {
			status:        StatusTradeTicketMismatched,
			expectedValue: true,
		}, {
			status:        StatusTimestampVerificationFailed,
			expectedValue: true,
		}, {
			status:        StatusTimestampInvalid,
			expectedValue: true,
		}, {
			status:        StatusRQServiceConnectionFailed,
			expectedValue: true,
		}, {
			status:        StatusSymbolFileNotFound,
			expectedValue: true,
		}, {
			status:        StatusSymbolFileInvalid,
			expectedValue: true,
		}, {
			status:        StatusSymbolNotFound,
			expectedValue: true,
		}, {
			status:        StatusSymbolMismatched,
			expectedValue: true,
		}, {
			status:        StatusSymbolsNotEnough,
			expectedValue: true,
		}, {
			status:        StatusFileDecodingFailed,
			expectedValue: true,
		}, {
			status:        StatusFileReadingFailed,
			expectedValue: true,
		}, {
			status:        StatusFileMismatched,
			expectedValue: true,
		}, {
			status:        StatusFileEmpty,
			expectedValue: true,
		}, {
			status:        StatusTaskCanceled,
			expectedValue: true,
		}, {
			status:        StatusTaskCompleted,
			expectedValue: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("status:%v/value:%v", testCase.status, testCase.expectedValue), func(t *testing.T) {
			t.Parallel()

			value := testCase.status.IsFailure()
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}
