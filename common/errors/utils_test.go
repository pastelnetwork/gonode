package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorStack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err     error
		contain string
	}{
		{
			err:     New("foo"),
			contain: "New(\"foo\")",
		}, {
			err:     fmt.Errorf("bar"),
			contain: "bar",
		}, {
			err:     fmt.Errorf("wrap error: %w", errors.New("inner error")),
			contain: "wrap error: inner error\ninner error",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			errStack := ErrorStack(testCase.err)
			assert.Contains(t, errStack, testCase.contain)
		})
	}
}

func TestIsContextCanceled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err        error
		isCanceled bool
	}{
		{
			err:        New(context.Canceled),
			isCanceled: true,
		},
		{
			err:        New("foo"),
			isCanceled: false,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.isCanceled, IsContextCanceled(testCase.err))
		})
	}
}

func TestRecover(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		f        func()
		expected interface{}
	}{
		{
			f: func() {
				panic(New("foo"))
			},
			expected: "foo",
		}, {
			f: func() {
				panic(struct {
					ID      int
					Message string
				}{
					ID:      1,
					Message: "baz",
				})
			},
			expected: "{1 baz}",
		}, {
			f: func() {
				fmt.Println("should not trigger Recover func")
			},
			expected: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			defer Recover(func(c error) {
				assert.Equal(t, testCase.expected, c.Error())
			})

			testCase.f()
		})
	}

}
