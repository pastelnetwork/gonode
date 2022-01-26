package nftregister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifiers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		rqIDs RQIDSList
		res   []string
	}{
		"simple": {
			rqIDs: RQIDSList{&RQIDS{ID: "a"}, &RQIDS{ID: "b"}},
			res:   []string{"a", "b"},
		},
		"empty": {
			rqIDs: RQIDSList{&RQIDS{}},
			res:   []string{""},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			identifiers := tc.rqIDs.Identifiers()
			assert.Equal(t, tc.res, identifiers)
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		rqIDs RQIDSList
		res   map[string][]byte
	}{
		"simple": {
			rqIDs: RQIDSList{&RQIDS{ID: "a", Content: []byte("abc")}, &RQIDS{ID: "b", Content: []byte("abc")}},
			res:   map[string][]byte{"a": []byte("abc"), "b": []byte("abc")},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			idMap := tc.rqIDs.ToMap()
			assert.Equal(t, tc.res, idMap)
		})
	}
}
