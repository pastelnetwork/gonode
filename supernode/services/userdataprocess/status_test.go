package artworkregister

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
				StatusConnected,
				StatusImageProbed,
				StatusImageAndThumbnailCoordinateUploaded,
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
			status:        StatusPrimaryMode,
			expectedValue: statusNames[StatusPrimaryMode],
		}, {
			status:        StatusSecondaryMode,
			expectedValue: statusNames[StatusSecondaryMode],
		}, {
			status:        StatusConnected,
			expectedValue: statusNames[StatusConnected],
		}, {
			status:        StatusImageProbed,
			expectedValue: statusNames[StatusImageProbed],
		}, {
			status:        StatusImageAndThumbnailCoordinateUploaded,
			expectedValue: statusNames[StatusImageAndThumbnailCoordinateUploaded],
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
