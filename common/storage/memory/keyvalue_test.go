package memory

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key           string
		expectedError error
		expectedValue []byte
	}{
		{
			key:           "exist",
			expectedError: nil,
			expectedValue: []byte("value-exist"),
		}, {
			key:           "not-exist",
			expectedError: storage.ErrKeyNotFound,
			expectedValue: nil,
		},
	}

	t.Run("group", func(t *testing.T) {
		db := &keyValue{
			values: map[string][]byte{
				"exist": []byte("value-exist"),
			},
		}

		for _, testCase := range testCases {
			testCase := testCase
			testName := fmt.Sprintf("key:%s/value:%v/err:%v", testCase.key, testCase.expectedValue, testCase.expectedError)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()

				val, err := db.Get(testCase.key)
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedValue, val)
			})
		}

	})
}

func TestSet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key                  string
		value                []byte
		expectedError        error
		expectedLengthValues int
	}{
		{
			key:                  "foo",
			value:                []byte("bar"),
			expectedError:        nil,
			expectedLengthValues: 1,
		}, {
			key:                  "foo",
			value:                []byte("grid"),
			expectedError:        nil,
			expectedLengthValues: 1,
		},
	}

	t.Run("group", func(t *testing.T) {
		db := &keyValue{
			values: make(map[string][]byte),
		}

		for _, testCase := range testCases {
			testCase := testCase
			testName := fmt.Sprintf("key:%s/value:%v/length:%d/err:%v", testCase.key, testCase.value, testCase.expectedLengthValues, testCase.expectedError)

			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				err := db.Set(testCase.key, testCase.value)
				assert.Equal(t, testCase.expectedError, err)
				value := db.values[testCase.key]
				assert.Equal(t, testCase.value, value)
				assert.Equal(t, testCase.expectedLengthValues, len(db.values))
			})
		}
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key                  string
		expectedError        error
		expectedLengthValues int
	}{
		{
			key:                  "exist-key",
			expectedError:        nil,
			expectedLengthValues: 1,
		}, {
			key:                  "not-exist-key",
			expectedError:        nil,
			expectedLengthValues: 1,
		},
	}

	t.Run("group", func(t *testing.T) {
		db := &keyValue{
			values: map[string][]byte{
				"exist-key":   []byte("bar"),
				"exist-key-1": []byte("foo"),
			},
		}
		for _, testCase := range testCases {
			testCase := testCase
			testName := fmt.Sprintf("key:%s/length:%d/err:%v", testCase.key, testCase.expectedLengthValues, testCase.expectedError)

			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				err := db.Delete(testCase.key)
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedLengthValues, len(db.values))
			})
		}
	})
}
