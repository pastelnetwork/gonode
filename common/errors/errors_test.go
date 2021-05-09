package errors_test

import (
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"gotest.tools/assert"
)

func TestNew(t *testing.T) {
	err := errors.New("foo")
	assert.Equal(t, "foo", err.Error())

	err = errors.New("bar")
	assert.Equal(t, "bar", err.Error())
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("foo")
	assert.Equal(t, "foo", err.Error())
}

func TestErrorStack(t *testing.T) {
	err := errors.Errorf("foo")
	assert.Equal(t, "foo", err.Error())
	assert.Equal(t, err.Error()+"\n"+err.Stack(), err.ErrorStack())
}
