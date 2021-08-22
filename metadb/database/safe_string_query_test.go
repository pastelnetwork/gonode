package database

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_safeStringQueryStruct(t *testing.T) {
	type testStruct struct {
		String      string
		Float64     float64
		Float64Star *float64
		Int         int
		IntStart    *int
	}
	hackString := "abc213ef ' OR 1==1 --"
	f10 := float64(10.0)
	i10 := int(10)

	processedHackString := "abc213ef '' OR 1==1 --"
	tests := []struct {
		v    interface{}
		want interface{}
	}{
		{
			v: &testStruct{
				String:      hackString,
				Float64:     f10,
				Float64Star: &f10,
				Int:         i10,
				IntStart:    &i10,
			},
			want: &testStruct{
				String:      processedHackString,
				Float64:     f10,
				Float64Star: &f10,
				Int:         i10,
				IntStart:    &i10,
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test_safeStringQueryStruct-%d", i), func(t *testing.T) {
			if got := safeStringQueryStruct(tt.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("safeStringQueryStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
