package memory

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/stretchr/testify/assert"
)

func NewTestDB() *keyValue {
	return &keyValue{
		values: map[string][]byte{
			"exist": []byte("bar"),
		},
	}
}

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
			expectedValue: []byte("bar"),
		}, {
			key:           "not-exist",
			expectedError: storage.ErrKeyNotFound,
			expectedValue: nil,
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, testCase := range testCases {

			testCase := testCase
			testName := fmt.Sprintf("key:%s/value:%v/err:%v", testCase.key, testCase.expectedValue, testCase.expectedError)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				db := NewTestDB()
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
			key:                  "exist",
			value:                []byte("baz"),
			expectedError:        nil,
			expectedLengthValues: 1,
		}, {
			key:                  "foo",
			value:                []byte("grid"),
			expectedError:        nil,
			expectedLengthValues: 2,
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, testCase := range testCases {

			testCase := testCase
			testName := fmt.Sprintf("key:%s/value:%v/length:%d/err:%v", testCase.key, testCase.value, testCase.expectedLengthValues, testCase.expectedError)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				db := NewTestDB()
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
			key:                  "exist",
			expectedError:        nil,
			expectedLengthValues: 0,
		}, {
			key:                  "not-exist",
			expectedError:        nil,
			expectedLengthValues: 1,
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, testCase := range testCases {

			testCase := testCase
			testName := fmt.Sprintf("key:%s/length:%d/err:%v", testCase.key, testCase.expectedLengthValues, testCase.expectedError)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				db := NewTestDB()
				err := db.Delete(testCase.key)
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedLengthValues, len(db.values))
			})
		}
	})
}
