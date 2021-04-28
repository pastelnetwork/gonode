// Package legroast imports LegRoast C API
package legroast

/*
#cgo CFLAGS: -g -Wall -I./include
#cgo windows,amd64 LDFLAGS: -L./lib -llegroast-win
#cgo linux,amd64 LDFLAGS: -static -L./lib -llegroast
#cgo darwin LDFLAGS: -L./lib -llegroast-darwin
#include "sign.h"
*/
import "C"

const (
	PKBytes = C.PK_BYTES
	SKBytes = C.SK_BYTES
)

// Keygen generates private and public keys
func Keygen() ([]byte, []byte) {
	pk := make([]byte, C.PK_BYTES)
	sk := make([]byte, C.SK_BYTES)

	C.keygen((*C.uchar)(&pk[0]), (*C.uchar)(&sk[0]))

	return pk, sk
}

// Sign signs the passed message with given keys
// It returns signed message data
func Sign(pk []byte, sk []byte, message []byte) []byte {
	signal := make([]byte, C.SIG_BYTES)
	var signal_len C.uint64_t

	C.sign((*C.uchar)(&sk[0]), (*C.uchar)(&pk[0]), (*C.uchar)(&message[0]), (C.uint64_t)(len(message)), (*C.uchar)(&signal[0]), (*C.uint64_t)(&signal_len))

	return signal
}

// Verify validates previously signed message
// It returns <= 0 if verification has failed
func Verify(pk []byte, message []byte, signal []byte) int {
	return (int)(C.verify((*C.uchar)(&pk[0]), (*C.uchar)(&message[0]), (C.uint64_t)(len(message)), (*C.uchar)(&signal[0])))
}
