package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusTypeNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key                int
		expectedStatusName string
	}{
		{
			key:                0,
			expectedStatusName: statusNames[StatusTaskStarted],
		}, {
			key:                1,
			expectedStatusName: statusNames[StatusConnected],
		},
	}

	statusTypes := StatusTypeNames()

	for _, testCase := range testCases {
		testCase := testCase
		statusType := statusTypes[testCase.key]

		assert.Equal(t, testCase.expectedStatusName, statusType)
	}
}
