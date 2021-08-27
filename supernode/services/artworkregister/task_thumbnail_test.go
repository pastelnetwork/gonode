package artworkregister

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestDeterminePreviewQuality(t *testing.T) {
	testCases := map[string]struct {
		size int
		want float32
	}{
		"lower": {
			size: 1000,
			want: 30.0,
		},
		"upper": {
			size: 2000,
			want: 10.0,
		},
		"lower-end": {
			size: 999,
			want: 30.0,
		},
		"bw": {
			size: 1500,
			want: 10.0 + 20.0*float32(2000-1500)/1000.0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := determinePreviewQuality(tc.size)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDetermineMediumQuality(t *testing.T) {
	testCases := map[string]struct {
		size int
		want float32
	}{
		"a": {
			size: 100,
			want: 100,
		},
		"b": {
			size: 399,
			want: 50.0,
		},
		"c": {
			size: 2000,
			want: 10.0,
		},
		"d": {
			size: 1000,
			want: 10.0 + 20.0*float32(2000-1000)/1000.0,
		},
		"e": {
			size: 401,
			want: 30.0,
		},
		"f": {
			size: 999,
			want: 30.0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := determineMediumQuality(tc.size)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMaxInt(t *testing.T) {
	testCases := map[string]struct {
		x, y, want int
	}{
		"x": {
			x:    1000,
			y:    30,
			want: 1000,
		},
		"y": {
			x:    312,
			y:    313,
			want: 313,
		},
		"equal": {
			x:    786,
			y:    786,
			want: 786,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := maxInt(tc.x, tc.y)
			assert.Equal(t, tc.want, got)
		})
	}
}
