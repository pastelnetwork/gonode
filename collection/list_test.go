package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListContains(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		list     []string
		element  string
		expected bool
	}{
		{[]string{}, "", false},
		{[]string{}, "foo", false},
		{[]string{"foo"}, "foo", true},
		{[]string{"bar", "foo", "baz"}, "foo", true},
		{[]string{"bar", "foo", "baz"}, "nope", false},
		{[]string{"bar", "foo", "baz"}, "", false},
	}

	for _, testCase := range testCases {
		actual := ListContains(testCase.list, testCase.element)
		assert.Equal(t, testCase.expected, actual, "For list %v and element %s", testCase.list, testCase.element)
	}
}
