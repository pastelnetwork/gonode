package artworkdownload

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusNames(t *testing.T) {
	// t.Parallel()

	testCases := []struct {
		expectedStatues []Status
	}{
		{
			expectedStatues: []Status{
				StatusTaskStarted,
				StatusConnected,
				StatusDownloaded,
				// Error
				StatusErrorFilesNotMatch,
				StatusErrorNotEnoughSuperNode,
				StatusErrorNotEnoughFiles,
				StatusErrorDownloadFailed,
				// Final
				StatusTaskRejected,
				StatusTaskCompleted,
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			// t.Parallel()

			var expectedNames []string
			for _, status := range testCase.expectedStatues {
				expectedNames = append(expectedNames, statusNames[status])
			}

			assert.Equal(t, expectedNames, StatusNames())
		})

	}
}

func TestStatusString(t *testing.T) {
	// t.Parallel()

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
		}, {
			status:        StatusDownloaded,
			expectedValue: statusNames[StatusDownloaded],
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: statusNames[StatusErrorFilesNotMatch],
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: statusNames[StatusErrorNotEnoughSuperNode],
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: statusNames[StatusErrorNotEnoughFiles],
		}, {
			status:        StatusErrorDownloadFailed,
			expectedValue: statusNames[StatusErrorDownloadFailed],
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
			// t.Parallel()

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
		}, {
			status:        StatusDownloaded,
			expectedValue: false,
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: false,
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: false,
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: false,
		}, {
			status:        StatusErrorDownloadFailed,
			expectedValue: false,
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
		}, {
			status:        StatusDownloaded,
			expectedValue: false,
		}, {
			status:        StatusErrorFilesNotMatch,
			expectedValue: true,
		}, {
			status:        StatusErrorNotEnoughSuperNode,
			expectedValue: true,
		}, {
			status:        StatusErrorNotEnoughFiles,
			expectedValue: true,
		}, {
			status:        StatusErrorDownloadFailed,
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
