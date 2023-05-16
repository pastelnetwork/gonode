package common_test

import (
	"testing"

	"github.com/pastelnetwork/gonode/hermes/common"
	"github.com/stretchr/testify/assert"
)

func TestIsP2PConnectionErrorTest(t *testing.T) {
	tests := []struct {
		testcase  string
		errString string
		expected  bool
	}{
		{
			testcase:  "when there's connection is closing error, should return true",
			errString: "ERROR app: error deleting rqFileID error=grpc delete err: rpc error: code = Canceled desc = grpc: the client connection is closing",
			expected:  true,
		},
		{
			testcase:  "when there's another, should return false",
			errString: "invalid string",
			expected:  false,
		},
	}

	for _, test := range tests {
		test := test // add this if there's subtest (t.Run)
		t.Run(test.testcase, func(t *testing.T) {

			if exp := common.IsP2PConnectionCloseError(test.errString); exp != test.expected {
				assert.Fail(t, "not the expected output")
			}
		})

	}
}
