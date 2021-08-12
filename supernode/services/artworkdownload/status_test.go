package artworkdownload

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

				// Process
				StatusFileDecoded,

				// Error
				StatusRequestTooLate,
				StatusArtRegGettingFailed,
				StatusArtRegDecodingFailed,
				StatusArtRegTicketInvalid,
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
			status:        StatusFileDecoded,
			expectedValue: statusNames[StatusFileDecoded],
		}, {
			status:        StatusRequestTooLate,
			expectedValue: statusNames[StatusRequestTooLate],
		}, {
			status:        StatusArtRegGettingFailed,
			expectedValue: statusNames[StatusArtRegGettingFailed],
		}, {
			status:        StatusArtRegDecodingFailed,
			expectedValue: statusNames[StatusArtRegDecodingFailed],
		}, {
			status:        StatusArtRegTicketInvalid,
			expectedValue: statusNames[StatusArtRegTicketInvalid],
		}, {
			status:        StatusListTradeTicketsFailed,
			expectedValue: statusNames[StatusListTradeTicketsFailed],
		}, {
			status:        StatusTradeTicketsNotFound,
			expectedValue: statusNames[StatusTradeTicketsNotFound],
		}, {
			status:        StatusTradeTicketMismatched,
			expectedValue: statusNames[StatusTradeTicketMismatched],
		}, {
			status:        StatusTimestampVerificationFailed,
			expectedValue: statusNames[StatusTimestampVerificationFailed],
		}, {
			status:        StatusTimestampInvalid,
			expectedValue: statusNames[StatusTimestampInvalid],
		}, {
			status:        StatusRQServiceConnectionFailed,
			expectedValue: statusNames[StatusRQServiceConnectionFailed],
		}, {
			status:        StatusSymbolFileNotFound,
			expectedValue: statusNames[StatusSymbolFileNotFound],
		}, {
			status:        StatusSymbolFileInvalid,
			expectedValue: statusNames[StatusSymbolFileInvalid],
		}, {
			status:        StatusSymbolNotFound,
			expectedValue: statusNames[StatusSymbolNotFound],
		}, {
			status:        StatusSymbolMismatched,
			expectedValue: statusNames[StatusSymbolMismatched],
		}, {
			status:        StatusSymbolsNotEnough,
			expectedValue: statusNames[StatusSymbolsNotEnough],
		}, {
			status:        StatusFileDecodingFailed,
			expectedValue: statusNames[StatusFileDecodingFailed],
		}, {
			status:        StatusFileReadingFailed,
			expectedValue: statusNames[StatusFileReadingFailed],
		}, {
			status:        StatusFileMismatched,
			expectedValue: statusNames[StatusFileMismatched],
		}, {
			status:        StatusFileEmpty,
			expectedValue: statusNames[StatusFileEmpty],
		}, {
			status:        StatusTaskCanceled,
			expectedValue: statusNames[StatusTaskCanceled],
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
			status:        StatusArtRegGettingFailed,
			expectedValue: false,
		}, {
			status:        StatusArtRegDecodingFailed,
			expectedValue: false,
		}, {
			status:        StatusArtRegTicketInvalid,
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
			status:        StatusArtRegGettingFailed,
			expectedValue: true,
		}, {
			status:        StatusArtRegDecodingFailed,
			expectedValue: true,
		}, {
			status:        StatusArtRegTicketInvalid,
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
