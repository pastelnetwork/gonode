package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusTypeNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectedStatues []StatusType
	}{
		{
			expectedStatues: []StatusType{
				StatusTaskStarted,
				StatusConnected,
				StatusUploadedImage,
				StatusTicketAccepted,
				StatusTicketRegistered,
				StatusTicketActivated,
				StatusErrorTooLowFee,
				StatusErrorFGPTNotMatch,
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

			assert.Equal(t, expectedNames, StatusTypeNames())
		})

	}
}

func TestStatusTypeString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		status        StatusType
		expectedValue string
	}{
		{
			status:        StatusTaskStarted,
			expectedValue: statusNames[StatusTaskStarted],
		}, {
			status:        StatusConnected,
			expectedValue: statusNames[StatusConnected],
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
			status:        StatusErrorTooLowFee,
			expectedValue: statusNames[StatusErrorTooLowFee],
		}, {
			status:        StatusErrorFGPTNotMatch,
			expectedValue: statusNames[StatusErrorFGPTNotMatch],
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
