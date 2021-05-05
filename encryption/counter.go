package encryption

import "errors"

const counterLen = 12

var (
	errInvalidCounter = errors.New("invalid counter")
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
		return nil, errInvalidCounter
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

// NewOutCounter returns an outgoing counter initialized to the starting sequence
// number for the client/server side of a connection.
func NewOutCounter(overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	c.value[counterLen-1] = 0x80
	return
}

// NewInCounter returns an incoming counter initialized to the starting sequence
// number for the client/server side of a connection. This is used in ALTS record
// to check that incoming counters are as expected, since ALTS record guarantees
// that messages are unwrapped in the same order that the peer wrapped them.
func NewInCounter(overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	c.value[counterLen-1] = 0x80
	return
}
