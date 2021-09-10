package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/storage"
	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func newTestDB() *keyValue {
	c := cache.New(-1*time.Minute, 12*time.Hour)
	c.Set("exist", []byte("bar"), cache.NoExpiration)
	return &keyValue{
		store: c,
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
			expectedError: storage.ErrKeyValueNotFound,
			expectedValue: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		testName := fmt.Sprintf("key:%s/value:%v/err:%v", testCase.key, testCase.expectedValue, testCase.expectedError)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			db := newTestDB()

			val, err := db.Get(testCase.key)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedValue, val)
		})
	}

}

func TestSet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key           string
		value         []byte
		expectedError error
	}{
		{
			key:           "exist",
			value:         []byte("baz"),
			expectedError: nil,
		}, {
			key:           "foo",
			value:         []byte("grid"),
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		testName := fmt.Sprintf("key:%s/value:%v/err:%v", testCase.key, testCase.value, testCase.expectedError)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			db := newTestDB()

			err := db.Set(testCase.key, testCase.value)
			assert.Equal(t, testCase.expectedError, err)

			value, ok := db.store.Get(testCase.key)
			assert.True(t, ok, "not found new key")
			assert.Equal(t, testCase.value, value.([]byte))
		})
	}

}

func TestDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key           string
		expectedError error
	}{
		{
			key:           "exist",
			expectedError: nil,
		}, {
			key:           "not-exist",
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		testName := fmt.Sprintf("key:%s/err:%v", testCase.key, testCase.expectedError)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			db := newTestDB()

			err := db.Delete(testCase.key)
			assert.Equal(t, testCase.expectedError, err)

			_, ok := db.store.Get(testCase.key)
			assert.False(t, ok, "found deleted key")
		})
	}
}
