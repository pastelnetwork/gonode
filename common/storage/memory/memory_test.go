package memory

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	testCase := map[string]struct {
		operation              string
		needCheckExpectedValue bool
		expectedError          bool
		keyValue               string
		inputValue             string
		expectedValue          string
	}{
		"1. Set": {
			operation:     "Set",
			expectedError: false,
			keyValue:      "test",
			inputValue:    "value",
		},
		"2. Get": {
			operation:              "Get",
			expectedError:          false,
			needCheckExpectedValue: true,
			keyValue:               "test",
			expectedValue:          "value",
		},
		"3. Delete": {
			operation:     "Delete",
			expectedError: false,
			keyValue:      "test",
		},
		"4. Get Deleted key": {
			operation:     "Get",
			expectedError: true,
			keyValue:      "test",
		},
	}

	memStore := NewKeyValue()
	var (
		err  error
		data []byte
	)
	//sorting map for consistently order
	//https://blog.golang.org/maps#TOC_7.
	var keys []string
	for k := range testCase {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		//reset all variable
		err, data = nil, nil

		tc := testCase[name]
		t.Run(name, func(ts *testing.T) {
			switch tc.operation {
			case "Set":
				err = memStore.Set(tc.keyValue, []byte(tc.inputValue))
			case "Get":
				data, err = memStore.Get(tc.keyValue)
			case "Delete":
				err = memStore.Delete(tc.keyValue)
			}

			if tc.expectedError {
				assert.Error(ts, err)
			} else {
				assert.NoError(ts, err)
			}

			if tc.needCheckExpectedValue {
				assert.Equal(ts, tc.expectedValue, string(data))
			}

		})

	}
}
