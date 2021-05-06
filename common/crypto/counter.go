package crypto

import (
	"errors"
)

var (
	ErrInvalidCounter = errors.New("invalid counter")
)

const (
	counterLen = 12
)

// Counter is a 96-bit, little-endian counter.
type Counter struct {
	value       [counterLen]byte
	invalid     bool
	overflowLen int
}

// Value returns the current value of the counter as a byte slice.
func (c *Counter) Value() ([]byte, error) {
	if c.invalid {
		return nil, ErrInvalidCounter
	}
	return c.value[:], nil
}

// Inc increments the counter and checks for overflow.
func (c *Counter) Inc() {
	// If the counter is already invalid, there is no need to increase it.
	if c.invalid {
		return
	}
	i := 0
	for ; i < c.overflowLen; i++ {
		c.value[i]++
		if c.value[i] != 0 {
			break
		}
	}
	if i == c.overflowLen {
		c.invalid = true
	}
}

// NewCounter returns a counter initialized to the starting sequence
// number for the client/server side of a connection.
func NewCounter(overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	c.value[counterLen-1] = 0x80
	return
}
