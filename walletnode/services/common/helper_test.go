package common

import (
	b64 "encoding/base64"
	"testing"

	json "github.com/json-iterator/go"

	"github.com/stretchr/testify/assert"
)

func TestInIntRange(t *testing.T) {
	t.Parallel()

	min := 5
	max := 10
	tests := map[string]struct {
		min     *int
		max     *int
		val     int
		inRange bool
	}{
		"simple":            {min: &min, max: &max, val: 7, inRange: true},
		"min inclusive":     {min: &min, max: &max, val: 5, inRange: true},
		"max inclusive":     {min: &min, max: &max, val: 10, inRange: true},
		"min false":         {min: &min, max: &max, val: 4, inRange: false},
		"max false":         {min: &min, max: &max, val: 11, inRange: false},
		"min nil":           {min: nil, max: &max, val: 9, inRange: true},
		"max nil":           {min: &min, max: nil, val: 11, inRange: true},
		"min max nil lower": {min: nil, max: nil, val: -1, inRange: true},
		"min max nil upper": {min: nil, max: nil, val: 1000, inRange: true},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := InIntRange(tc.val, tc.min, tc.max)
			assert.Equal(t, tc.inRange, got)
		})
	}
}

func TestInFloatRange(t *testing.T) {
	t.Parallel()

	min := 5.0
	max := 10.0
	tests := map[string]struct {
		min     *float64
		max     *float64
		val     float64
		inRange bool
	}{
		"simple":            {min: &min, max: &max, val: 7.0, inRange: true},
		"min inclusive":     {min: &min, max: &max, val: 5.0, inRange: true},
		"max inclusive":     {min: &min, max: &max, val: 10.0, inRange: true},
		"min false":         {min: &min, max: &max, val: 4.0, inRange: false},
		"max false":         {min: &min, max: &max, val: 11.0, inRange: false},
		"min nil":           {min: nil, max: &max, val: 9.0, inRange: true},
		"max nil":           {min: &min, max: nil, val: 11.0, inRange: true},
		"min max nil lower": {min: nil, max: nil, val: -1.0, inRange: true},
		"min max nil upper": {min: nil, max: nil, val: 1000.0, inRange: true},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := InFloatRange(tc.val, tc.min, tc.max)
			assert.Equal(t, tc.inRange, got)
		})
	}
}

func TestFromBase64(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		TestVarStr string
		TestVarInt int
	}

	testVarA := testStruct{TestVarStr: "test string /a", TestVarInt: -1}
	testBytesA, err := json.Marshal(testVarA)
	assert.Nil(t, err)

	testVarB := testStruct{TestVarStr: "", TestVarInt: 100}
	testBytesB, err := json.Marshal(testVarB)
	assert.Nil(t, err)

	tests := map[string]struct {
		encodedStr string
		out        *testStruct
		want       *testStruct
	}{
		"a": {encodedStr: b64.StdEncoding.EncodeToString([]byte(testBytesA)), out: &testStruct{}, want: &testVarA},
		"b": {encodedStr: b64.StdEncoding.EncodeToString([]byte(testBytesB)), out: &testStruct{}, want: &testVarB},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Nil(t, FromBase64(tc.encodedStr, tc.out))
			assert.Equal(t, tc.want.TestVarStr, tc.out.TestVarStr)
			assert.Equal(t, tc.want.TestVarInt, tc.out.TestVarInt)
		})
	}
}
