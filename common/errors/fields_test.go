package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		fields Fields
		value  string
	}{
		{
			fields: Fields{"foo": "bar"},
			value:  "foo=bar",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			value := testCase.fields.String()
			assert.Equal(t, testCase.value, value)
		})
	}
}

func TestExtractFields(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err    error
		fields Fields
	}{
		{
			err:    New("foo").WithField("bar", "baz"),
			fields: Fields{"bar": "baz"},
		}, {
			err:    fmt.Errorf("native error"),
			fields: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			fields := ExtractFields(testCase.err)
			assert.Equal(t, testCase.fields, fields)
		})
	}
}
