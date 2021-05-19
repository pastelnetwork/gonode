package errors_test

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err      *errors.Error
		expected string
	}{
		{
			err:      errors.New("foo"),
			expected: "foo",
		}, {
			err:      errors.New("bar"),
			expected: "bar",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.err.Error())
		})
	}
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err      *errors.Error
		expected string
	}{
		{
			err:      errors.Errorf("foo: %s", "bar"),
			expected: "foo: bar",
		}, {
			err:      errors.Errorf("foo: %t", true),
			expected: "foo: true",
		}, {
			err:      errors.Errorf("foo: %d", 1),
			expected: "foo: 1",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.err.Error())
		})
	}
}

func TestStack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err     *errors.Error
		contain string
	}{
		{
			err:     errors.New("stack trace"),
			contain: "errors.New(\"stack trace\")",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Contains(t, testCase.err.Stack(), testCase.contain)
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err error
	}{
		{
			err: fmt.Errorf("native error"),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			newError := errors.New(testCase.err)
			assert.Equal(t, testCase.err, newError.Unwrap())
		})
	}

}

func TestErrorStack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err      *errors.Error
		expected string
	}{
		{
			err:      errors.Errorf("foo: %w", errors.New("bar")),
			expected: "foo: bar",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.expected+"\n"+testCase.err.Stack(), testCase.err.ErrorStack())
		})
	}
}

func TestWithField(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		field string
		value interface{}
	}{
		{
			field: "foo",
			value: "bar",
		}, {
			field: "bar",
			value: "baz",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			newError := errors.New("error").WithField(testCase.field, testCase.value)
			fields := errors.ExtractFields(newError)

			value, ok := fields[testCase.field]
			assert.Equal(t, testCase.value, value)
			assert.True(t, ok)

		})
	}
}
