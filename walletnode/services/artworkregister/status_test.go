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
				StatusConnected,
				StatusImageProbed,
				StatusImageAndThumbnailUploaded,
				StatusTicketAccepted,
				StatusTicketRegistered,
				StatusTicketActivated,
				ErrorInsufficientFee,
				StatusErrorFingerprintsNotMatch,
				StatusErrorThumbnailHashsesNotMatch,
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
		}, {
			status:        StatusImageProbed,
			expectedValue: statusNames[StatusImageProbed],
		}, {
			status:        StatusImageAndThumbnailUploaded,
			expectedValue: statusNames[StatusImageAndThumbnailUploaded],
		}, {
			status:        StatusTicketAccepted,
			expectedValue: statusNames[StatusTicketAccepted],
		}, {
			status:        StatusTicketRegistered,
			expectedValue: statusNames[StatusTicketRegistered],
		}, {
			status:        StatusTicketActivated,
			expectedValue: statusNames[StatusTicketActivated],
		}, {
			status:        ErrorInsufficientFee,
			expectedValue: statusNames[ErrorInsufficientFee],
		}, {
			status:        StatusErrorFingerprintsNotMatch,
			expectedValue: statusNames[StatusErrorFingerprintsNotMatch],
		}, {
			status:        StatusErrorThumbnailHashsesNotMatch,
			expectedValue: statusNames[StatusErrorThumbnailHashsesNotMatch],
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
