package errors

import (
	"context"
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
			err:        context.Canceled,
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
