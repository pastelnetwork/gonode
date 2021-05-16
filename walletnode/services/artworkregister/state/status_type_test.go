package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusTypeNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectedStatusTypeName string
	}{
		{
			expectedStatusTypeName: statusNames[StatusTaskStarted],
		}, {
			expectedStatusTypeName: statusNames[StatusConnected],
		}, {
			expectedStatusTypeName: statusNames[StatusTicketAccepted],
		}, {
			expectedStatusTypeName: statusNames[StatusTicketRegistered],
		}, {
			expectedStatusTypeName: statusNames[StatusTicketActivated],
		}, {
			expectedStatusTypeName: statusNames[StatusErrorTooLowFee],
		}, {
			expectedStatusTypeName: statusNames[StatusErrorFGPTNotMatch],
		}, {
			expectedStatusTypeName: statusNames[StatusTaskRejected],
		}, {
			expectedStatusTypeName: statusNames[StatusTaskCompleted],
		},
	}

	statusTypes := StatusTypeNames()

	for key, testCase := range testCases {
		testCase := testCase
		key := key

		t.Run(fmt.Sprintf("key:%d/type:%s", key, testCase.expectedStatusTypeName), func(t *testing.T) {
			t.Parallel()

			statusType := statusTypes[key]
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
