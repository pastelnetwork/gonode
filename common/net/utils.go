package net

import (
	"encoding/binary"
	"fmt"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseFramedMsg parse the provided buffer and returns a frame of the format
// msgLength+msg and any remaining bytes in that buffer.
func ParseFramedMsg(b []byte, maxLen uint32) ([]byte, []byte, error) {
	// If the size field is not complete, return the provided buffer as
	// remaining buffer.
	if len(b) < MsgLenFieldSize {
		return nil, b, nil
	}
	msgLenField := b[:MsgLenFieldSize]
	length := binary.LittleEndian.Uint32(msgLenField)
	if length > maxLen {
		return nil, nil, fmt.Errorf("received the frame length %d larger than the limit %d", length, maxLen)
	}
	if len(b) < int(length)+4 { // account for the first 4 msg length bytes.
		// Frame is not complete yet.
		return nil, b, nil
	}
	return b[:MsgLenFieldSize+length], b[MsgLenFieldSize+length:], nil
}

func prepareLine(r []byte) []byte {
	return append(r, []byte("\n")...)
}