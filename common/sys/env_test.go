package sys

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envVarNameUsedForTest = "test-parse-debug-mode-test"

func TestGetBoolEnv(t *testing.T) {
	// We can't run this test in parallel, as it sets env vars
	// t.Parallel()

	var testCases = []struct {
		envVarValue string
		fallback    bool
		expected    bool
	}{
		{"", false, false},
		{"foo", false, false},
		{"false", false, false},
		{"False", false, false},
		{"FALSE", false, false},
		{"true", false, true},
		{"True", false, true},
		{"TRUE", false, true},
		{"", true, true},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(fmt.Sprintf("val-%s-expected-%t", testCase.envVarValue, testCase.expected), func(t *testing.T) {
			os.Setenv(envVarNameUsedForTest, testCase.envVarValue)
			defer os.Unsetenv(envVarNameUsedForTest)

			actual := GetBoolEnv(envVarNameUsedForTest, testCase.fallback)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
