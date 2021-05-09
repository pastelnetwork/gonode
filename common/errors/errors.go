package errors

import (
	"fmt"

	goerrors "github.com/go-errors/errors"
)

const defaultCallersSkip = 2

// Error is a custom error
type Error struct {
	goerror *goerrors.Error
	fields  Fields
}

// Error returns the underlying error's message.
func (err *Error) Error() string {
	return err.goerror.Error()
}

// Stack returns a string that contains the callstack.
func (err *Error) Stack() string {
	return string(err.goerror.Stack())
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (err *Error) ErrorStack() string {
	return err.Error() + "\n" + string(err.goerror.Stack())
}

// Unwrap unpacks wrapped native errors
func (err *Error) Unwrap() error {
	return err.goerror.Err
}

// WithField adds a new field
func (err *Error) WithField(key string, value interface{}) *Error {
	err.fields[key] = value
	return err
}

// Errorf creates a new error with the given message. You can use it
// as a drop-in replacement for fmt.Errorf() to provide descriptive
// errors in return values.
func Errorf(format string, vals ...interface{}) *Error {
	return newWithSkip(fmt.Errorf(format, vals...), defaultCallersSkip)
}

// New makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The stacktrace will point to the line of code that
// called New.
func New(err interface{}) *Error {
	return newWithSkip(err, defaultCallersSkip)
}

func newWithSkip(err interface{}, skip int) *Error {
	return &Error{
		goerror: goerrors.Wrap(err, skip),
		fields:  make(Fields),
	}
}
