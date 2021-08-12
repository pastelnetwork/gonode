package conn

import "github.com/pastelnetwork/gonode/common/net/credentials/alts"

const bitMask = 0x80

// NewOutCounter returns an outgoing counter initialized to the starting sequence
// number for the client/server side of a connection.
func NewOutCounter(s alts.Side, overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	if s == alts.ServerSide {
		// Server counters in ALTS record have the little-endian high bit
		// set.
		c.value[counterLen-1] = bitMask
	}
	return
}

// NewInCounter returns an incoming counter initialized to the starting sequence
// number for the client/server side of a connection. This is used in ALTS record
// to check that incoming counters are as expected, since ALTS record guarantees
// that messages are unwrapped in the same order that the peer wrapped them.
func NewInCounter(s alts.Side, overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	if s == alts.ClientSide {
		// Server counters in ALTS record have the little-endian high bit
		// set.
		c.value[counterLen-1] = bitMask
	}
	return
}

// CounterFromValue creates a new counter given an initial value.
func CounterFromValue(value []byte, overflowLen int) (c Counter) {
	c.overflowLen = overflowLen
	copy(c.value[:], value)
	return
}

// CounterSide returns the connection side (client/server) a sequence counter is
// associated with.
func CounterSide(c []byte) alts.Side {
	if c[counterLen-1]&bitMask == bitMask {
		return alts.ServerSide
	}
	return alts.ClientSide
}
