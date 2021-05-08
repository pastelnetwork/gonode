package errors_test

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
)

func TestNew(t *testing.T) {

	err := errors.New("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = errors.New(fmt.Errorf("foo"))

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	if err.Goerr().ErrorStack() != err.Goerr().TypeName()+" "+err.Error()+"\n"+string(err.Goerr().Stack()) {
		t.Errorf("ErrorStack is in the wrong format")
	}
}

func TestErrorf(t *testing.T) {

	err := errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	if err.Goerr().ErrorStack() != err.Goerr().TypeName()+" "+err.Error()+"\n"+string(err.Goerr().Stack()) {
		t.Errorf("ErrorStack is in the wrong format")
	}
}

func TestErrorStack(t *testing.T) {
	err := errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	if err.ErrorStack() != err.Error()+"\n"+string(err.Goerr().Stack()) {
		t.Errorf("ErrorStack is in the wrong format")
	}
}

func TestUnwrap(t *testing.T) {

	err := errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = errors.Errorf("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	if err.Unwrap() != err.Goerr().Err {
		t.Errorf("Unwrap returns the wrong error")
	}
}
