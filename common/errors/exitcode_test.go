package errors

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type exitStatusTest struct {
	code    int
	message string
}

func (e *exitStatusTest) ExitStatus() int {
	return e.code
}

func (e *exitStatusTest) Error() string {
	return e.message
}

func TestExitCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		err      error
		exitCode int
	}{
		{
			err: &exitStatusTest{
				code:    2,
				message: "exit status",
			},
			exitCode: 2,
		}, {
			err: &exec.ExitError{
				ProcessState: &os.ProcessState{},
			},
			exitCode: 1,
		}, {
			err:      New("foo"),
			exitCode: 1,
		}, {
			err:      nil,
			exitCode: 0,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			exitCode := ExitCode(testCase.err)
			assert.Equal(t, testCase.exitCode, exitCode)
		})
	}
}
