package state

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			status:        StatusPrimaryMode,
			expectedValue: statusNames[StatusPrimaryMode],
		}, {
			status:        StatusSecondaryMode,
			expectedValue: statusNames[StatusSecondaryMode],
		}, {
			status:        StatusAcceptedNodes,
			expectedValue: statusNames[StatusAcceptedNodes],
		}, {
			status:        StatusConnectedToNode,
			expectedValue: statusNames[StatusConnectedToNode],
		}, {
			status:        StatusImageUploaded,
			expectedValue: statusNames[StatusImageUploaded],
		}, {
			status:        StatusTaskCanceled,
			expectedValue: statusNames[StatusTaskCanceled],
		}, {
			status:        StatusTaskCompleted,
			expectedValue: statusNames[StatusTaskCompleted],
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedValue, testCase.status.String())
		})
	}
}
