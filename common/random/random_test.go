package random

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringIsMostlyRandom(t *testing.T) {
	t.Parallel()

	// Ensure that there is no overlap in 32 character random strings generated 100 times
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		newStr, err := String(32, Base62Chars)
		require.NoError(t, err)
		_, hasSeen := seen[newStr]
		require.False(t, hasSeen)
		seen[newStr] = true
	}
}

func TestStringRespectsStrLen(t *testing.T) {
	t.Parallel()

	for i := 0; i < 40; i++ {
		newStr, err := String(i, Base62Chars)
		require.NoError(t, err)
		assert.Equal(t, len(newStr), i)
	}
}

func TestStringRespectsAllowedChars(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		newStr, err := String(10, Digits)
		require.NoError(t, err)
		// Since the new string should only be composed of digits, if String respects allowed chars you should
		// always be able to convert the string to an int
		_, err = strconv.Atoi(newStr)
		require.NoError(t, err)
	}
}
