// Package random provides utilities and functions for generating random data.
package random

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Character sets that you can use when passing into String
const (
	Digits       = "0123456789"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"
	SpecialChars = "<>[]{}()-_*%&/?\"'\\"
)

// Base62Chars is set of Base64 charsets
var Base62Chars = Digits + UpperLetters + LowerLetters

// String generates a random string of length strLength, composing only of characters in allowedChars. Based on
// code here: http://stackoverflow.com/a/9543797/483528
// For convenience, the random package exposes various character sets you can use for the allowedChars parameter. Here
// are a few examples:
//
// // Only lower case chars + digits
// random.String(6, random.Digits + random.LowerLetters)
//
// // alphanumerics + special chars
// random.String(6, random.Base62Chars + random.SpecialChars)
//
// // Only alphanumerics (base62)
// random.String(6, random.Base62Chars)
//
// // Only abc
// random.String(6, "abc")
func String(strLength int, allowedChars string) (string, error) {
	var out bytes.Buffer

	for i := 0; i < strLength; i++ {
		id, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		if err != nil {
			return out.String(), errors.New(err)
		}
		out.WriteByte(allowedChars[id.Int64()])
	}

	return out.String(), nil
}
