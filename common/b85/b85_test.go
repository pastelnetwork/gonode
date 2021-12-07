package b85

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in []byte
	}{
		"simple": {
			in: []byte("test-string-to-be-b85encoded"),
		},
		"complex": {
			in: []byte("test-string-complex-<*&%$@!+=*()[];:'"),
		},
		"numeric": {
			in: []byte("1234567890"),
		},
		"alpha-numeric-complex": {
			in: []byte("abcdefghijklmnopqrstuvwABCDEFGHIJKLMNOP!^&%*(&^%$$;:'`~1234567890"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			encoded := Encode(tc.in)

			decoded, err := Decode(encoded)
			assert.Nil(t, err)
			assert.Equal(t, tc.in, decoded)
		})
	}
}

func TestEncode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   []byte
		want string
	}{
		"simple": {
			in:   []byte("test-string-to-be-b85encoded"),
			want: "FCfN8/TZ#SBl7Q8FDia?AM%@N2.^Z8De*Ei",
		},
		"complex": {
			in:   []byte("test-string-complex-<*&%$@!+=*()[];:'"),
			want: "FCfN8/TZ#SBl7Q8@rH4'Ch7iC4=V[(,X<M'4Xqj/>?s<O-N",
		},

		"numeric": {
			in:   []byte("1234567890"),
			want: "0etOA2)[BQ3A:",
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			encoded := Encode(tc.in)
			assert.Equal(t, tc.want, encoded)

			/*
				decoded, err := Decode(encoded)
				assert.Nil(t, err)
				assert.Equal(t, tc.in, decoded)
			*/
		})
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want []byte
	}{
		"simple": {
			in:   "FCfN8/TZ#SBl7Q8FDia?AM%@N2.^Z8De*Ei",
			want: []byte("test-string-to-be-b85encoded"),
		},
		"complex": {
			in:   "FCfN8/TZ#SBl7Q8@rH4'Ch7iC4=V[(,X<M'4Xqj/>?s<O-N",
			want: []byte("test-string-complex-<*&%$@!+=*()[];:'"),
		},

		"numeric": {
			want: []byte("1234567890"),
			in:   "0etOA2)[BQ3A:F5",
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			decoded, err := Decode(tc.in)
			assert.Nil(t, err)

			assert.Equal(t, tc.want, decoded)
		})
	}
}
