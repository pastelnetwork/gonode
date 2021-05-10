package memory

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		operation     string
		err           error
		keyValue      string
		inputValue    []byte
		expectedValue []byte
	}{
		{
			operation:  "Set",
			err:        nil,
			keyValue:   "test",
			inputValue: []byte("value"),
		},
		{
			operation:     "Get",
			err:           nil,
			keyValue:      "test",
			expectedValue: []byte("value"),
		}, {
			operation: "Delete",
			err:       nil,
			keyValue:  "test",
		}, {
			operation:     "Get",
			err:           storage.ErrKeyNotFound,
			keyValue:      "test",
			expectedValue: nil,
		},
	}

	memStore := NewKeyValue()

	for _, testCase := range testCases {
		testCase := testCase
		var (
			err  error
			data []byte
		)
		//run test in sequence order
		t.Run(fmt.Sprintf("operation-%s-expected-val%v", testCase.operation, testCase.expectedValue), func(t *testing.T) {
			switch testCase.operation {
			case "Set":
				err = memStore.Set(testCase.keyValue, testCase.inputValue)
			case "Get":
				data, err = memStore.Get(testCase.keyValue)
			case "Delete":
				err = memStore.Delete(testCase.keyValue)
			}
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.expectedValue, data)
		})

	}
}
