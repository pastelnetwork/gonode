package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusTypeNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key                    int
		expectedStatusTypeName string
	}{
		{
			key:                    0,
			expectedStatusTypeName: statusNames[StatusTaskStarted],
		}, {
			key:                    1,
			expectedStatusTypeName: statusNames[StatusConnected],
		}, {
			key:                    2,
			expectedStatusTypeName: statusNames[StatusTicketAccepted],
		}, {
			key:                    3,
			expectedStatusTypeName: statusNames[StatusTicketRegistered],
		}, {
			key:                    4,
			expectedStatusTypeName: statusNames[StatusTicketActivated],
		}, {
			key:                    5,
			expectedStatusTypeName: statusNames[StatusErrorTooLowFee],
		}, {
			key:                    6,
			expectedStatusTypeName: statusNames[StatusErrorFGPTNotMatch],
		}, {
			key:                    7,
			expectedStatusTypeName: statusNames[StatusTaskRejected],
		}, {
			key:                    8,
			expectedStatusTypeName: statusNames[StatusTaskCompleted],
		},
	}

	statusTypes := StatusTypeNames()

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("key:%d/type:%s", testCase.key, testCase.expectedStatusTypeName), func(t *testing.T) {
			t.Parallel()

			statusType := statusTypes[testCase.key]

			assert.Equal(t, testCase.expectedStatusTypeName, statusType)
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
